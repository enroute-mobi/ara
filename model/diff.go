package model

// Unused for now

import (
	"errors"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

var timeType = reflect.TypeOf(time.Time{})

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

type DiffResult struct {
	DiffMap map[string]interface{}
	Equal   bool
}

func Equal(x, y interface{}) (*DiffResult, error) {
	if x == nil || y == nil {
		return nil, errors.New("use of Equal with nil value")
	}
	v1 := handlePtr(reflect.ValueOf(x))
	v2 := handlePtr(reflect.ValueOf(y))
	if v1.Type() != v2.Type() {
		return nil, errors.New("use of Equal with different type values")
	}

	result := &DiffResult{
		DiffMap: make(map[string]interface{}),
		Equal:   true,
	}

	for i, n := 0, v1.NumField(); i < n; i++ {
		// Ignore unexported fields
		if !v1.Field(i).CanInterface() {
			continue
		}
		// Ignore fields with tag `diffignore:"true"`
		if checkTag(v1.Type().Field(i)) {
			continue
		}

		if !deepValueEqual(v1.Field(i), v2.Field(i), make(map[visit]bool)) {
			result.DiffMap[v1.Type().Field(i).Name] = v2.Field(i).Interface()
			result.Equal = false
		}
	}

	return result, nil
}

func handlePtr(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		return v.Elem()
	}
	return v
}

func checkTag(v reflect.StructField) bool {
	b, _ := strconv.ParseBool(v.Tag.Get("diffignore"))
	return b
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
			t1 := v1.Interface().(time.Time)
			t2 := v2.Interface().(time.Time)
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
	default:
		return v1.Interface() == v2.Interface()
	}
}
