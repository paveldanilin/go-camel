package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
	"github.com/paveldanilin/go-camel/camel/expr"
)

type ChoiceProcessor struct {
	// stepName is a logical name of current operation.
	stepName  string
	cases     []*choiceWhen
	otherwise camel.Processor
}

func Choice() *ChoiceProcessor {
	return &ChoiceProcessor{
		stepName: "choice{}",
		cases:    []*choiceWhen{},
	}
}

func (p *ChoiceProcessor) WithStepName(stepName string) *ChoiceProcessor {
	p.stepName = stepName
	return p
}

func (p *ChoiceProcessor) When(predicate camel.Expr, processor camel.Processor) *ChoiceProcessor {
	if predicate == nil {
		panic(fmt.Errorf("camel: choice: when: predicate must be not nil"))
	}
	if processor == nil {
		panic(fmt.Errorf("camel: choice: when: processor must be not nil"))
	}
	p.cases = append(p.cases, &choiceWhen{
		predicate: expr.Predicate(predicate),
		processor: processor,
	})
	return p
}

func (p *ChoiceProcessor) Otherwise(processor camel.Processor) *ChoiceProcessor {
	p.otherwise = processor
	return p
}

func (p *ChoiceProcessor) Process(exchange *camel.Exchange) {
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

// choiceWhen represents a single 'when' check of ChoiceProcessor
type choiceWhen struct {
	predicate camel.Predicate
	processor camel.Processor
}

func (when *choiceWhen) match(exchange *camel.Exchange) (bool, error) {
	return when.predicate.Test(exchange)
}
