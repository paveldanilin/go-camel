package template

import (
	"bytes"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"strings"
)

type Template struct {
	tmpl     *template.Template
	vars     []string
	template string
}

func (t *Template) Render(data map[string]any) (string, error) {
	var buf bytes.Buffer
	err := t.tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Vars returns an unique array of variables (${var_name}).
func (t *Template) Vars() []string {
	return t.vars
}

// Template returns original template as string.
func (t *Template) Template() string {
	return t.template
}

// Parse parses input string and create Template instance.
// Variable delimiters are: '${' - begin delimiter, '}' - end delimiter.
//
//	Example: `Hello, ${username}!`
func Parse(input string) (*Template, error) {
	var builder strings.Builder
	vars := map[string]struct{}{}
	resolvers := make(map[string][]pathStep)
	runes := []rune(input)
	n := len(runes)
	i := 0

	for i < n {
		if runes[i] == '$' && i+1 < n && runes[i+1] == '{' {
			i += 2
			pathStart := i
			for ; i < n && runes[i] != '}'; i++ {
			}
			if i == n {
				return nil, fmt.Errorf("unclosed variable at position %d", pathStart-2)
			}
			path := string(runes[pathStart:i])
			if path == "" {
				return nil, fmt.Errorf("empty variable at position %d", pathStart-2)
			}

			steps, err := parsePath(path)
			if err != nil {
				return nil, err
			}
			resolvers[path] = steps

			vars[path] = struct{}{}

			// {{lookup .path}}
			builder.WriteString(`{{lookup . "`)
			builder.WriteString(strings.ReplaceAll(path, `"`, `\"`))
			builder.WriteString(`"}}`)
			i++
			continue
		}
		builder.WriteRune(runes[i])
		i++
	}

	// use Closure to capture resolver
	lookupFunc := func(ctx any, path string) (string, error) {
		return lookup(ctx, path, resolvers)
	}

	tmpl := template.New("rendered").Funcs(template.FuncMap{
		"lookup": lookupFunc,
	})
	parsed, err := tmpl.Parse(builder.String())
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}

	varNames := make([]string, 0, len(vars))
	for varName := range vars {
		varNames = append(varNames, varName)
	}

	return &Template{
		tmpl:     parsed,
		vars:     varNames,
		template: input,
	}, nil
}

// HasVars checks if input string contains variables like '${var_name}'.
func HasVars(input string) bool {
	if len(input) < 3 { // "${}"
		return false
	}

	i := 0
	n := len(input)

	for {
		// search for next "$"
		dollarPos := strings.Index(input[i:], "$")
		if dollarPos == -1 {
			return false
		}
		i += dollarPos

		// escaping: count backslash before "$"
		escapeCount := 0
		for j := i - 1; j >= 0 && input[j] == '\\'; j-- {
			escapeCount++
		}
		if escapeCount%2 == 1 { // escaped \
			i++ // SKIP "$"
			continue
		}

		// check for  "${"
		if i+1 < n && input[i+1] == '{' {
			i += 2
			pathStart := i
			// search for next "}" (unescaped)
			foundClose := false
			for ; i < n; i++ {
				if input[i] == '}' {
					// check escaping before "}"
					escapeCount = 0
					for j := i - 1; j >= pathStart && input[j] == '\\'; j-- {
						escapeCount++
					}
					if escapeCount%2 == 0 { // unescaped
						foundClose = true
						break
					}
				}
			}
			if foundClose && i > pathStart { // valid var
				return true
			}
			// If not found or empty, proceed from the current i
			continue
		}
		i++ // SKIP single "$"
	}
}

