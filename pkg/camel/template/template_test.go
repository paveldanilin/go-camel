package template

import (
	"github.com/stretchr/testify/assert"
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

func TestTemplate_Render(t *testing.T) {
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
	assert.NoError(t, err)

	result, err := tmpl.Render(data)
	assert.NoError(t, err)

	expectedString := "Hello, Alice! Age: 30. Friend: Charlie. Nested: second. Special: \\${escaped\\}. Profile: XYZ-01; SRV-001; PROD-002;"
	assert.Equal(t, expectedString, result)
}

func TestHasVars_NoVars(t *testing.T) {
	input := "Just simple string with $ sign {x}"

	result := HasVars(input)
	assert.False(t, result)
}

func TestHasVars_Escaped(t *testing.T) {
	input := "This is simple string \\${xxx}"

	result := HasVars(input)
	assert.False(t, result)
}

func TestHasVars(t *testing.T) {
	input := "Hello ${user_name}"

	result := HasVars(input)
	assert.True(t, result)
}

func TestVars_NoVars(t *testing.T) {
	input := "Hello, user!"

	vars, err := Vars(input)
	assert.NoError(t, err)
	assert.Empty(t, vars)
}

func TestVars(t *testing.T) {
	input := "Hello, ${username}! Today is ${day}."

	vars, err := Vars(input)
	assert.NoError(t, err)
	assert.Len(t, vars, 2)
}
