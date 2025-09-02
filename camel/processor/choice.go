package processor

import (
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/expr"
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
	p.cases = append(p.cases, &choiceWhen{predicate: expr.Predicate(predicate), processor: processor})
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
		caseMatched, err := whenCase.match(exchange)
		if err != nil {
			// In case of error stop processing choice
			exchange.Error = err
			return
		}

		if caseMatched {
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
	predicate camel.Predicate
	processor camel.Processor
}

func (when *choiceWhen) match(exchange *camel.Exchange) (bool, error) {
	return when.predicate.Test(exchange)
}
