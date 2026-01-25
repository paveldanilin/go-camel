package processor

import (
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/template"
	"log/slog"
)

type Log struct {
	level       slog.Level
	logger      *slog.Logger
	msg         string             // msg represents a templated message
	msgTemplate *template.Template // will be used if msg contains variables like: ${var_name} or ${person.headerName} or ${person.products[0].headerName}.
}

func NewLog(msg string, level slog.Level, logger *slog.Logger) (*Log, error) {
	p := &Log{
		msg:    msg,
		level:  level,
		logger: logger,
	}

	if template.HasVars(msg) {
		t, err := template.Parse(msg)
		if err != nil {
			return nil, err
		}
		p.msgTemplate = t
	}

	return p, nil
}

func (p *Log) Process(e *exchange.Exchange) {
	// Message does not have variables, nothing to interpolate
	if p.msgTemplate == nil {
		p.logger.Log(e.Context(), p.level, p.msg)
		return
	}

	// Message has variables
	// resolve variables, render and send to the logger
	msg, err := p.msgTemplate.Render(e.AsMap())
	if err != nil {
		e.SetError(err)
		return
	}

	p.logger.Log(e.Context(), p.level, msg)
}
