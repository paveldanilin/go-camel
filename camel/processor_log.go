package camel

import (
	"github.com/paveldanilin/go-camel/template"
)

type logProcessor struct {
	id string

	// msg represents a templated message
	msg   string
	level LogLevel

	// will be used if msg contains variables like: ${var_name} or ${person.name} or ${person.products[0].id}.
	tpl    *template.Template
	logger Logger
}

func newLogProcessor(id, msg string, level LogLevel, logger Logger) *logProcessor {
	p := &logProcessor{
		id:     id,
		msg:    msg,
		level:  level,
		logger: logger,
	}
	if template.HasVars(msg) {
		t, err := template.Parse(msg)
		if err != nil {
			// failed to create processor
			panic(err)
		}
		p.tpl = t
	}
	return p
}

func (p *logProcessor) Process(exchange *Exchange) {
	if p.tpl == nil {
		p.logger.Log(exchange.ctx, p.level, p.msg)
		return
	}

	// resolve variables, render and send to the logger
	msg, err := p.tpl.Render(exchange.asMap())
	if err != nil {
		exchange.SetError(err)
		return
	}

	p.logger.Log(exchange.ctx, p.level, msg)
}
