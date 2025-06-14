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

func (p *SetPayloadProcessor) Process(message *camel.Message) error {

	value, err := p.value.Eval(message)
	if err != nil {
		return nil
	}

	message.SetPayload(value)

	return nil
}
