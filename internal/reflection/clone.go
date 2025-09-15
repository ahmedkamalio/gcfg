// Package reflection provides utilities for working with Go's reflection system.
// It includes functions for deep cloning of values, handling nil checks, and other
// reflection-based operations.
package reflection

import "reflect"

// Clone creates and returns a deep copy of the input value.
func Clone[T any](v T) T {
	if IsNil(v) {
		var zero T

		return zero
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		elem := Clone(rv.Elem().Interface())
		newPtr := reflect.New(rv.Type().Elem())
		newPtr.Elem().Set(reflect.ValueOf(elem))

		return newPtr.Interface().(T)

	case reflect.Slice:
		newSlice := reflect.MakeSlice(rv.Type(), rv.Len(), rv.Cap())
		for i := range rv.Len() {
			newSlice.Index(i).Set(reflect.ValueOf(Clone(rv.Index(i).Interface())))
		}

		return newSlice.Interface().(T)

	case reflect.Map:
		newMap := reflect.MakeMapWithSize(rv.Type(), rv.Len())
		for _, key := range rv.MapKeys() {
			newMap.SetMapIndex(
				reflect.ValueOf(Clone(key.Interface())),
				reflect.ValueOf(Clone(rv.MapIndex(key).Interface())),
			)
		}

		return newMap.Interface().(T)

	default:
		return v
	}
}
