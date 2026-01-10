package exchange

const (
	PropertyEndpointUri = "CAMEL_ENDPOINT_URI"
	PropertyRouteName   = "CAMEL_ROUTE_NAME"
	PropertyStepName    = "CAMEL_ROUTE_STEP_NAME"
	PropertyLoopIndex   = "CAMEL_LOOP_INDEX"
	PropertyLoopBreak   = "CAMEL_LOOP_BREAK"
)

type Properties map[string]any

func (p Properties) Get(name string) (any, bool) {
	if v, exists := p[name]; exists {
		return v, true
	}
	return nil, false
}

func (p Properties) Set(name string, value any) {
	p[name] = value
}

func (p Properties) SetAll(kv map[string]any) {
	clear(p)

	for k, vl := range kv {
		p[k] = vl
	}
}

func (p Properties) All() map[string]any {
	return p
}

func (p Properties) Remove(name string) {
	delete(p, name)
}

func (p Properties) Has(name string) bool {
	_, exists := p[name]
	return exists
}

func (p Properties) Copy() Properties {
	if p == nil {
		return nil
	}

	cp := make(Properties, len(p))
	for k, v := range p {
		if copier, isCopier := v.(Copier); isCopier {
			cp[k] = copier.Copy()
		} else {
			cp[k] = v
		}
	}

	return cp
}
