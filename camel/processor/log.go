package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

type LogProcessor struct {
	stepName string
	prefix   string
}

func LogMessage(prefix string) *LogProcessor {
	return &LogProcessor{
		stepName: fmt.Sprintf("log{prefix:%s}", prefix),
		prefix:   prefix,
	}
}

func (p *LogProcessor) Process(exchange *camel.Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	fmt.Printf("%s body=%+v; headers=%+v\n", p.prefix, exchange.Message().Body, exchange.Message().Headers())
}
