package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type RemoveProperty struct {
	propertyNames []string
}

func NewRemoveProperty(names ...string) *RemoveProperty {
	return &RemoveProperty{
		propertyNames: names,
	}
}

func (p *RemoveProperty) Process(e *exchange.Exchange) {
	for _, propertyName := range p.propertyNames {
		e.RemoveProperty(propertyName)
	}
}
