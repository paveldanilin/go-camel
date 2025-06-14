package camel

type ContextProvider interface {
	Component(componentId string) Component
	Endpoint(uri string) Endpoint
	Route(routeId string) *Route
}

type Message struct {
	context ContextProvider
	payload any
	headers map[string]any
}

func NewMessage() *Message {
	return &Message{
		payload: nil,
		headers: map[string]any{},
	}
}

func NewMessageWithContext(context ContextProvider) *Message {
	return &Message{
		context: context,
		payload: nil,
		headers: map[string]any{},
	}
}

func (m *Message) Context() ContextProvider {

	return m.context
}

func (m *Message) Payload() any {

	return m.payload
}

func (m *Message) SetPayload(payload any) {

	m.payload = payload
}

func (m *Message) SetHeader(name string, value any) {

	m.headers[name] = value
}

func (m *Message) Header(name string) (any, bool) {

	if v, exists := m.headers[name]; exists {
		return v, true
	}

	return nil, false
}

func (m *Message) Headers() map[string]any {

	return m.headers
}
