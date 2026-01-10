package exchange

type Variables map[string]any

func (v Variables) Get(name string) (any, bool) {
	if v, exists := v[name]; exists {
		return v, true
	}
	return nil, false
}

func (v Variables) Set(name string, value any) {
	v[name] = value
}

func (v Variables) SetAll(kv map[string]any) {
	clear(v)

	for k, vl := range kv {
		v[k] = vl
	}
}

func (v Variables) All() map[string]any {
	return v
}

func (v Variables) Remove(name string) {
	delete(v, name)
}

func (v Variables) Has(name string) bool {
	_, exists := v[name]
	return exists
}

func (v Variables) Copy() Variables {
	if v == nil {
		return nil
	}

	cp := make(Variables, len(v))
	for k, vv := range v {
		if copier, isCopier := vv.(Copier); isCopier {
			cp[k] = copier.Copy()
		} else {
			cp[k] = vv
		}
	}

	return cp
}
