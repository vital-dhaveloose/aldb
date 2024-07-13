package aldb

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

var ErrorNullCast = errors.New("null-cast")
var ErrorWrongType = errors.New("wrong-type")

type Value interface {
	// IsNull returns whether the Value represents null.
	IsNull() bool

	CastBool() (bool, error)
	CastInt64() (int64, error)
	CastFloat64() (float64, error)
	CastString() (string, error)
	CastTime() (time.Time, error)
	CastSlice() ([]any, error)
	CastMap() (map[string]any, error)

	CastRef() (VersionRef, error)

	CastAny() any

	// Get returns the value that is at the given path relative to this value. If this
	// Value represents an error, the result is unchanged. If the path doesn't exist in
	// this Value, the result represents an error.
	//TODO specify type of error
	Get(path Path) (Value, error)

	Set(path Path, val any) (Value, error)
}

func Val(itf any) Value {
	return &itfWrappingValue{
		actual: convertRecursively(itf),
	}
}

func convertRecursively(cand any) any {
	if cand == nil {
		return nil
	}
	switch c := cand.(type) {
	// recursive
	case []any:
		res := make([]any, len(c))
		for i := range c {
			res[i] = convertRecursively(c[i])
		}
		return res
	case map[string]any:
		res := make(map[string]any, len(c))
		for k := range c {
			res[k] = convertRecursively(c[k])
		}
		return res
	// type conversions
	case int:
		return int64(c)
	case int32:
		return int64(c)
	case float32:
		return float64(c)
	case map[any]any: // for result of yaml parsing
		res := make(map[string]any, len(c))
		for k := range c {
			res[fmt.Sprintf("%v", k)] = convertRecursively(c[k])
		}
		return res
	}
	return cand
}

func Construct(paths []Path, vals []any) (out Value, err error) {
	out = &itfWrappingValue{}
	for i, pth := range paths {
		if pth.Len() == 0 {
			break
		}
		out, err = out.Set(pth, vals[i])
		if err != nil {
			return nil, err
		}
	}
	return
}

type itfWrappingValue struct {
	actual any
}

func (v *itfWrappingValue) IsNull() bool {
	return v.actual == nil
}

//#region Cast*

func (v *itfWrappingValue) CastBool() (bool, error) {
	if v.actual == nil {
		return false, ErrorNullCast
	}
	b, castOk := v.actual.(bool)
	if !castOk {
		return false, ErrorWrongType
	}
	return b, nil
}

func (v *itfWrappingValue) CastInt64() (int64, error) {
	if v.actual == nil {
		return 0, ErrorNullCast
	}
	i, castOk := v.actual.(int64)
	if !castOk {
		return 0, ErrorWrongType
	}
	return i, nil
}

func (v *itfWrappingValue) CastFloat64() (float64, error) {
	if v.actual == nil {
		return 0, ErrorNullCast
	}
	fl, castOk := v.actual.(float64)
	if !castOk {
		return 0, ErrorWrongType
	}
	return fl, nil
}

func (v *itfWrappingValue) CastString() (string, error) {
	if v.actual == nil {
		return "", ErrorNullCast
	}
	str, castOk := v.actual.(string)
	if !castOk {
		return "", ErrorWrongType
	}
	return str, nil
}

func (v *itfWrappingValue) CastTime() (time.Time, error) {
	if v.actual == nil {
		return time.Time{}, ErrorNullCast
	}
	t, castOk := v.actual.(time.Time)
	if !castOk {
		return time.Time{}, ErrorWrongType
	}
	return t, nil
}

func (v *itfWrappingValue) CastSlice() ([]any, error) {
	if v.actual == nil {
		return nil, ErrorNullCast
	}
	slc, castOk := v.actual.([]any)
	if !castOk {
		return nil, ErrorWrongType
	}
	return slc, nil
}

func (v *itfWrappingValue) CastMap() (map[string]any, error) {
	if v.actual == nil {
		return nil, ErrorNullCast
	}
	m, castOk := v.actual.(map[string]any)
	if !castOk {
		return nil, ErrorWrongType
	}
	return m, nil
}

func (v *itfWrappingValue) CastRef() (VersionRef, error) {
	panic("itfWrappingValue#CastRef not implemented")
}

func (v *itfWrappingValue) CastAny() any {
	return v.actual
}

//#endregion

func (v *itfWrappingValue) Get(path Path) (foundVal Value, err error) {
	found := false
	checkFound := func(root Value, pathFromRoot Path, val Value) (newVal Value, err error) {
		if path.Equal(pathFromRoot) {
			foundVal = val
			found = true
			return val, ErrorStop
		}
		return val, nil
	}
	_, err = TraverseDepthFirst(v, checkFound)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("NOTFOUND") //TODO better error
	}
	return
}

