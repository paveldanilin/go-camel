package to

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/template"
)

type endpointRegistry interface {
	Endpoint(uri string) api.Endpoint
}

type toProcessor struct {
	routeName        string
	name             string
	env              api.Env
	uri              string
	uriVars          []string
	uriTpl           *template.Template
	endpointRegistry endpointRegistry
}

func NewProcessor(routeName, name, uri string, endpointRegistry endpointRegistry, env api.Env) *toProcessor {
	uriVars, _ := template.Vars(uri)
	var uriTpl *template.Template
	if len(uriVars) > 0 {
		uriTpl, _ = template.Parse(uri)
	}

	return &toProcessor{
		routeName:        routeName,
		name:             name,
		env:              env,
		uri:              uri,
		uriVars:          uriVars,
		uriTpl:           uriTpl,
		endpointRegistry: endpointRegistry,
	}
}

func (p *toProcessor) Name() string {
	return p.name
}

func (p *toProcessor) RouteName() string {
	return p.routeName
}

func (p *toProcessor) Process(e *exchange.Exchange) {
	uri := p.uri
	// Dynamic URI
	if p.uriTpl != nil {
		data := e.AsMap()
		env := make(map[string]string, len(p.uriVars))
		for _, varName := range p.uriVars {
			if varValue, varExists := p.env.LookupVar(varName); varExists {
				env[varName] = varValue
			}
		}
		data["env"] = env
		var renderErr error
		uri, renderErr = p.uriTpl.Render(data)
		if renderErr != nil {
			e.SetError(fmt.Errorf("failed to resolve dynamic uri '%s': %w", p.uri, renderErr))
			return
		}
	}

	endpoint := p.endpointRegistry.Endpoint(uri)
	if endpoint == nil {
		e.SetError(fmt.Errorf("endpoint not found for uri '%s'", uri))
		return
	}

	producer, err := endpoint.CreateProducer()
	if err != nil {
		e.SetError(err)
		return
	}

	producer.Process(e)
}
