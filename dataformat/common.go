package dataformat

import "reflect"

func newInstanceOfType(t any) any {
	typ := reflect.TypeOf(t)
	val := reflect.ValueOf(t)

	var target any
	if typ.Kind() != reflect.Ptr {
		target = reflect.New(typ).Interface()
	} else {
		// if pointer, but nil, initialize
		if val.IsNil() {
			target = reflect.New(typ.Elem()).Interface()
		} else {
			target = t
		}
	}

	return target
}
