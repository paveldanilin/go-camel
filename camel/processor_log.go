package camel

import (
	"fmt"
)

type logProcessor struct {
	stepName string
	prefix   string
}

func newLogProcessor(prefix string) *logProcessor {
	return &logProcessor{
		stepName: fmt.Sprintf("log{prefix:%s}", prefix),
		prefix:   prefix,
	}
}

func (p *logProcessor) Process(exchange *Exchange) {
	if !exchange.On(p.stepName) {
		return
	}

	fmt.Printf("%s body=%+v; headers=%+v\n", p.prefix, exchange.Message().Body, exchange.Message().Headers())
}
