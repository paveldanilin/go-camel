package exchange

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"io"
	"reflect"
	"strings"
	"time"
)

const CamelHeaderMessageHistory = "CAMEL_HEADER_MESSAGE_HISTORY"

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

func (m *Message) RemoveHeader(name string) {
	m.headers.Remove(name)
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

// Copier allows user types provide its own implementation for copying valueExpression
type Copier interface {
	Copy() any
}

// copyValue tries to copy the given valueExpression, used for Message.Body copy
func copyValue(v any) any {
	if v == nil {
		return nil
	}

	if c, ok := v.(Copier); ok {
		return c.Copy()
	}

	switch x := v.(type) {
	case bool, int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64,
		complex64, complex128,
		string:
		return x // immutable

	case time.Time:
		return x // by valueExpression

	case error:
		return x

	case []byte:
		if x == nil {
			return []byte(nil)
		}
		cp := make([]byte, len(x))
		copy(cp, x)
		return cp

	case *bytes.Buffer:
		if x == nil {
			return (*bytes.Buffer)(nil)
		}
		b := x.Bytes()
		cp := make([]byte, len(b))
		copy(cp, b)
		return bytes.NewBuffer(cp)

	case *bytes.Reader:
		if x == nil {
			return (*bytes.Reader)(nil)
		}
		all := make([]byte, x.Size())
		n, _ := x.ReadAt(all, 0)
		all = all[:n]
		cp := make([]byte, len(all))
		copy(cp, all)
		// Unable to restore position, will start from 0
		return bytes.NewReader(cp)

	case *strings.Builder:
		if x == nil {
			return (*strings.Builder)(nil)
		}
		s := x.String()
		var b strings.Builder
		b.Grow(len(s))
		b.WriteString(s)
		return &b

	case Map:
		return x.Copy()

	case map[string]any:
		if x == nil {
			return map[string]any(nil)
		}
		cp := make(map[string]any, len(x))
		for k, v2 := range x {
			cp[k] = copyValue(v2)
		}
		return cp

	case map[string]string:
		if x == nil {
			return map[string]string(nil)
		}
		cp := make(map[string]string, len(x))
		for k, v2 := range x {
			cp[k] = v2
		}
		return cp

	case []string:
		if x == nil {
			return []string(nil)
		}
		cp := make([]string, len(x))
		copy(cp, x)
		return cp

	case []int:
		if x == nil {
			return []int(nil)
		}
		cp := make([]int, len(x))
		copy(cp, x)
		return cp

	case []any:
		if x == nil {
			return []any(nil)
		}
		cp := make([]any, len(x))
		for i := range x {
			cp[i] = copyValue(x[i])
		}
		return cp

	// Unable to copy
	case io.Reader, io.ReadCloser:
		return x

	// Context is immutable
	case context.Context:
		return x
	}

	// Fallback: reflection slices/maps/pointers
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return v
		}
		elem := rv.Elem()
		cp := reflect.New(elem.Type())
		cp.Elem().Set(reflect.ValueOf(copyValue(elem.Interface())))
		return cp.Interface()

	case reflect.Slice:
		if rv.IsNil() {
			return v
		}
		l := rv.Len()
		cp := reflect.MakeSlice(rt, l, l)
		for i := 0; i < l; i++ {
			cp.Index(i).Set(reflect.ValueOf(copyValue(rv.Index(i).Interface())))
		}
		return cp.Interface()

	case reflect.Array:
		l := rv.Len()
		cp := reflect.New(rt).Elem()
		for i := 0; i < l; i++ {
			cp.Index(i).Set(reflect.ValueOf(copyValue(rv.Index(i).Interface())))
		}
		return cp.Interface()

	case reflect.Map:
		if rv.IsNil() {
			return v
		}
		cp := reflect.MakeMapWithSize(rt, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key()
			val := iter.Value()
			cp.SetMapIndex(k, reflect.ValueOf(copyValue(val.Interface())))
		}
		return cp.Interface()

	case reflect.Struct:
		cp := reflect.New(rt).Elem()
		for i := 0; i < rt.NumField(); i++ {
			if !cp.Field(i).CanSet() {
				continue
			}
			fv := rv.Field(i)
			switch fv.Kind() {
			case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
				reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
				reflect.String:
				cp.Field(i).Set(fv)
			default:
				cp.Field(i).Set(reflect.ValueOf(copyValue(fv.Interface())))
			}
		}
		return cp.Interface()

	case reflect.Interface:
		if rv.IsNil() {
			return v
		}
		inner := rv.Elem().Interface()
		return copyValue(inner)

	default:
		// Other types: fn, chan, unsafe.Pointer
		return v
	}
}
