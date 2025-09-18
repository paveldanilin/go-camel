package camel

import (
	"fmt"
)

type choiceProcessor struct {
	// stepName is a logical name of current operation.
	stepName  string
	cases     []*choiceWhen
	otherwise Processor
}

func newChoiceProcessor() *choiceProcessor {
	return &choiceProcessor{
		stepName: "choice{}",
		cases:    []*choiceWhen{},
	}
}

func (p *choiceProcessor) WithStepName(stepName string) *choiceProcessor {
	p.stepName = stepName
	return p
}

func (p *choiceProcessor) When(predicate Expr, processor Processor) *choiceProcessor {
	if predicate == nil {
		panic(fmt.Errorf("camel: choice: when: predicate must be not nil"))
	}
	if processor == nil {
		panic(fmt.Errorf("camel: choice: when: processor must be not nil"))
	}
	p.cases = append(p.cases, &choiceWhen{
		predicate: newPredicateExpr(predicate),
		processor: processor,
	})
	return p
}

func (p *choiceProcessor) Otherwise(processor Processor) *choiceProcessor {
	p.otherwise = processor
	return p
}

func (p *choiceProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

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

// choiceWhen represents a single 'when' check of choiceProcessor
type choiceWhen struct {
	predicate Predicate
	processor Processor
}

func (when *choiceWhen) match(exchange *Exchange) (bool, error) {
	return when.predicate.Test(exchange)
}
