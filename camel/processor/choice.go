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

func (p *ChoiceProcessor) When(predicate camel.Expr, processor camel.Processor) *ChoiceProcessor {
	p.cases = append(p.cases, &choiceWhen{predicate: predicate, processor: processor})
	return p
}

func (p *ChoiceProcessor) Otherwise(processor camel.Processor) *ChoiceProcessor {
	p.otherwise = processor
	return p
}

func (p *ChoiceProcessor) Process(exchange *camel.Exchange) {
	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	whenMatched := false

	for _, whenCase := range p.cases {
		if whenCase.match(exchange) {
			whenMatched = true
			whenCase.processor.Process(exchange)
			break
		}
	}

	if !whenMatched && p.otherwise != nil {
		if err := exchange.CheckCancelOrTimeout(); err != nil {
			exchange.Error = err
			return
		}
		p.otherwise.Process(exchange)
	}
}

// choiceWhen represents a single when check of ChoiceProcessor
type choiceWhen struct {
	predicate camel.Expr
	processor camel.Processor
}

func (when *choiceWhen) match(exchange *camel.Exchange) bool {
	v, err := when.predicate.Eval(exchange)
	if err != nil {
		panic(err)
	}

	return v.(bool)
}
