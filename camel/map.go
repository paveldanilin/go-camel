package camel

type Map map[string]any

func newMap() Map {
	return Map{}
}

func (m Map) Get(name string) (any, bool) {
	if v, exists := m[name]; exists {
		return v, true
	}
	return nil, false
}

func (m Map) Set(name string, value any) {
	m[name] = value
}

func (m Map) SetAll(kv map[string]any) {
	clear(m)

	for k, vl := range kv {
		m[k] = vl
	}
}

func (m Map) All() map[string]any {
	return m
}

func (m Map) Copy() Map {
	if m == nil {
		return nil
	}

	cp := make(Map, len(m))
	for k, v := range m {
		if copier, isCopier := v.(Copier); isCopier {
			cp[k] = copier.Copy()
		} else {
			cp[k] = v
		}
	}

	return cp
}
