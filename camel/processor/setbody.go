package processor

import "github.com/paveldanilin/go-camel/camel"

type SetBodyProcessor struct {
	value camel.Expr
}

func SetBody(value camel.Expr) *SetBodyProcessor {
	return &SetBodyProcessor{
		value: value,
	}
}

func (p *SetBodyProcessor) Process(message *camel.Message) {

	value, err := p.value.Eval(message)
	if err != nil {
		message.Error = err
		return
	}

	message.Body = value
}
