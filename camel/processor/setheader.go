package processor

import "github.com/paveldanilin/go-camel/camel"

type SetHeaderProcessor struct {
	name  string
	value camel.Expr
}

func SetHeader(name string, value camel.Expr) *SetHeaderProcessor {
	return &SetHeaderProcessor{
		name:  name,
		value: value,
	}
}

func (p *SetHeaderProcessor) Process(message *camel.Message) {

	value, err := p.value.Eval(message)
	if err != nil {
		message.SetError(err)
		return
	}

	message.SetHeader(p.name, value)
}
