package aldb

import (
	"sort"

	"golang.org/x/exp/slices"
)

//#region Sets

func SetIntersection[C comparable](in, andIn []C) []C {
	return FilterSlice(in, func(_ int, item C) bool {
		return slices.Contains(andIn, item)
	})
}

func SetDiff[C comparable](in, butNotIn []C) []C {
	return FilterSlice(in, func(_ int, item C) bool {
		return !slices.Contains(butNotIn, item)
	})
}

//#endregion

//#region Slices

func TransformSlice[I, O any](ins []I, transform func(idx int, in I) O) []O {
	if ins == nil {
		return nil
	}
	out := make([]O, len(ins))
	for i := range ins {
		out[i] = transform(i, ins[i])
	}
	return out
}

func FilterSlice[A any](in []A, keep func(idx int, item A) bool) []A {
	if in == nil {
		return nil
	}
	out := make([]A, 0, len(in))
	for i := range in {
		if keep(i, in[i]) {
			out = append(out, in[i])
		}
	}
	return out
}

func ToIndexMap[V any](in []V) map[int]V {
	if in == nil {
		return nil
	}
	out := map[int]V{}
	for i := range in {
		out[i] = in[i]
	}
	return out
}

//#endregion

//#region Maps

func Group[K comparable, V any](values []V, toKey func(v V) K) map[K][]V {
	if values == nil {
		return nil
	}
	out := map[K][]V{}
	for i := range values {
		val := values[i]
		key := toKey(val)
		if _, found := out[key]; !found {
			out[key] = []V{val}
		} else {
			out[key] = append(out[key], val)
		}
	}
	return out
}

func FilterMap[K comparable, V any](in map[K]V, keep func(key K, value V) bool) map[K]V {
	if in == nil {
		return nil
	}
	out := map[K]V{}
	for k := range in {
		if keep(k, in[k]) {
			out[k] = in[k]
		}
	}
	return out
}

func GetMapKeys[K comparable, V any](m map[K]V) []K {
	if m == nil {
		return nil
	}
	out := make([]K, len(m))
	i := 0
	for k := range m {
		out[i] = k
		i++
	}
	return out
}

func GetMapValuesSortedByKey[K comparable, V any](m map[K]V, compare CompareFunc[K]) []V {
	if m == nil {
		return nil
	}
	if compare == nil {
		compare = CompareKnownTypes[K]
	}
	out := make([]V, len(m))
	keys := make([]K, len(m))
	i := 0
	for key, val := range m {
		keys[i] = key
		out[i] = val
		i++
	}
	sort.Slice(out, func(i, j int) bool {
		less, _ := compare(keys[i], keys[j])
		return less
	})
	return out
}

//endregion

//#region Interval

type Boundary[A any] struct {
	closed   bool
	value    A
	included bool
}

func (b Boundary[A]) Value() (val A, inclusive bool, closed bool) {
	return b.value, b.included, b.closed
}

func (b Boundary[A]) IsOpen() bool {
	return !b.closed
}

func Incl[A any](value A) Boundary[A] {
	return Boundary[A]{closed: true, value: value, included: true}
}

func Excl[A any](value A) Boundary[A] {
	return Boundary[A]{closed: true, value: value, included: false}
}

func Open[A any]() Boundary[A] {
	return Boundary[A]{closed: false}
}

type Interval[A any] struct {
	Lower, Upper Boundary[A]
}

func (i Interval[A]) Contains(a A, compare CompareFunc[A]) bool {
	return PassesOver(a, i.Lower, compare) && PassesUnder(a, i.Upper, compare)
}

func PassesOver[A any](value A, boundary Boundary[A], compare CompareFunc[A]) bool {
	if boundary.IsOpen() {
		return true
	}
	if compare == nil {
		compare = CompareKnownTypes[A]
	}
	valueOverBoundary, equal := compare(boundary.value, value)
	return (equal && boundary.included) || valueOverBoundary
}

func PassesUnder[A any](value A, boundary Boundary[A], compare CompareFunc[A]) bool {
	if boundary.IsOpen() {
		return true
	}
	if compare == nil {
		compare = CompareKnownTypes[A]
	}
	valueUnderBoundary, equal := compare(value, boundary.value)
	return (equal && boundary.included) || valueUnderBoundary
}

type CompareFunc[A any] func(left, right A) (leftIsSmaller, equal bool)

func CompareKnownTypes[A any](left, right A) (leftIsSmaller bool, equal bool) {
	switch any(left).(type) {
	case string:
		l := any(left).(string)
		r := any(right).(string)
		return l < r, l == r
	case int:
		l := any(left).(int)
		r := any(right).(int)
		return l < r, l == r
	default:
		return false, false
	}
}

//#endregion
