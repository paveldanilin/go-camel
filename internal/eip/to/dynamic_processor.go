package to

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/api"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/template"
	"strings"
)

type endpointRegistry interface {
	Endpoint(uri string) api.Endpoint
}

type dynamicToProcessor struct {
	routeName        string
	name             string
	env              api.Env
	uriVars          []string
	uriEnvVars       []string
	uriTpl           *template.Template
	endpointRegistry endpointRegistry
}

func NewDynamicProcessor(routeName, name string, uriTpl *template.Template, uriVars []string, endpointRegistry endpointRegistry, env api.Env) *dynamicToProcessor {
	// Build list of env vars (prefix 'env.')
	var uriEnvVars []string
	for _, varName := range uriVars {
		if strings.HasPrefix(varName, "env.") {
			uriEnvVars = append(uriEnvVars, varName)
		}
	}

	return &dynamicToProcessor{
		routeName:        routeName,
		name:             name,
		env:              env,
		uriVars:          uriVars,
		uriEnvVars:       uriEnvVars,
		uriTpl:           uriTpl,
		endpointRegistry: endpointRegistry,
	}
}

func (p *dynamicToProcessor) Name() string {
	return p.name
}

func (p *dynamicToProcessor) RouteName() string {
	return p.routeName
}

func (p *dynamicToProcessor) Process(e *exchange.Exchange) {
	// Prepare data for dynamic resolving
	// data[] = {
	//	'id':		 	exchange.Message's id
	//	'exchangeId':	exchange.Exchange's id
	//	'body':			exchange.Message's body
	//	'header':		exchange.Message's headers
	//	'error':		exchange.Exchange's error
	//	'property':		exchange.Exchange's properties
	//	'env:			map[string]string{}
	//}
	data := e.AsMap()
	// Resolve env variables
	if len(p.uriEnvVars) > 0 {
		env := make(map[string]string, len(p.uriEnvVars))
		for _, varName := range p.uriEnvVars {
			// Resolve without prefix 'env.'
			if varValue, varExists := p.env.LookupVar(strings.TrimLeft(varName, "env.")); varExists {
				env[varName] = varValue
			}

		}
		data["env"] = env
	}

	// Resolve uri variables
	uri, err := p.uriTpl.Render(data)
	if err != nil {
		e.SetError(fmt.Errorf("failed to resolve dynamic uri '%s': %w", p.uriTpl.Template(), err))
		return
	}

	// TODO: cache [URI]=Producer?
	// Resolve endpoint
	endpoint := p.endpointRegistry.Endpoint(uri)
	if endpoint == nil {
		e.SetError(fmt.Errorf("endpoint not found for uri '%s'", uri))
		return
	}

	// Create producer
	producer, err := endpoint.CreateProducer()
	if err != nil {
		e.SetError(err)
		return
	}

	producer.Process(e)
}
