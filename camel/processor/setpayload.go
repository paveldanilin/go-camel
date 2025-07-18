package processor

import "github.com/paveldanilin/go-camel/camel"

type SetPayloadProcessor struct {
	value camel.Expr
}

func SetPayload(value camel.Expr) *SetPayloadProcessor {
	return &SetPayloadProcessor{
		value: value,
	}
}

func (p *SetPayloadProcessor) Process(message *camel.Message) {

	value, err := p.value.Eval(message)
	if err != nil {
		message.SetError(err)
		return
	}

	message.SetPayload(value)
}
