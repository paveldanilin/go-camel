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
	predicate camel.Expr
	processor camel.Processor
}

func (when *choiceWhen) match(message *camel.Message) bool {

	v, err := when.predicate.Eval(message)
	if err != nil {
		panic(err)
	}

	return v.(bool)
}

func (p *ChoiceProcessor) When(predicate camel.Expr, processor camel.Processor) *ChoiceProcessor {

	p.cases = append(p.cases, &choiceWhen{predicate: predicate, processor: processor})
	return p
}

func (p *ChoiceProcessor) Otherwise(processor camel.Processor) *ChoiceProcessor {

	p.otherwise = processor
	return p
}

func (p *ChoiceProcessor) Process(message *camel.Message) {

	whenMatched := false
	for _, when := range p.cases {
		if when.match(message) {
			whenMatched = true
			when.processor.Process(message)
			break
		}
	}

	if !whenMatched && p.otherwise != nil {
		p.otherwise.Process(message)
	}
}
