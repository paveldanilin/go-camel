package camel

import (
	"testing"
)

type jsonModel struct {
	Page   int
	Fruits []string
}

func TestJsonProcessorMarshal(t *testing.T) {
	j := JsonProcessor{model: jsonModel{Page: 2, Fruits: []string{"orange", "abricot"}}, operation: "marshal"}
	jsonProcMarshal := newJsonProcessor(j)
	exchange := NewExchange(nil, nil)
	expectedBody := `{"Page":2,"Fruits":["orange","abricot"]}`

	jsonProcMarshal.Process(exchange)

	if exchange.Message().Body != expectedBody {
		t.Errorf("TestJsonProcessorMarshal() = %v; want body %v", exchange.Message().Body, expectedBody)
	}
}

// func TestJsonProcessorUnmarshal(t *testing.T) {
// 	// var body map[string]interface{}
// 	j := JsonProcessor{json: `{"Page": 1, "Fruits": ["apple", "peach"]}`, model: make(map[string]*jsonModel), operation: "unmarshal"}
// 	jsonProcUnmarshal := newJsonProcessor(j)
// 	exchange := NewExchange(nil, nil)
// 	// expectedBody := map[Fruits:[apple peach] Page:1]

// 	jsonProcUnmarshal.Process(exchange)

// 	// b, ok := exchange.Message().Body.(map[string]interface{})
// 	// if ok {
// 	// 	body = b
// 	// }

// 	// t.Log(body, expectedBody)

// 	// if body["page"] != expectedBody {
// 	// 	t.Errorf("TestJsonProcessorUnmarshal() = %v; want body %v", body["page"], expectedBody)
// 	// }
// }
