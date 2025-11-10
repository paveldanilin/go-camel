package env

type MapEnv struct {
	vars map[string]string
}

func NewMapEnv(vars map[string]string) *MapEnv {
	return &MapEnv{vars: vars}
}

func (env *MapEnv) LookupVar(name string) (string, bool) {
	v, exists := env.vars[name]
	return v, exists
}
