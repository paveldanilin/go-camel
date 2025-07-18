package camel

type RuntimeProvider interface {
	Component(componentId string) Component
	Endpoint(uri string) Endpoint
	Route(routeId string) *Route
}

type MessageHeaders struct {
	headers map[string]any
}

func (h *MessageHeaders) Get(name string) (any, bool) {

	if v, exists := h.headers[name]; exists {
		return v, true
	}

	return nil, false
}

func (h *MessageHeaders) Set(name string, value any) {

	h.headers[name] = value
}

func (m *MessageHeaders) SetAll(headers map[string]any) {

	clear(m.headers)

	for k, v := range headers {
		m.headers[k] = v
	}
}

func (h *MessageHeaders) Headers() map[string]any {

	return h.headers
}

type Message struct {
	id      string
	runtime RuntimeProvider
	headers MessageHeaders
	payload any
	err     error
}

func NewMessage() *Message {
	return &Message{
		payload: nil,
		headers: MessageHeaders{
			headers: map[string]any{},
		},
	}
}

func NewMessageWithRuntime(runtime RuntimeProvider) *Message {
	return &Message{
		runtime: runtime,
		payload: nil,
		headers: MessageHeaders{
			headers: map[string]any{},
		},
	}
}

func (m *Message) ID() string {

	return m.id
}

func (m *Message) Runtime() RuntimeProvider {

	return m.runtime
}

func (m *Message) Payload() any {

	return m.payload
}

func (m *Message) SetPayload(payload any) {

	m.payload = payload
}

func (m *Message) MessageHeaders() *MessageHeaders {

	return &m.headers
}

func (m *Message) SetHeader(name string, value any) {

	m.headers.Set(name, value)
}

func (m *Message) Header(name string) (any, bool) {

	return m.headers.Get(name)
}

func (m *Message) HasHeader(name string) bool {

	_, exists := m.Header(name)
	return exists
}

func (m *Message) MustHeader(name string) any {

	if v, exists := m.headers.Get(name); exists {
		return v
	}

	panic("message header not found: '" + name + "'")
}

func (m *Message) Error() error {

	return m.err
}

func (m *Message) SetError(err error) {

	m.err = err
}

func (m *Message) IsError() bool {

	return m.err != nil
}
