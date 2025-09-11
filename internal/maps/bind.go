package maps

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// Bind binds src (map[string]any) into dest which must be a pointer to struct.
// It recursively assigns values handling nested structs, slices, arrays, maps and pointers.
// Field matching: `json` tag (if present) then case-insensitive field name.
func Bind(src map[string]any, dest any) error {
	if dest == nil {
		return errors.New("dest is nil")
	}

	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("dest must be a non-nil pointer to a struct")
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return errors.New("dest must point to a struct")
	}

	typeInfo := buildStructFieldMap(rv.Type())

	for k, v := range src {
		if fi, ok := typeInfo[k]; ok {
			fv := rv.Field(fi.Index)
			if !fv.CanSet() {
				// unexported field
				continue
			}

			err := setValue(fv, v)
			if err != nil {
				return fmt.Errorf("field %s: %w", fi.Name, err)
			}
		}
	}

	return nil
}

type fieldInfo struct {
	Name  string
	Index int
	Tag   string
}

// buildStructFieldMap creates a lookup for "keys" to fields using json tag then case-insensitive name.
func buildStructFieldMap(t reflect.Type) map[string]fieldInfo {
	out := map[string]fieldInfo{}

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		// skip unexported fields
		if sf.PkgPath != "" {
			continue
		}

		jsonTag := sf.Tag.Get("json")
		name := sf.Name

		key := strings.ToLower(name)
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" && parts[0] != "-" {
				out[parts[0]] = fieldInfo{Name: sf.Name, Index: i, Tag: jsonTag}
			}
		}
		// fallback by lowercased field name if not already present
		if _, exists := out[key]; !exists {
			out[key] = fieldInfo{Name: sf.Name, Index: i, Tag: ""}
		}
	}

	return out
}

