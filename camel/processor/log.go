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

func (p *LogProcessor) Process(message *camel.Message) {

	fmt.Printf("%s body=%+v; headers=%+v\n", p.prefix, message.Body, message.MessageHeaders())
}
