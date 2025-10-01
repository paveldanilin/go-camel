package camel

import (
	"fmt"
)

// choiceWhen represents a single 'when' check of choiceProcessor
type choiceWhen struct {
	predicate Predicate
	processor Processor
}

func (when *choiceWhen) match(exchange *Exchange) (bool, error) {
	return when.predicate.Test(exchange)
}

type choiceProcessor struct {
	id        string
	cases     []*choiceWhen
	otherwise Processor
}

func newChoiceProcessor(id string) *choiceProcessor {
	return &choiceProcessor{
		id:    id,
		cases: []*choiceWhen{},
	}
}

func (p *choiceProcessor) getId() string {
	return p.id
}

func (p *choiceProcessor) addWhen(predicate Expr, processor Processor) *choiceProcessor {
	if predicate == nil {
		panic(fmt.Errorf("camel: choice: when: predicate must be not nil"))
	}
	if processor == nil {
		panic(fmt.Errorf("camel: choice: when: processor must be not nil"))
	}
	p.cases = append(p.cases, &choiceWhen{
		predicate: newPredicateFromExpr(predicate),
		processor: processor,
	})
	return p
}

func (p *choiceProcessor) setOtherwise(processor Processor) *choiceProcessor {
	p.otherwise = processor
	return p
}

func (p *choiceProcessor) Process(exchange *Exchange) {
	for _, whenCase := range p.cases {
		caseMatched, err := whenCase.match(exchange)
		if err != nil {
			// In case of error stop processing choice
			exchange.SetError(err)
			return
		}

		if caseMatched {
			whenCase.processor.Process(exchange)
			return
		}
	}

	// No one case was found
	if p.otherwise != nil {
		p.otherwise.Process(exchange)
	}
}
