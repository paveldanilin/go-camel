package camel

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// TODO: interface Context

type Context struct {
	// TODO: components
	routes     map[string]*Route
	endpoints  map[string]Endpoint
	components map[string]Component
	consumers  []Consumer
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewContext() *Context {
	ctx, cancel := context.WithCancel(context.Background())
	return &Context{
		routes:     map[string]*Route{},
		endpoints:  map[string]Endpoint{},
		components: map[string]Component{},
		consumers:  []Consumer{},
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (ctx *Context) RegisterComponent(c Component) error {

	if _, exists := ctx.components[c.Id()]; exists {
		return errors.New("component already registered: " + c.Id())
	}

	if contextAware, isContextAware := c.(ContextAware); isContextAware {
		contextAware.SetContext(ctx)
	}

	ctx.components[c.Id()] = c

	return nil
}

func (ctx *Context) Component(componentId string) Component {

	if c, exists := ctx.components[componentId]; exists {
		return c
	}

	return nil
}

func (ctx *Context) RegisterRoute(r *Route) error {

	if _, exists := ctx.routes[r.id]; exists {
		return errors.New("route already registered: " + r.id)
	}

	ctx.routes[r.id] = r

	return nil
}

func (ctx *Context) Endpoint(uri string) Endpoint {
	return ctx.endpoints[uri]
}

func (ctx *Context) Send(uri string, payload any, headers map[string]any) (*Message, error) {

	if endpoint, exists := ctx.endpoints[uri]; exists {
		if messenger, isMessenger := endpoint.(Messenger); isMessenger {
			// TODO: message pooling?
			message := NewMessage()
			message.context = ctx
			message.payload = payload
			message.headers = headers
			err := messenger.SendMessage(message)
			if err != nil {
				return nil, err
			}
			return message, nil
		}
		return nil, errors.New("endpoint is not supporting messaging")
	}

	return nil, errors.New("endoint not found for uri: " + uri)
}

func (ctx *Context) Route(routeId string) *Route {

	if r, exists := ctx.routes[routeId]; exists {
		return r
	}

	return nil
}

func parseURI(uri string) []string {
	// "direct:foo" -> ["direct", "foo"]
	for i := 0; i < len(uri); i++ {
		if uri[i] == ':' {
			return []string{uri[:i], uri[i+1:]}
		}
	}
	return []string{uri}
}

func (ctx *Context) Start() error {

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	for _, route := range ctx.routes {

		parts := parseURI(route.from)
		if len(parts) != 2 {
			return fmt.Errorf("invalid URI format: %s", route.from)
		}

		componentName, uri := parts[0], parts[1]
		component, exists := ctx.components[componentName]
		if !exists {
			return fmt.Errorf("component not found: %s", componentName)
		}

		if endpoint, exists := ctx.endpoints[route.from]; exists {
			_, err := endpoint.Consumer(route.producer)
			if err != nil {
				return fmt.Errorf("zzz: %w", err)
			}
		} else {
			endpoint, err := component.Endpoint(uri)
			if err != nil {
				return fmt.Errorf("failed to create endpoint: %w", err)
			}
			ctx.endpoints[route.from] = endpoint
			consumer, err := endpoint.Consumer(route.producer)
			if err != nil {
				return fmt.Errorf("zzz: %w", err)
			}
			ctx.consumers = append(ctx.consumers, consumer)
		}
	}

	for _, consumer := range ctx.consumers {
		if err := consumer.Start(); err != nil {
			return err
		}
	}

	return nil

}

func (ctx *Context) Stop() error {

	ctx.cancel()
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	for _, consumer := range ctx.consumers {
		if err := consumer.Stop(); err != nil {
			return err
		}
	}

	ctx.consumers = nil
	ctx.endpoints = nil
	ctx.components = nil
	ctx.routes = nil

	return nil
}
