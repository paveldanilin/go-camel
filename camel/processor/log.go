package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/camel"
)

type LogProcessor struct {
	prefix string
}

func LogMessage(prefix string) *LogProcessor {
	return &LogProcessor{
		prefix: prefix,
	}
}

func (p *LogProcessor) Process(exchange *camel.Exchange) {
	if err := exchange.CheckCancelOrTimeout(); err != nil {
		exchange.Error = err
		return
	}

	fmt.Printf("%s body=%+v; headers=%+v\n", p.prefix, exchange.Message().Body, exchange.Message().Headers())
}
