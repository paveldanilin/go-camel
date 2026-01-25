package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

// ChoiceWhen represents a single 'when' check of Choice
type ChoiceWhen struct {
	predicate Predicate
	processor Processor
}

func (when *ChoiceWhen) match(e *exchange.Exchange) (bool, error) {
	return when.predicate.Test(e)
}

type Choice struct {
	cases     []*ChoiceWhen
	otherwise Processor
}

func NewChoice() *Choice {
	return &Choice{
		cases: []*ChoiceWhen{},
	}
}

func (p *Choice) AddWhen(predicate Predicate, then Processor) error {
	if predicate == nil {
		return fmt.Errorf("predicate is nil")
	}
	if then == nil {
		return fmt.Errorf("then processor is nil")
	}
	p.cases = append(p.cases, &ChoiceWhen{
		predicate: predicate,
		processor: then,
	})
	return nil
}

func (p *Choice) SetOtherwise(otherwise Processor) *Choice {
	p.otherwise = otherwise
	return p
}

func (p *Choice) Process(e *exchange.Exchange) {
	for _, whenCase := range p.cases {
		caseMatched, err := whenCase.match(e)
		if err != nil {
			// In case of error stop processing choice
			e.SetError(err)
			return
		}

		if caseMatched {
			whenCase.processor.Process(e)
			return
		}
	}

	// No one case was found
	if p.otherwise != nil {
		p.otherwise.Process(e)
	}
}
