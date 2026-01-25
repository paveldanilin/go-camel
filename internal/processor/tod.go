package processor

import (
	"fmt"
	"github.com/paveldanilin/go-camel/pkg/camel/exchange"
	"github.com/paveldanilin/go-camel/pkg/camel/template"
	"strings"
)

type Env interface {
	LookupVar(name string) (string, bool)
}

type EndpointProducerFactory interface {
	CreateEndpointProducer(uri string) (Processor, error)
}

type ToD struct {
	env                     Env
	uriVars                 []string
	uriEnvVars              []string
	uriTpl                  *template.Template
	endpointProducerFactory EndpointProducerFactory
}

func NewToD(uriTpl *template.Template, uriVars []string, endpointProducerFactory EndpointProducerFactory, env Env) *ToD {
	// Build list of env vars (prefix 'env.')
	var uriEnvVars []string
	for _, varName := range uriVars {
		if strings.HasPrefix(varName, "env.") {
			uriEnvVars = append(uriEnvVars, varName)
		}
	}

	return &ToD{
		env:                     env,
		uriVars:                 uriVars,
		uriEnvVars:              uriEnvVars,
		uriTpl:                  uriTpl,
		endpointProducerFactory: endpointProducerFactory,
	}
}

func (p *ToD) Process(e *exchange.Exchange) {
	// Prepare data for dynamic resolving
	// data[] = {
	//	'id':		 	exchange.Message's id
	//	'exchangeId':	exchange.Exchange's id
	//	'body':			exchange.Message's body
	//	'header':		exchange.Message's headers
	//	'error':		exchange.Exchange's error
	//	'property':		exchange.Exchange's properties
	//	'variable':		exchange.Exchange's variables
	//	'env:			map[string]string{}
	//}
	data := e.AsMap()
	// Resolve env variables
	if len(p.uriEnvVars) > 0 {
		env := make(map[string]string, len(p.uriEnvVars))
		for _, varName := range p.uriEnvVars {
			// Resolve without prefix 'env.'
			normalizedEnvName := strings.TrimLeft(varName, "env.")
			if varValue, varExists := p.env.LookupVar(normalizedEnvName); varExists {
				env[normalizedEnvName] = varValue
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

	producer, err := p.endpointProducerFactory.CreateEndpointProducer(uri)
	if err != nil {
		e.SetError(err)
		return
	}

	producer.Process(e)
}
