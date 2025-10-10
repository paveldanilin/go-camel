package template

import (
	"testing"
)

type service struct {
	ID   string
	Name string
}

type product struct {
	ID   string
	Name string
}

type profile struct {
	ID       string
	Services []service
	Products []product
}

func TestTemplate(t *testing.T) {
	data := map[string]any{
		"a": "valueA",
		"person": map[string]any{
			"name":    "Alice",
			"age":     30,
			"friends": []any{"Bob", "Charlie"},
			"nestedList": []any{
				map[string]any{"item": "first"},
				map[string]any{"item": "second"},
			},
		},
		"special": "${escaped}",
		"profile": profile{
			ID: "XYZ-01",
			Services: []service{
				{"S001", "SRV-001"},
				{"S002", "SRV-002"},
				{"S003", "SRV-003"},
			},
			Products: []product{
				{"P0001", "PROD-001"},
				{"P0002", "PROD-002"},
				{"P0003", "PROD-003"},
			},
		},
	}

	templateStr := "Hello, ${person.name}! Age: ${person.age}. Friend: ${person.friends[1]}. Nested: ${person.nestedList[1].item}. Special: ${special}. Profile: ${profile.id}; ${profile.services[0].name}; ${profile.products[1].name};"
	tmpl, err := Parse(templateStr)
	if err != nil {
		t.Errorf("Parse error: %s", err)
	}

	result, err := tmpl.Render(data)
	if err != nil {
		t.Errorf("Render error: %s", err)
	}

	expectedString := "Hello, Alice! Age: 30. Friend: Charlie. Nested: second. Special: \\${escaped\\}. Profile: XYZ-01; SRV-001; PROD-002;"
	if expectedString != result {
		t.Errorf("expected string '%s', but got '%s'", expectedString, result)
	}
}

func TestHasVars_NoVars(t *testing.T) {
	input := "Just simple string with $ sign {x}"

	result := HasVars(input)
	expectedValue := false

	if expectedValue != result {
		t.Errorf("expected bool value %v, but got %v", expectedValue, result)
	}
}

func TestHasVars_Escaped(t *testing.T) {
	input := "This is simple string \\${xxx}"

	result := HasVars(input)
	expectedValue := false

	if expectedValue != result {
		t.Errorf("expected bool value %v, but got %v", expectedValue, result)
	}
}

func TestHasVars(t *testing.T) {
	input := "Hello ${user_name}"

	result := HasVars(input)
	expectedValue := true

	if expectedValue != result {
		t.Errorf("expected bool value %v, but got %v", expectedValue, result)
	}
}
