package camel

import (
	"encoding/json"
	"fmt"
)

type JsonProcessor struct {
	json      string
	model     any
	operation string
}

func newJsonProcessor(j JsonProcessor) *JsonProcessor {
	return &JsonProcessor{
		json:      j.json,
		model:     j.model,
		operation: j.operation,
	}
}

// type jsonString struct {
// 	value string
// 	model any
// }

// type jsonBits struct {
// 	value []byte
// 	model any
// }

// type JsonOperation interface {
// 	marshal() string
// 	unmarshal() any
// }

func (j *JsonProcessor) unmarshal() {
	if err := json.Unmarshal([]byte(j.json), &j.model); err != nil {
		fmt.Println("error:", err)
	}
}

// func (j *jsonBits) unmarshal() {
// 	fmt.Println(2)
// 	if err := json.Unmarshal(j.value, &j.model); err != nil {
// 		fmt.Println("error:", err)
// 	}
// }

func (j *JsonProcessor) marshal() {
	data, err := json.Marshal(&j.model)
	if err != nil {
		fmt.Println("json error:", err)
	} else {
		j.json = string(data)
		// fmt.Println("string(res)", data, string(data))
	}
}

func (j *JsonProcessor) handleMarshal() {
	fmt.Println("MARSHAL OPERATION")
	j.marshal()
}

func (j *JsonProcessor) handleUnnmarshal() {
	fmt.Println("UNMARSHAL OPERATION")
	j.unmarshal()
}

func (j *JsonProcessor) Process(exchange *Exchange) {
	switch j.operation {
	case "marshal":
		j.handleMarshal()
		exchange.Message().Body = j.json
	case "unmarshal":
		j.handleUnnmarshal()
		exchange.Message().Body = j.model
	default:
		fmt.Println("operation should be marshal or unmarshal")
	}
}