//nolint:cyclop
func setValue(dst reflect.Value, v any) error {
	// handle pointer destination by allocating if nil
	for dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(dst.Type().Elem()))
		}

		dst = dst.Elem()
	}

	if !dst.CanSet() {
		return errors.New("destination not settable")
	}

	if v == nil {
		// zero value
		dst.Set(reflect.Zero(dst.Type()))

		return nil
	}

	srcVal := reflect.ValueOf(v)

	switch dst.Kind() {
	case reflect.Struct:
		// if src is map[string]any -> recurse
		if m, ok := v.(map[string]any); ok {
			// create an addressable struct value to pass into Bind-like logic
			// we'll iterate fields manually instead of calling Bind to avoid type checks
			t := dst.Type()

			fieldMap := buildStructFieldMap(t)
			for key, val := range m {
				// try tag key then lowercased name
				if fi, ok := fieldMap[key]; ok {
					fv := dst.Field(fi.Index)
					if !fv.CanSet() {
						continue
					}

					err := setValue(fv, val)
					if err != nil {
						return fmt.Errorf("struct field %s: %w", fi.Name, err)
					}
				} else if fi, ok := fieldMap[strings.ToLower(key)]; ok {
					fv := dst.Field(fi.Index)
					if !fv.CanSet() {
						continue
					}

					err := setValue(fv, val)
					if err != nil {
						return fmt.Errorf("struct field %s: %w", fi.Name, err)
					}
				}
			}

			return nil
		}
		// if src is a struct assignable
		if srcVal.Type().AssignableTo(dst.Type()) {
			dst.Set(srcVal)

			return nil
		}
		// else try to convert if possible
		if srcVal.Type().ConvertibleTo(dst.Type()) {
			dst.Set(srcVal.Convert(dst.Type()))

			return nil
		}

		return fmt.Errorf("cannot set struct from %T", v)

	case reflect.Map:
		// src must be map[string]any (or map[<key>]<value> convertible)
		if m, ok := v.(map[string]any); ok {
			newMap := reflect.MakeMapWithSize(dst.Type(), len(m))
			keyType := dst.Type().Key()

			elemType := dst.Type().Elem()
			for mk, mv := range m {
				kv := reflect.New(keyType).Elem()
				// set key (try convert string to key)
				if err := setSimpleValueFromString(kv, mk); err != nil {
					// try direct string assign if key is string
					if keyType.Kind() == reflect.String {
						kv.SetString(mk)
					} else {
						return fmt.Errorf("map key conversion error: %w", err)
					}
				}

				ev := reflect.New(elemType).Elem()
				if err := setValue(ev, mv); err != nil {
					return fmt.Errorf("map value for key %s: %w", mk, err)
				}

				newMap.SetMapIndex(kv, ev)
			}

			dst.Set(newMap)

			return nil
		}
		// if src is a map with reflected type that can be converted
		if srcVal.Type().AssignableTo(dst.Type()) {
			dst.Set(srcVal)

			return nil
		}

		return fmt.Errorf("cannot set map from %T", v)

	case reflect.Slice:
		// expect src to be []any or something convertible
		if arr, ok := v.([]any); ok {
			slice := reflect.MakeSlice(dst.Type(), len(arr), len(arr))
			for i := range arr {
				err := setValue(slice.Index(i), arr[i])
				if err != nil {
					return fmt.Errorf("slice index %d: %w", i, err)
				}
			}

			dst.Set(slice)

			return nil
		}
		// if src is slice/array assignable/convertible
		if srcVal.Kind() == reflect.Slice || srcVal.Kind() == reflect.Array {
			// try direct conversion if types align
			if srcVal.Type().AssignableTo(dst.Type()) {
				dst.Set(srcVal)

				return nil
			}
			// fallback: iterate elements
			l := srcVal.Len()

			slice := reflect.MakeSlice(dst.Type(), l, l)
			for i := range l {
				elem := srcVal.Index(i).Interface()

				err := setValue(slice.Index(i), elem)
				if err != nil {
					return fmt.Errorf("slice element %d: %w", i, err)
				}
			}

			dst.Set(slice)

			return nil
		}

		return fmt.Errorf("cannot set slice from %T", v)

	case reflect.Array:
		// handle arrays similarly but must match length
		if arr, ok := v.([]any); ok {
			if len(arr) != dst.Len() {
				return fmt.Errorf("array length mismatch: dest %d src %d", dst.Len(), len(arr))
			}

			for i := 0; i < dst.Len(); i++ {
				err := setValue(dst.Index(i), arr[i])
				if err != nil {
					return fmt.Errorf("array index %d: %w", i, err)
				}
			}

			return nil
		}
		// try assignable
		if srcVal.Type().AssignableTo(dst.Type()) {
			dst.Set(srcVal)

			return nil
		}

		return fmt.Errorf("cannot set array from %T", v)

	case reflect.Interface:
		// put raw value into interface if assignable
		if srcVal.Type().AssignableTo(dst.Type()) || dst.Type().NumMethod() == 0 {
			dst.Set(srcVal)

			return nil
		}
		// create a value of the interface's concrete type if possible
		dst.Set(srcVal)

		return nil

	default:
		// basic kinds: Bool, Int*, Uint*, Float*, String
		return setBasicKind(dst, v)
	}
}

//nolint:cyclop
func setBasicKind(dst reflect.Value, v any) error {
	switch dst.Kind() {
	case reflect.Bool:
		b, err := toBool(v)
		if err != nil {
			return err
		}

		dst.SetBool(b)

		return nil
	case reflect.String:
		s, err := toString(v)
		if err != nil {
			return err
		}

		dst.SetString(s)

		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := toInt64(v)
		if err != nil {
			return err
		}

		if !withinIntRange(i, dst.Type().Bits()) {
			return fmt.Errorf("integer %d overflows %s", i, dst.Type().Kind().String())
		}

		dst.SetInt(i)

		return nil
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		u, err := toUint64(v)
		if err != nil {
			return err
		}

		if !withinUintRange(u, dst.Type().Bits()) {
			return fmt.Errorf("unsigned integer %d overflows %s", u, dst.Type().Kind().String())
		}

		dst.SetUint(u)

		return nil
	case reflect.Float32, reflect.Float64:
		f, err := toFloat64(v)
		if err != nil {
			return err
		}

		dst.SetFloat(f)

		return nil
	default:
		return fmt.Errorf("unsupported kind %s for value %T", dst.Kind().String(), v)
	}
}

