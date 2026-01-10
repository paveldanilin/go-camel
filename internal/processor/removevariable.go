package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
)

type RemoveVariable struct {
	variableNames []string
}

func NewRemoveVariable(names ...string) *RemoveVariable {
	return &RemoveVariable{
		variableNames: names,
	}
}

func (p *RemoveVariable) Process(e *exchange.Exchange) {
	for _, name := range p.variableNames {
		e.RemoveVariable(name)
	}
}
