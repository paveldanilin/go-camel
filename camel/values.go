package camel

type Values struct {
	values map[string]any
}

func newValues() Values {
	return Values{
		values: map[string]any{},
	}
}

func (v *Values) Get(name string) (any, bool) {
	if v, exists := v.values[name]; exists {
		return v, true
	}
	return nil, false
}

func (v *Values) Set(name string, value any) {
	v.values[name] = value
}

func (v *Values) SetAll(values map[string]any) {
	clear(v.values)

	for k, vl := range values {
		v.values[k] = vl
	}
}

func (v *Values) All() map[string]any {
	return v.values
}