// helpers for conversions.
func toBool(v any) (bool, error) {
	switch x := v.(type) {
	case bool:
		return x, nil
	case string:
		return strconv.ParseBool(x)
	case float64:
		return x != 0, nil
	case float32:
		return x != 0, nil
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(x).Int() != 0, nil
	case uint, uint8, uint16, uint32, uint64:
		return reflect.ValueOf(x).Uint() != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

func toString(v any) (string, error) {
	switch x := v.(type) {
	case string:
		return x, nil
	case []byte:
		return string(x), nil
	default:
		// fallback to fmt
		return fmt.Sprintf("%v", v), nil
	}
}

//nolint:cyclop
func toInt64(v any) (int64, error) {
	switch x := v.(type) {
	case int:
		return int64(x), nil
	case int8:
		return int64(x), nil
	case int16:
		return int64(x), nil
	case int32:
		return int64(x), nil
	case int64:
		return x, nil
	case uint:
		return int64(x), nil
	case uint8:
		return int64(x), nil
	case uint16:
		return int64(x), nil
	case uint32:
		return int64(x), nil
	case uint64:
		if x > math.MaxInt64 {
			return 0, fmt.Errorf("uint64 %d overflows int64", x)
		}

		return int64(x), nil
	case float64:
		return int64(x), nil
	case float32:
		return int64(x), nil
	case string:
		return strconv.ParseInt(x, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

//nolint:cyclop
func toUint64(v any) (uint64, error) {
	switch x := v.(type) {
	case uint:
		return uint64(x), nil
	case uint8:
		return uint64(x), nil
	case uint16:
		return uint64(x), nil
	case uint32:
		return uint64(x), nil
	case uint64:
		return x, nil
	case int:
		if x < 0 {
			return 0, fmt.Errorf("negative int %d cannot convert to uint64", x)
		}

		return uint64(x), nil
	case int8:
		if x < 0 {
			return 0, fmt.Errorf("negative int8 %d cannot convert to uint64", x)
		}

		return uint64(x), nil
	case int16:
		if x < 0 {
			return 0, fmt.Errorf("negative int16 %d cannot convert to uint64", x)
		}

		return uint64(x), nil
	case int32:
		if x < 0 {
			return 0, fmt.Errorf("negative int32 %d cannot convert to uint64", x)
		}

		return uint64(x), nil
	case int64:
		if x < 0 {
			return 0, fmt.Errorf("negative int64 %d cannot convert to uint64", x)
		}

		return uint64(x), nil
	case float64:
		if x < 0 {
			return 0, fmt.Errorf("negative float %f cannot convert to uint64", x)
		}

		return uint64(x), nil
	case string:
		return strconv.ParseUint(x, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to uint64", v)
	}
}

//nolint:cyclop
func toFloat64(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case int8:
		return float64(x), nil
	case int16:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case uint:
		return float64(x), nil
	case uint8:
		return float64(x), nil
	case uint16:
		return float64(x), nil
	case uint32:
		return float64(x), nil
	case uint64:
		return float64(x), nil
	case string:
		return strconv.ParseFloat(x, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

func withinIntRange(i int64, bits int) bool {
	switch bits {
	case 8:
		return i >= math.MinInt8 && i <= math.MaxInt8
	case 16:
		return i >= math.MinInt16 && i <= math.MaxInt16
	case 32:
		return i >= math.MinInt32 && i <= math.MaxInt32
	case 64:
		return true
	default:
		return true
	}
}

func withinUintRange(u uint64, bits int) bool {
	switch bits {
	case 8:
		return u <= math.MaxUint8
	case 16:
		return u <= math.MaxUint16
	case 32:
		return u <= math.MaxUint32
	case 64:
		return true
	default:
		return true
	}
}

// setSimpleValueFromString tries to set a reflect.Value from a string (used for map keys).
func setSimpleValueFromString(dst reflect.Value, s string) error {
	switch dst.Kind() {
	case reflect.String:
		dst.SetString(s)

		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, dst.Type().Bits())
		if err != nil {
			return err
		}

		dst.SetInt(i)

		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, dst.Type().Bits())
		if err != nil {
			return err
		}

		dst.SetUint(u)

		return nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}

		dst.SetBool(b)

		return nil
	default:
		return fmt.Errorf("unsupported key type %s", dst.Kind().String())
	}
}
