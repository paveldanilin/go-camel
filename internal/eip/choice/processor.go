package choice

import (
	"fmt"
	"github.com/paveldanilin/go-camel/internal/expression"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

// choiceWhen represents a single 'when' check of choiceProcessor
type choiceWhen struct {
	predicate expression.Predicate
	processor api.Processor
}

func (when *choiceWhen) match(e *exchange.Exchange) (bool, error) {
	return when.predicate.Test(e)
}

type choiceProcessor struct {
	routeName string
	name      string
	cases     []*choiceWhen
	otherwise api.Processor
}

func NewProcessor(routeName, name string) *choiceProcessor {
	return &choiceProcessor{
		routeName: routeName,
		name:      name,
		cases:     []*choiceWhen{},
	}
}

func (p *choiceProcessor) RouteName() string {
	return p.routeName
}

func (p *choiceProcessor) Name() string {
	return p.name
}

func (p *choiceProcessor) AddWhen(predicate expression.Expression, processor api.Processor) *choiceProcessor {
	if predicate == nil {
		panic(fmt.Errorf("choice: when: predicate must be not nil"))
	}
	if processor == nil {
		panic(fmt.Errorf("choice: when: processor must be not nil"))
	}
	p.cases = append(p.cases, &choiceWhen{
		predicate: expression.NewPredicateFromExpression(predicate),
		processor: processor,
	})
	return p
}

func (p *choiceProcessor) SetOtherwise(processor api.Processor) *choiceProcessor {
	p.otherwise = processor
	return p
}

func (p *choiceProcessor) Process(e *exchange.Exchange) {
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