func Vars(input string) ([]string, error) {
	vars := map[string]struct{}{}
	runes := []rune(input)
	n := len(runes)
	i := 0

	for i < n {
		if runes[i] == '$' && i+1 < n && runes[i+1] == '{' {
			i += 2
			pathStart := i
			for ; i < n && runes[i] != '}'; i++ {
			}
			if i == n {
				return nil, fmt.Errorf("unclosed variable at position %d", pathStart-2)
			}
			path := string(runes[pathStart:i])
			if path == "" {
				return nil, fmt.Errorf("empty variable at position %d", pathStart-2)
			}

			vars[path] = struct{}{}
			i++
			continue
		}
		i++
	}

	varNames := make([]string, 0, len(vars))
	for varName := range vars {
		varNames = append(varNames, varName)
	}

	return varNames, nil
}

func Render(input string, data map[string]any) (string, error) {
	t, err := Parse(input)
	if err != nil {
		return "", err
	}
	return t.Render(data)
}

type pathStep struct {
	isIndex bool
	key     string // for map/struct
	index   int    // for slice/array
}

func lookup(ctx any, path string, resolvers map[string][]pathStep) (string, error) {
	steps, ok := resolvers[path]
	if !ok {
		return "", fmt.Errorf("no resolver for %q (check path in template)", path)
	}
	val, err := resolve(ctx, steps)
	if err != nil {
		return "", err
	}
	strVal := fmt.Sprintf("%v", val)

	strVal = strings.ReplaceAll(strVal, "${", "\\${")
	strVal = strings.ReplaceAll(strVal, "}", "\\}")
	return strVal, nil
}

func resolve(ctx any, steps []pathStep) (any, error) {
	current := ctx
	for _, step := range steps {
		if step.isIndex {
			v := reflect.ValueOf(current)
			if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
				return nil, fmt.Errorf("not a slice for index %d", step.index)
			}
			if step.index >= v.Len() {
				return nil, fmt.Errorf("index %d out of bounds", step.index)
			}
			current = v.Index(step.index).Interface()
			continue
		}

		// first trying to use 'map[string]any' (faster)
		if m, ok := current.(map[string]any); ok {
			val, exists := m[step.key]
			if !exists {
				return nil, fmt.Errorf("key %q not found", step.key)
			}
			current = val
			continue
		}

		// use reflect for struct/map
		v := reflect.ValueOf(current)
		switch v.Kind() {
		case reflect.Struct:
			found := false
			for i := 0; i < v.NumField(); i++ {
				fieldName := v.Type().Field(i).Name
				if strings.EqualFold(fieldName, step.key) {
					current = v.Field(i).Interface()
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("field %q not found (case-insensitive)", step.key)
			}
		case reflect.Map:
			key := reflect.ValueOf(step.key)
			val := v.MapIndex(key)
			if !val.IsValid() {
				return nil, fmt.Errorf("key %q not found", step.key)
			}
			current = val.Interface()
		default:
			return nil, fmt.Errorf("unsupported type for key %q", step.key)
		}
	}
	return current, nil
}

func parsePath(path string) ([]pathStep, error) {
	var steps []pathStep
	runes := []rune(path)
	n := len(runes)
	i := 0

	for i < n {
		segStart := i
		for ; i < n && runes[i] != '.' && runes[i] != '['; i++ {
		}
		seg := string(runes[segStart:i])
		if seg == "" && !(i < n && runes[i] == '[') {
			return nil, fmt.Errorf("invalid path segment")
		}

		if i < n && runes[i] == '[' {
			i++
			idxStart := i
			for ; i < n && runes[i] != ']'; i++ {
			}
			if i == n || runes[i] != ']' {
				return nil, fmt.Errorf("unclosed index")
			}
			idxStr := string(runes[idxStart:i])
			idx, err := strconv.Atoi(idxStr)
			if err != nil || idx < 0 {
				return nil, fmt.Errorf("invalid index %q: %w", idxStr, err)
			}
			i++
			if seg != "" {
				steps = append(steps, pathStep{isIndex: false, key: seg})
			}
			steps = append(steps, pathStep{isIndex: true, index: idx})
		} else if seg != "" {
			steps = append(steps, pathStep{isIndex: false, key: seg})
		}

		if i < n && runes[i] == '.' {
			i++
		} else if i < n {
			return nil, fmt.Errorf("unexpected char after segment")
		}
	}
	return steps, nil
}
