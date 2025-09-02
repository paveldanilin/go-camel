package camel

import (
	"github.com/google/uuid"
)

type Message struct {
	id      string
	headers Map
	Body    any
}

func NewMessage() *Message {
	return &Message{
		id:      uuid.NewString(),
		headers: newMap(),
	}
}

func (m *Message) Id() string {
	return m.id
}

func (m *Message) Headers() *Map {
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
	panic("camel: message header not found: '" + name + "'")
}

func (m *Message) Copy() *Message {
	if m == nil {
		return nil
	}

	var headersCopy Map
	if m.headers != nil {
		headersCopy = m.headers.Copy()
	}

	return &Message{
		id:      uuid.NewString(),
		headers: headersCopy,
		Body:    copyValue(m.Body),
	}
}
