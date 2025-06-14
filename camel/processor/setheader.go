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

func (p *SetHeaderProcessor) Process(message *camel.Message) error {

	value, err := p.value.Eval(message)
	if err != nil {
		return err
	}

	message.SetHeader(p.name, value)

	return nil
}
