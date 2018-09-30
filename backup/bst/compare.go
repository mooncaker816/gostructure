package bst

import (
	"reflect"
)

// Comparator 比较器
type Comparator func(a, b interface{}) int

// BasicCompare 比较大小
func BasicCompare(a, b interface{}) int {

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	if va.Kind() != vb.Kind() {
		panic("can not compare between different types")
	}
	switch va.Kind() {
	case
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		if va.Int() == vb.Int() {
			return 0
		}
		if va.Int() < vb.Int() {
			return -1
		}
		return 1
	case
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		if va.Uint() == vb.Uint() {
			return 0
		}
		if va.Uint() < vb.Uint() {
			return -1
		}
		return 1
	case
		reflect.Float32,
		reflect.Float64:
		if va.Float() == vb.Float() {
			return 0
		}
		if va.Float() < vb.Float() {
			return -1
		}
		return 1
	case
		reflect.String:
		if va.String() == vb.String() {
			return 0
		}
		if va.String() < vb.String() {
			return -1
		}
		return 1
	default:
		panic("type not support for comparing,should customize compare method")
	}
}
