package removeproperty

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type removeProperty struct {
	routeName     string
	name          string
	propertyNames []string
}

func NewProcessor(routeName, name string, propertyNames ...string) *removeProperty {
	return &removeProperty{
		routeName:     routeName,
		name:          name,
		propertyNames: propertyNames,
	}
}

func (p *removeProperty) Name() string {
	return p.name
}

func (p *removeProperty) RouteName() string {
	return p.routeName
}

func (p *removeProperty) Process(e *exchange.Exchange) {
	for _, propertyName := range p.propertyNames {
		e.RemoveProperty(propertyName)
	}
}
