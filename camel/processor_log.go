package camel

import (
	"fmt"
)

type logProcessor struct {
	id     string
	prefix string
}

func newLogProcessor(id, prefix string) *logProcessor {
	return &logProcessor{
		id:     id,
		prefix: prefix,
	}
}

func (p *logProcessor) Process(exchange *Exchange) {
	fmt.Printf("%s body=%+v; headers=%+v\n", p.prefix, exchange.Message().Body, exchange.Message().Headers())
}
