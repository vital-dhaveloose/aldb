package aldb

import (
	"errors"
	"fmt"
	"strings"
)

type Ref interface {
	IsComplete() bool
	String() string
	Parse(s string) error
	//TODO json and yaml marshalling and parsing/scanning
	IdMap() map[string]string
	FromIdMap(ids map[string]string) bool
}

type PatternRef struct {
	Pattern     string
	IdFuncs     map[string]func() string
	FromIdsFunc func(ids map[string]string) bool
}

func (r *PatternRef) IsComplete() bool {
	if r == nil {
		return false
	}
	ids := r.IdMap()
	for _, field := range getPatternFields(r.Pattern) {
		if val, found := ids[field]; !found || len(val) == 0 {
			return false
		}
	}
	return true
}

func (r *PatternRef) String() string {
	return insertIds(r.Pattern, r.IdMap())
}

func (r *PatternRef) Parse(s string) error {
	if r == nil {
		return errors.New("cannot parse into a nil *PatternRef")
	}
	ids, err := getIdsFromString(r.Pattern, s)
	if err != nil {
		return err
	}
	r.FromIdsFunc(ids)
	return nil
}

func (r *PatternRef) IdMap() map[string]string {
	if r == nil {
		return nil
	}
	m := map[string]string{}
	for idName, idFunc := range r.IdFuncs {
		m[idName] = idFunc()
	}
	return m
}

//#region FooRef based on PatternRef

// Working with interface instead of simple struct
// Construction is

type FooRef interface {
	Ref
	FooId() string
}

type fooRefImpl struct {
	PatternRef
	fooId string
}

func CreateFooRef(fooId string) FooRef {
	r := &fooRefImpl{fooId: fooId}
	r.PatternRef = PatternRef{
		Pattern: "foos/{fooId}",
		IdFuncs: map[string]func() string{
			"fooId": r.FooId,
		},
		FromIdsFunc: r.FromIdMap,
	}
	return r
}

func (r *fooRefImpl) FooId() string {
	return r.fooId
}

func (r *fooRefImpl) FromIdMap(ids map[string]string) bool {
	found := false
	r.fooId, found = ids["fooId"]
	return found && len(r.fooId) > 0
}

//#endregion

//#region FooRef2 from scratch, using shared functions

type FooRef2 struct {
	FooId string
}

func (r *FooRef2) IsComplete() bool {
	return r != nil && len(r.FooId) > 0
}

func (r *FooRef2) String() string {
	return insertIds("foos/{fooId}", r.IdMap())
}

func (r *FooRef2) Parse(s string) error {
	if r == nil {
		return errors.New("cannot parse into a nil *FooRef2")
	}
	ids, err := getIdsFromString("foos/{fooId}", s)
	if err != nil {
		return err
	}
	r.FromIdMap(ids)
	return nil
}

func (r *FooRef2) IdMap() map[string]string {
	if !r.IsComplete() {
		return nil
	}
	return map[string]string{
		"fooId": r.FooId,
	}
}

func (r *FooRef2) FromIdMap(ids map[string]string) bool {
	found := false
	r.FooId, found = ids["fooId"]
	return found
}

//#endregion

func getIdsFromString(pattern string, s string) (map[string]string, error) {
	patternSegments := strings.Split(pattern, "/")
	segments := strings.Split(s, "/")
	if len(patternSegments) != len(segments) {
		return nil, fmt.Errorf("cannot get ids from string, number of parts don't match. pattern='%s', string='%s'", pattern, s)
	}
	out := map[string]string{}
	for i, patternPart := range patternSegments {
		if field, isField := patternSegmentToField(patternPart); isField {
			out[field] = segments[i]
		} else if patternPart != segments[i] {
			return nil, fmt.Errorf("segment %d of string doesn't match pattern. pattern='%s', string='%s'", i, pattern, s)
		}
	}
	return out, nil
}

func insertIds(pattern string, ids map[string]string) string {
	out := pattern
	for _, field := range getPatternFields(pattern) {
		id, found := ids[field]
		if !found {
			panic(fmt.Sprintf("no id found for field %s in pattern %s", field, pattern))
		}
		out = strings.Replace(out, "{"+field+"}", id, 1)
	}
	return out
}

func getPatternFields(pattern string) []string {
	segments := strings.Split(pattern, "/")
	out := make([]string, 0, len(segments))
	for _, segment := range segments {
		if field, isField := patternSegmentToField(segment); isField {
			out = append(out, field)
		}
	}
	return out
}

func patternSegmentToField(segment string) (field string, isField bool) {
	if !strings.HasPrefix(segment, "{") || !strings.HasSuffix(segment, "}") {
		return "", false
	}
	return segment[1 : len(segment)-1], true
}
