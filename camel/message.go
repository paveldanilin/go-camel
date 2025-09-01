package camel

import "github.com/google/uuid"

type Message struct {
	id      string
	headers Values
	Body    any
}

func NewMessage() *Message {
	return &Message{
		id:      uuid.NewString(),
		headers: newValues(),
	}
}

func (m *Message) Id() string {
	return m.id
}

func (m *Message) Headers() *Values {
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
