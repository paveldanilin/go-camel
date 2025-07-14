package processor

import (
	"github.com/paveldanilin/go-camel/camel"
)

type ChoiceCase struct {
	expr      camel.Expr
	processor camel.Processor
}

func NewChoiceCase(expr camel.Expr, processor camel.Processor) *ChoiceCase {

	return &ChoiceCase{expr: expr, processor: processor}
}

func (cc *ChoiceCase) IsMatch(message *camel.Message) bool {

	v, err := cc.expr.Eval(message)
	if err != nil {
		panic(err)
	}

	return v.(bool)
}

type ChoiceProcessor struct {
	cases []*ChoiceCase
}

func Choice() *ChoiceProcessor {
	return &ChoiceProcessor{
		cases: []*ChoiceCase{},
	}
}

func (p *ChoiceProcessor) AddCase(c *ChoiceCase) *ChoiceProcessor {

	p.cases = append(p.cases, c)
	return p
}

func (p *ChoiceProcessor) Process(message *camel.Message) error {

	for _, choiceCase := range p.cases {
		if choiceCase.IsMatch(message) {
			choiceCase.processor.Process(message)
		}
	}

	return nil
}

