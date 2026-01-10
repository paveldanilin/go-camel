package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

// ChoiceWhen represents a single 'when' check of Choice
type ChoiceWhen struct {
	predicate expression.Predicate
	processor api.Processor
}

func (when *ChoiceWhen) match(e *exchange.Exchange) (bool, error) {
	return when.predicate.Test(e)
}

type Choice struct {
	cases     []*ChoiceWhen
	otherwise api.Processor
}

func NewChoice() *Choice {
	return &Choice{
		cases: []*ChoiceWhen{},
	}
}

func (p *Choice) When(predicate expression.Expression, processor api.Processor) *Choice {
	if predicate == nil {
		panic(fmt.Errorf("choice: when: predicate must be not nil"))
	}
	if processor == nil {
		panic(fmt.Errorf("choice: when: processor must be not nil"))
	}
	p.cases = append(p.cases, &ChoiceWhen{
		predicate: expression.NewPredicateFromExpression(predicate),
		processor: processor,
	})
	return p
}

func (p *Choice) SetOtherwise(processor api.Processor) *Choice {
	p.otherwise = processor
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
