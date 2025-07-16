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

func (p *LogProcessor) Process(message *camel.Message) error {

	fmt.Printf("%s payload=%+v; headers=%+v\n", p.prefix, message.Payload(), message.MessageHeaders())

	return nil
}
