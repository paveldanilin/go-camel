package camel

import "github.com/paveldanilin/go-camel/pkg/camel/route"

func NewRoute(name, from string) *route.Builder {
	return route.NewBuilder(name, from)
}
