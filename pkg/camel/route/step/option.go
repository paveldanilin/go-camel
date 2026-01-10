package step

import "github.com/paveldanilin/go-camel/pkg/camel/api"

type step interface {
	setStepName(name string)
}

type Option func(step)

func WithName(name string) Option {
	return func(s step) {
		s.setStepName(name)
	}
}

func ApplyOptions(routeStep api.RouteStep, options ...Option) {
	if s, isStep := routeStep.(step); isStep {
		for _, opt := range options {
			opt(s)
		}
	}
}
