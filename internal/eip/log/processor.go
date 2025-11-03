package log

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/logger"
	"github.com/paveldanilin/go-camel/pkg/camel/template"
)

type logProcessor struct {
	routeName string
	name      string
	msg       string // msg represents a templated message
	level     logger.LogLevel
	tpl       *template.Template // will be used if msg contains variables like: ${var_name} or ${person.headerName} or ${person.products[0].headerName}.
	logger    logger.Logger
}

func NewProcessor(routeName, name, msg string, level logger.LogLevel, logger logger.Logger) *logProcessor {
	p := &logProcessor{
		routeName: routeName,
		name:      name,
		msg:       msg,
		level:     level,
		logger:    logger,
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

func (p *logProcessor) Name() string {
	return p.name
}

func (p *logProcessor) Process(e *exchange.Exchange) {
	if p.tpl == nil {
		p.logger.Log(e.Context(), p.level, p.msg)
		return
	}

	// resolve variables, render and send to the logger
	msg, err := p.tpl.Render(e.AsMap())
	if err != nil {
		e.SetError(err)
		return
	}

	p.logger.Log(e.Context(), p.level, msg, "step", p.routeName, "processor", p.name)
}
