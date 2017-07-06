package model

import (
	"reflect"
	"time"
	"unsafe"
)

var timeType = reflect.TypeOf(time.Time{})

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

func Equal(x, y interface{}) (map[string]interface{}, bool) {
	results := make(map[string]interface{})
	// if x == nil || y == nil {
	// 	return
	// }
	a1 := reflect.ValueOf(x)
	a2 := reflect.ValueOf(y)
	// if v1.Type() != v2.Type() {
	// 	return
	// }
	var v1, v2 reflect.Value
	if a1.Kind() == reflect.Ptr {
		v1 = a1.Elem()
	} else {
		v1 = a1
	}
	if a2.Kind() == reflect.Ptr {
		v2 = a2.Elem()
	} else {
		v2 = a2
	}

	equal := true
	for i, n := 0, v1.NumField(); i < n; i++ {
		if !deepValueEqual(v1.Field(i), v2.Field(i), make(map[visit]bool)) {
			results[v1.Type().Field(i).Name] = bypassCanInterface(v2.Field(i)).Interface()
			equal = false
		}
	}

	return results, equal
}

func deepValueEqual(v1, v2 reflect.Value, visited map[visit]bool) bool {
	if !v1.IsValid() || !v2.IsValid() {
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		return false
	}

	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		addr2 := unsafe.Pointer(v2.UnsafeAddr())
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Slice:
		// We treat a nil slice the same as an empty slice.
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), visited) {
				return false
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			return v1.IsNil() == v2.IsNil()
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited)
	case reflect.Ptr:
		return deepValueEqual(v1.Elem(), v2.Elem(), visited)
	case reflect.Struct:
		if v1.Type() == timeType {
			// Special case for time - we ignore the time zone.
			t1 := interfaceOf(v1).(time.Time)
			t2 := interfaceOf(v2).(time.Time)
			return t1.Equal(t2)
		}
		for i, n := 0, v1.NumField(); i < n; i++ {
			if !deepValueEqual(v1.Field(i), v2.Field(i), visited) {
				return false
			}
		}
		return true
	case reflect.Map:
		if v1.IsNil() != v2.IsNil() {
			return false
		}
		if v1.Len() != v2.Len() {
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() || !deepValueEqual(v1.MapIndex(k), v2.MapIndex(k), visited) {
				return false
			}
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		return false
	// case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	// 	if v1.Int() != v2.Int() {
	// 		return false
	// 	}
	// 	return true
	// case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	// 	if v1.Uint() != v2.Uint() {
	// 		return false
	// 	}
	// 	return true
	// case reflect.Float32, reflect.Float64:
	// 	if v1.Float() != v2.Float() {
	// 		return false
	// 	}
	// 	return true
	// case reflect.Complex64, reflect.Complex128:
	// 	if v1.Complex() != v2.Complex() {
	// 		return false
	// 	}
	// 	return true
	// case reflect.Bool:
	// 	if v1.Bool() != v2.Bool() {
	// 		return false
	// 	}
	// 	return true
	// case reflect.String:
	// 	if v1.String() != v2.String() {
	// 		return false
	// 	}
	// 	return true
	// case reflect.Chan, reflect.UnsafePointer:
	// 	if v1.Pointer() != v2.Pointer() {
	// 		return false
	// 	}
	// 	return true
	default:
		return bypassCanInterface(v1).Interface() == bypassCanInterface(v2).Interface()
	}
}

// interfaceOf returns v.Interface() even if v.CanInterface() == false.
// This enables us to call fmt.Printf on a value even if it's derived
// from inside an unexported field.
// See https://code.google.com/p/go/issues/detail?id=8965
// for a possible future alternative to this hack.
func interfaceOf(v reflect.Value) interface{} {
	if !v.IsValid() {
		return nil
	}
	return bypassCanInterface(v).Interface()
}

type flag uintptr

var flagRO flag

// constants copied from reflect/value.go
const (
	// The value of flagRO up to and including Go 1.3.
	flagRO1p3 = 1 << 0

	// The value of flagRO from Go 1.4.
	flagRO1p4 = 1 << 5
)

var flagValOffset = func() uintptr {
	field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag")
	if !ok {
		panic("reflect.Value has no flag field")
	}
	return field.Offset
}()

func flagField(v *reflect.Value) *flag {
	return (*flag)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + flagValOffset))
}

// bypassCanInterface returns a version of v that
// bypasses the CanInterface check.
func bypassCanInterface(v reflect.Value) reflect.Value {
	if !v.IsValid() || v.CanInterface() {
		return v
	}
	*flagField(&v) &^= flagRO
	return v
}

func init() {
	field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag")
	if !ok {
		panic("reflect.Value has no flag field")
	}
	if field.Type.Kind() != reflect.TypeOf(flag(0)).Kind() {
		panic("reflect.Value flag field has changed kind")
	}
	var t struct {
		a int
		A int
	}
	vA := reflect.ValueOf(t).FieldByName("A")
	va := reflect.ValueOf(t).FieldByName("a")
	flagA := *flagField(&vA)
	flaga := *flagField(&va)

	// Infer flagRO from the difference between the flags
	// for the (otherwise identical) fields in t.
	flagRO = flagA ^ flaga
	if flagRO != flagRO1p3 && flagRO != flagRO1p4 {
		panic("reflect.Value read-only flag has changed semantics")
	}
}
