package converter

import (
	"errors"
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"reflect"
	"sync"
)

type registry struct {
	mu         sync.RWMutex
	converters map[reflect.Type]map[reflect.Type]any // map[fromType][toType]Converter[From, To]
	cache      map[string]chainPath                  // cacheKey: "from.String():to.String()" -> chainPath
	named      map[string]reflect.Type
}

// chainPath - cache for cache of converters.
type chainPath struct {
	steps []chainStep
}

type chainStep struct {
	toType reflect.Type
	conv   any
}

func NewRegistry() *registry {
	return &registry{
		converters: make(map[reflect.Type]map[reflect.Type]any),
		cache:      make(map[string]chainPath),
		named:      make(map[string]reflect.Type),
	}
}

func (r *registry) Register(conv any) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	convType := reflect.TypeOf(conv)
	if convType.Kind() != reflect.Func && !convType.Implements(reflect.TypeOf((*api.Converter[any, any])(nil)).Elem()) {
		return errors.New("provided conv must implement Converter[From, To]")
	}

	// extract From and To from Converter[From, To]
	convertMethod := reflect.ValueOf(conv).MethodByName("Convert")
	if !convertMethod.IsValid() {
		return errors.New("converter must have Convert method")
	}

	fromType := convertMethod.Type().In(0)
	toType := convertMethod.Type().Out(0)

	if fromType == nil || toType == nil {
		return errors.New("cannot register converter with nil types")
	}

	if r.converters[fromType] == nil {
		r.converters[fromType] = make(map[reflect.Type]any)
	}
	r.converters[fromType][toType] = conv

	// invalidate cache
	r.cache = make(map[string]chainPath)
	r.named[toType.String()] = toType
	return nil
}

func (r *registry) Type(name string) (reflect.Type, bool) {
	t, exists := r.named[name]
	return t, exists
}

func (r *registry) CanConvert(fromType, toType reflect.Type) bool {
	_, err := r.findChain(fromType, toType)
	return err == nil
}

func (r *registry) Convert(value any, toType reflect.Type, params map[string]any) (any, error) {
	if value == nil {
		return nil, errors.New("cannot convert nil value")
	}
	fromType := reflect.TypeOf(value)
	if fromType == toType {
		return value, nil // Нет нужды в конвертации
	}

	// test cache
	r.mu.RLock()
	cacheKey := fmt.Sprintf("%s:%s", fromType.String(), toType.String())
	if path, ok := r.cache[cacheKey]; ok {
		r.mu.RUnlock()
		return r.executeChain(value, path, params)
	}
	r.mu.RUnlock()

	// get chain
	path, err := r.findChain(fromType, toType)
	if err != nil {
		return nil, err
	}

	// put cache
	r.mu.Lock()
	r.cache[cacheKey] = path
	r.mu.Unlock()

	return r.executeChain(value, path, params)
}

// findChain finds the shortes path of converters (BFS).
func (r *registry) findChain(start, target reflect.Type) (chainPath, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type node struct {
		current reflect.Type
		path    []chainStep
		depth   int
	}

	queue := []node{{current: start, path: nil, depth: 0}}
	visited := map[reflect.Type]bool{start: true}
	maxDepth := 5 // max depath limit

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.current == target {
			return chainPath{steps: curr.path}, nil
		}

		if curr.depth >= maxDepth {
			continue
		}

		if targets, ok := r.converters[curr.current]; ok {
			for nextType, conv := range targets {
				if !visited[nextType] {
					visited[nextType] = true
					newPath := append([]chainStep(nil), curr.path...)
					newPath = append(newPath, chainStep{toType: nextType, conv: conv})
					queue = append(queue, node{current: nextType, path: newPath, depth: curr.depth + 1})
				}
			}
		}
	}

	return chainPath{}, errors.New("no conversion path found for types: " + start.String() + " to " + target.String())
}

// executeChain executes the chain of converters.
func (r *registry) executeChain(value any, path chainPath, params map[string]any) (any, error) {
	current := value
	for _, step := range path.steps {
		convValue := reflect.ValueOf(step.conv)
		convertMethod := convValue.MethodByName("Convert")
		if !convertMethod.IsValid() {
			return nil, errors.New("invalid converter in chain")
		}

		expectedFromType := convertMethod.Type().In(0)
		currentType := reflect.TypeOf(current)
		if currentType != expectedFromType {
			if currentType.ConvertibleTo(expectedFromType) {
				current = reflect.ValueOf(current).Convert(expectedFromType).Interface()
			} else {
				return nil, fmt.Errorf("type mismatch in chain: expected %s, got %s", expectedFromType, currentType)
			}
		}

		results := convertMethod.Call([]reflect.Value{reflect.ValueOf(current), reflect.ValueOf(params)})
		if len(results) != 2 {
			return nil, errors.New("converter must return (To, error)")
		}

		if !results[1].IsNil() {
			return nil, results[1].Interface().(error)
		}

		current = results[0].Interface()
	}
	return current, nil
}
