package exchange

const MessageHeaderHistory = "CAMEL_MESSAGE_HEADER_HISTORY"

type MessageHeaders map[string]any

func (h MessageHeaders) Get(name string) (any, bool) {
	if v, exists := h[name]; exists {
		return v, true
	}
	return nil, false
}

func (h MessageHeaders) Set(name string, value any) {
	h[name] = value
}

func (h MessageHeaders) SetAll(kv map[string]any) {
	clear(h)

	for k, vl := range kv {
		h[k] = vl
	}
}

func (h MessageHeaders) All() map[string]any {
	return h
}

func (h MessageHeaders) Remove(name string) {
	delete(h, name)
}

func (h MessageHeaders) Has(name string) bool {
	_, exists := h[name]
	return exists
}

func (h MessageHeaders) Copy() MessageHeaders {
	if h == nil {
		return nil
	}

	cp := make(MessageHeaders, len(h))
	for k, v := range h {
		if copier, isCopier := v.(Copier); isCopier {
			cp[k] = copier.Copy()
		} else {
			cp[k] = v
		}
	}

	return cp
}
