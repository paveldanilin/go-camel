package env

type Map struct {
	vars map[string]string
}

func NewMap(vars map[string]string) *Map {
	return &Map{vars: vars}
}

func (env *Map) LookupVar(name string) (string, bool) {
	v, exists := env.vars[name]
	return v, exists
}