func (v *itfWrappingValue) Set(path Path, val any) (Value, error) {
	head, tail := path.HeadTail()
	if len(head) == 0 {
		return v, nil
	}
	if idx, isIdx := head.Index(); isIdx {
		return v.SetItem(idx, tail, val)
	}
	return v.SetEntry(head.Key(), tail, val)
}

func (v *itfWrappingValue) SetItem(idx int, pth Path, val any) (Value, error) {
	var slc []any
	var err error
	if v.IsNull() {
		slc = make([]any, idx+1)
	} else {
		slc, err = v.CastSlice()
		if err != nil {
			return nil, err
		}
	}
	if len(slc) < idx+1 {
		old := slc
		slc = make([]any, idx+1)
		copy(slc, old)
	}
	slc[idx], err = Val(slc[idx]).Set(pth, val)
	if err != nil {
		return nil, err
	}
	return Val(slc), nil
}

func (v *itfWrappingValue) SetEntry(key string, pth Path, val any) (Value, error) {
	var mep map[string]any
	var err error
	if v.IsNull() {
		mep = make(map[string]any)
	} else {
		mep, err = v.CastMap()
		if err != nil {
			return nil, err
		}
	}
	mep[key], err = Val(mep[key]).Set(pth, val)
	if err != nil {
		return nil, err
	}
	return Val(mep), nil
}

//#region Path

type Path interface {
	Len() int
	HeadTail() (Segment, Path)
	Append(segs ...Segment) Path
	Segments() []Segment
	Equal(other Path) bool
	Less(other Path) bool
	String() string
}

type Segment string

func (s Segment) Key() string {
	return string(s)
}

func (s Segment) Index() (int, bool) {
	idx, err := strconv.Atoi(string(s))
	if err == nil {
		return idx, true
	}
	return -1, false
}

func KeySegment(key string) Segment {
	return Segment(key)
}

func IndexSegment(idx int) Segment {
	return Segment(strconv.FormatInt(int64(idx), 10))
}

func EmptyPath() Path {
	panic("EmptyPath not implemented yet")
}

//#endregion

//#region Traverse

var (
	ErrorBreak = errors.New("BREAK")
	ErrorStop  = errors.New("STOP")
)

// VisitFunc is the function that will be called first when a step is made in the traversion.
// The return value must be the new value for the current location in the structure, use val to
// only inspect (not manipulate) the structure.
// If the traversion must be stopped completely, return (val, ErrorStop).
// If the traversion of the above array or dictionary must be stopped, return (val, ErrorBreak).
type VisitFunc func(root Value, pathFromRoot Path, val Value) (newVal Value, err error)

func TraverseDepthFirst(root Value, visit VisitFunc) (newVal Value, err error) {
	newVal, err = traverseDepthFirstStep(root, EmptyPath(), root, visit)
	if err != nil && !errors.Is(err, ErrorBreak) && !errors.Is(err, ErrorStop) {
		return nil, err
	}
	return newVal, nil
}

func traverseDepthFirstStep(root Value, pathFromRoot Path, val Value, visit VisitFunc) (newVal Value, err error) {
	newVal, err = visit(root, pathFromRoot, val)
	if err != nil {
		if err == ErrorStop {
			return
		}
		return nil, err
	}
	if newVal.IsNull() {
		return
	}
	switch nv := newVal.CastAny().(type) {
	case bool:
		return
	case int64:
		return
	case float64:
		return
	case string:
		return
	case time.Time:
		return
	//TODO VersionRef ?
	case []any:
		for i := range nv {
			newItemval, err2 := traverseDepthFirstStep(root, pathFromRoot.Append(IndexSegment(i)), Val(nv[i]), visit)
			if err2 != nil && !errors.Is(err2, ErrorBreak) && !errors.Is(err2, ErrorStop) {
				return nil, err2
			}
			nv[i] = newItemval.CastAny()
			if errors.Is(err2, ErrorBreak) {
				break
			}
			if errors.Is(err2, ErrorStop) {
				return newVal, err2
			}
		}
		return
	case map[string]any:
		for k := range nv {
			newEntryVal, err2 := traverseDepthFirstStep(root, pathFromRoot.Append(KeySegment(k)), Val(nv[k]), visit)
			if err2 != nil && !errors.Is(err2, ErrorBreak) && !errors.Is(err2, ErrorStop) {
				return nil, err2
			}
			nv[k] = newEntryVal.CastAny()
			if errors.Is(err2, ErrorBreak) {
				break
			}
			if errors.Is(err2, ErrorStop) {
				return newVal, err2
			}
		}
		return
	}
	return nil, fmt.Errorf("unsupported-type: %s", reflect.TypeOf(newVal).Name())
}

//#endregion
