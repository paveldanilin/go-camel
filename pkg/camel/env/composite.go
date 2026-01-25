package env

type Composite struct {
	envs []Env
}

func NewComposite(envs ...Env) *Composite {
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
