package env

import "github.com/paveldanilin/go-camel/pkg/camel/api"

type Composite struct {
	envs []api.Env
}

func NewComposite(envs ...api.Env) *Composite {
	return &Composite{envs: envs}
}

func (c *Composite) LookupVar(name string) (string, bool) {
	for _, env := range c.envs {
		if val, exists := env.LookupVar(name); exists {
			return val, true
		}
	}
	return "", false
}
