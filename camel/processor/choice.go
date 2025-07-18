package processor

import (
	"github.com/paveldanilin/go-camel/camel"
)

type ChoiceProcessor struct {
	cases     []*choiceWhen
	otherwise camel.Processor
}

func Choice() *ChoiceProcessor {

	return &ChoiceProcessor{
		cases: []*choiceWhen{},
	}
}

type choiceWhen struct {
	cond      camel.Expr
	processor camel.Processor
}

func (when *choiceWhen) isMatch(message *camel.Message) bool {

	v, err := when.cond.Eval(message)
	if err != nil {
		panic(err)
	}

	return v.(bool)
}

func (p *ChoiceProcessor) When(cond camel.Expr, processor camel.Processor) *ChoiceProcessor {

	p.cases = append(p.cases, &choiceWhen{cond: cond, processor: processor})
	return p
}

func (p *ChoiceProcessor) Otherwise(processor camel.Processor) *ChoiceProcessor {

	p.otherwise = processor
	return p
}

func (p *ChoiceProcessor) Process(message *camel.Message) {

	whenMatched := false
	for _, when := range p.cases {
		if when.isMatch(message) {
			whenMatched = true
			when.processor.Process(message)
			break
		}
	}

	if !whenMatched && p.otherwise != nil {
		p.otherwise.Process(message)
	}
}
