package aldb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type Selector[C comparable] interface {
	SelectFrom(cs []C) []C
}

// ByIndexSelector is a Selector[string] that selects items based on their index in some ordered list.
// TODO specify index range and meaning of values
type ByIndexSelector[C comparable] interface {
	Selector[C]
	SelectMapFrom(ids []C) map[int]C
}

type serialiser interface {
	json.Marshaler
	json.Unmarshaler
	yaml.Marshaler
	yaml.Unmarshaler
	String() string
	//TODO FromString()
}

//#region AllSelector

// AllSelector is a Selector that will select all inputs. Its string representation is "*".
type AllSelector[C comparable] struct{}

var _ ByIndexSelector[string] = &AllSelector[string]{}
var _ serialiser = &AllSelector[string]{}

func (s *AllSelector[C]) SelectFrom(cs []C) []C {
	if cs == nil {
		return nil
	}
	out := make([]C, len(cs))
	copy(out, cs)
	return out
}

func (s *AllSelector[C]) SelectMapFrom(values []C) map[int]C {
	return ToIndexMap(values)
}

func (s *AllSelector[C]) String() string {
	return "*"
}

func (s *AllSelector[C]) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *AllSelector[C]) UnmarshalJSON(bts []byte) error {
	if string(bts) != `"*"` {
		return fmt.Errorf(`only "*" can be unmarshalled to AllSelector, not %s`, string(bts))
	}
	return nil
}

func (s *AllSelector[C]) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *AllSelector[C]) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region NoneSelector

// NoneSelector is a Selector that will accept nothing. Its string representation is "{}".
type NoneSelector[C comparable] struct{}

var _ ByIndexSelector[string] = &NoneSelector[string]{}
var _ serialiser = &NoneSelector[string]{}

func (s *NoneSelector[C]) SelectFrom(cs []C) []C {
	if cs == nil {
		return nil
	}
	return []C{}
}

func (s *NoneSelector[C]) SelectMapFrom(ids []C) map[int]C {
	if ids == nil {
		return nil
	}
	return map[int]C{}
}

func (s *NoneSelector[C]) String() string {
	return "{}"
}

func (s *NoneSelector[C]) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.String() + `"`), nil
}

func (s *NoneSelector[C]) UnmarshalJSON(bts []byte) error {
	if string(bts) != `"{}"` {
		return fmt.Errorf(`only "{}" can be unmarshalled to NoneSelector, not %s`, string(bts))
	}
	return nil
}

func (s *NoneSelector[C]) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *NoneSelector[C]) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region OneSelector

// OneSelector is a Selector that will accept exactly one input: its Value. Its string representation is "<Value>"
// where <Value> is the string representation of the value.
type OneSelector[C comparable] struct {
	Value C
}

var _ Selector[string] = &OneSelector[string]{}
var _ serialiser = &OneSelector[string]{}

func (s *OneSelector[C]) SelectFrom(cs []C) []C {
	return FilterSlice(cs, func(_ int, item C) bool {
		return item == s.Value
	})
}

func (s *OneSelector[C]) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf(`%v`, s.Value)
}

func (s *OneSelector[C]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot marshal nil *OneSelector")
	}
	return []byte(`"` + s.String() + `"`), nil
}

func (s *OneSelector[C]) UnmarshalJSON(bts []byte) error {
	str, err := unmarshalString(bts)
	if err != nil {
		return err
	}
	return parseValue(str, &s.Value)
}

func (s *OneSelector[C]) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *OneSelector[C]) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region SetSelector

// SetSelector is a Selector that will accept only values of its Set. Its string representation is "{<item1>, <item2>, ..., <itemN>}"
// where <itemI> is the string representation of the I'th item in the Set.
type SetSelector[C comparable] struct {
	Set []C
}

var _ Selector[string] = &SetSelector[string]{}
var _ serialiser = &SetSelector[string]{}

func (s *SetSelector[C]) SelectFrom(cs []C) []C {
	return SetIntersection(cs, s.Set)
}

func (s *SetSelector[C]) String() string {
	if s == nil {
		return "<nil>"
	}
	strs := TransformSlice(s.Set, func(_ int, c C) string {
		return fmt.Sprintf("%v", c)
	})
	return "{" + strings.Join(strs, ", ") + "}"
}

func (s *SetSelector[C]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot marshal nil *SetSelector")
	}
	return []byte(`"` + s.String() + `"`), nil
}

func (s *SetSelector[C]) UnmarshalJSON(bts []byte) error {
	str, err := unmarshalString(bts)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(str, `{`) {
		return fmt.Errorf(`cannot unmarshal as SetSelector string not starting with '{'`)
	}
	if !strings.HasSuffix(str, `}`) {
		return fmt.Errorf(`cannot unmarshal as SetSelector string not ending with '}'`)
	}
	str = str[1 : len(str)-1]
	parts := strings.Split(str, ",")
	s.Set = make([]C, len(parts))
	for i := range parts {
		err := parseValue(strings.TrimSpace(parts[i]), &s.Set[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SetSelector[C]) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *SetSelector[C]) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region NotSelector

// NotSelector is a Selector that will accept only values not on its BlackList. Its string representation is "!{<item1>, <item2>, ..., <itemN>}"
// where <itemI> is the string representation of the I'th item on the BlackList.
type NotSelector[C comparable] struct {
	BlackList []C
}

var _ Selector[string] = &NotSelector[string]{}
var _ serialiser = &NotSelector[string]{}

func (s *NotSelector[C]) SelectFrom(cs []C) []C {
	return SetDiff(cs, s.BlackList)
}

func (s *NotSelector[C]) String() string {
	if s == nil {
		return "<nil>"
	}
	strs := TransformSlice(s.BlackList, func(_ int, c C) string {
		return fmt.Sprintf("%v", c)
	})
	return "!{" + strings.Join(strs, ", ") + "}"
}

func (s *NotSelector[C]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot marshal nil *NotSelector")
	}
	return []byte(`"` + s.String() + `"`), nil
}

func (s *NotSelector[C]) UnmarshalJSON(bts []byte) error {
	str, err := unmarshalString(bts)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(str, `!{`) {
		return fmt.Errorf(`cannot unmarshal as SetSelector string not starting with '!{'`)
	}
	if !strings.HasSuffix(str, `}`) {
		return fmt.Errorf(`cannot unmarshal as SetSelector string not ending with '}'`)
	}
	str = str[2 : len(str)-1]
	parts := strings.Split(str, ",")
	s.BlackList = make([]C, len(parts))
	for i := range parts {
		err := parseValue(strings.TrimSpace(parts[i]), &s.BlackList[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *NotSelector[C]) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *NotSelector[C]) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region IntervalSelector

// IntervalSelector is a Selector that will accept only values between its boundaries. Its string representation is
// "[<from>, <to>]" where <from> and <to> are the string representations of the respective boundaries. If the bracket
// of a boundary is flipped (i.e. ']' instead of '[' or vice versa) the value of the boundary will not be included in
// the selection. If a boundary value is empty (e.g. "[a,]" or "[,z]") the boundary is considered open (no matter the
// type of bracket).
type IntervalSelector[C comparable] struct {
	From, To Boundary[C]
	// Compare is used to compare values to the boundary values. If nil CompareKnownTypes is used,
	// meaning the used generic type A must be supported by that function.
	Compare CompareFunc[C]
}

var _ Selector[string] = &IntervalSelector[string]{}
var _ serialiser = &IntervalSelector[string]{}

func (s *IntervalSelector[A]) SelectFrom(cs []A) []A {
	i := Interval[A]{Lower: s.From, Upper: s.To}
	return FilterSlice(cs, func(_ int, item A) bool {
		return i.Contains(item, s.Compare)
	})
}

func (s *IntervalSelector[C]) String() string {
	if s == nil {
		return "<nil>"
	}
	return intervalSelectorToString(s.From, s.To, func(c C) string {
		return fmt.Sprintf("%v", c)
	})
}

func (s *IntervalSelector[C]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot marshal nil *IntervalSelector")
	}
	return []byte(`"` + s.String() + `"`), nil
}

func (s *IntervalSelector[C]) UnmarshalJSON(bts []byte) (err error) {
	s.From, s.To, err = unmarshalIntervalSelectorJSON(bts, parseValue[C])
	return
}

func (s *IntervalSelector[C]) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *IntervalSelector[C]) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region RegexSelector

// RegexSelector is a Selector for strings that only selects strings matching its Regex. Its string representation
// is "<regex>" where <regex> is a regex starting with '^' and ending with '$'.
type RegexSelector struct {
	Regex *regexp.Regexp
}

var _ Selector[string] = &RegexSelector{}
var _ serialiser = &RegexSelector{}

func (s *RegexSelector) SelectFrom(cs []string) []string {
	return FilterSlice(cs, func(_ int, item string) bool {
		return s.Regex.MatchString(item)
	})
}

func (s *RegexSelector) String() string {
	if s == nil {
		return "<nil>"
	}
	str := s.Regex.String()
	if !strings.HasPrefix(str, "^") {
		str = "^" + str
	}
	if !strings.HasSuffix(str, "$") {
		str = str + "$"
	}
	return str
}

func (s *RegexSelector) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot marshal nil *RegexSelector")
	}
	if s.Regex == nil {
		return nil, fmt.Errorf("cannot marshal RegexSelector with no regex")
	}
	return []byte(`"` + s.String() + `"`), nil
}

func (s *RegexSelector) UnmarshalJSON(bts []byte) error {
	str, err := unmarshalString(bts)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(str, "^") {
		return fmt.Errorf(`cannot unmarshal as RegexSelector string not starting with '^'`)
	}
	if !strings.HasSuffix(str, "$") {
		return fmt.Errorf(`cannot unmarshal as RegexSelector string not ending with '$'`)
	}
	s.Regex, err = regexp.Compile(str)
	return err
}

func (s *RegexSelector) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *RegexSelector) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region UnionSelector

type UnionSelector[C comparable] struct {
	SubSelectors []Selector[C]
}

var _ Selector[string] = UnionSelector[string]{}

func (s UnionSelector[C]) SelectFrom(cs []C) []C {
	if cs == nil {
		return nil
	}
	unselected := make([]C, len(cs))
	copy(unselected, cs)
	out := make([]C, 0, len(cs))
	for i := range s.SubSelectors {
		selected := s.SubSelectors[i].SelectFrom(unselected)
		unselected = SetDiff(unselected, selected)
		out = append(out, selected...)
	}
	slices.SortFunc(out, func(left, right C) bool {
		leftIdx := slices.Index(cs, left)
		rightIdx := slices.Index(cs, right)
		return leftIdx < rightIdx
	})
	return out
}

//#endregion

//#region IntersectionSelector

type IntersectionSelector[C comparable] struct {
	SubSelectors []Selector[C]
}

var _ Selector[string] = IntersectionSelector[string]{}

func (s IntersectionSelector[C]) SelectFrom(cs []C) []C {
	if cs == nil {
		return nil
	}
	out := make([]C, len(cs))
	copy(out, cs)
	for i := range s.SubSelectors {
		if len(out) == 0 {
			break
		}
		out = s.SubSelectors[i].SelectFrom(out)
	}
	return out
}

//#endregion

//#region OneIndexSelector

type OneIndexSelector[C comparable] struct {
	Index int
}

var _ ByIndexSelector[string] = &OneIndexSelector[string]{}
var _ serialiser = &OneIndexSelector[string]{}

func (s *OneIndexSelector[C]) SelectFrom(cs []C) []C {
	return GetMapValuesSortedByKey(s.SelectMapFrom(cs), nil)
}

func (s *OneIndexSelector[C]) SelectMapFrom(cs []C) map[int]C {
	if cs == nil {
		return nil
	}
	if len(cs) == 0 {
		return map[int]C{}
	}
	actI, outOfBounds := toActualIndex(s.Index, len(cs))
	if outOfBounds {
		return map[int]C{}
	}
	return map[int]C{s.Index: cs[actI]}
}

func (s *OneIndexSelector[C]) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf(`#%d`, s.Index)
}

func (s *OneIndexSelector[C]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot marshal nil *OneIndexSelector")
	}
	return []byte(`"` + s.String() + `"`), nil
}

func (s *OneIndexSelector[C]) UnmarshalJSON(bts []byte) error {
	str, err := unmarshalString(bts)
	if err != nil {
		return err
	}
	err = parseIndex(str, &s.Index)
	return err
}

func (s *OneIndexSelector[C]) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *OneIndexSelector[C]) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region IndexSetSelector

type IndexSetSelector[C comparable] struct {
	IndexSet []int
}

var _ ByIndexSelector[string] = &IndexSetSelector[string]{}
var _ serialiser = &IndexSetSelector[string]{}

func (s *IndexSetSelector[C]) SelectFrom(cs []C) []C {
	return GetMapValuesSortedByKey(s.SelectMapFrom(cs), nil)
}

func (s *IndexSetSelector[C]) SelectMapFrom(cs []C) map[int]C {
	if cs == nil {
		return nil
	}
	if len(cs) == 0 {
		return map[int]C{}
	}
	out := map[int]C{}
	for _, extI := range s.IndexSet {
		actI, outOfBounds := toActualIndex(extI, len(cs))
		if outOfBounds {
			continue
		}
		out[extI] = cs[actI]
	}
	return out
}

func (s *IndexSetSelector[C]) String() string {
	if s == nil {
		return "<nil>"
	}
	strs := TransformSlice(s.IndexSet, func(_ int, setIdx int) string {
		return "#" + strconv.Itoa(setIdx)
	})
	return "{" + strings.Join(strs, ", ") + "}"
}

func (s *IndexSetSelector[C]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot marshal nil *IndexSetSelector")
	}
	return []byte(`"` + s.String() + `"`), nil
}

func (s *IndexSetSelector[C]) UnmarshalJSON(bts []byte) error {
	str, err := unmarshalString(bts)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(str, `{`) {
		return fmt.Errorf(`cannot unmarshal as IndexSetSelector string not starting with '{'`)
	}
	if !strings.HasSuffix(str, `}`) {
		return fmt.Errorf(`cannot unmarshal as IndexSetSelector string not ending with '}'`)
	}
	str = str[1 : len(str)-1]
	parts := strings.Split(str, ",")
	s.IndexSet = make([]int, len(parts))
	for i := range parts {
		err = parseIndex(strings.TrimSpace(parts[i]), &s.IndexSet[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *IndexSetSelector[C]) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *IndexSetSelector[C]) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region IndexIntervalSelector

type IndexIntervalSelector[C comparable] struct {
	From, To Boundary[int]
}

var _ ByIndexSelector[string] = &IndexIntervalSelector[string]{}
var _ serialiser = &IndexIntervalSelector[string]{}

func (s *IndexIntervalSelector[C]) SelectFrom(cs []C) []C {
	return GetMapValuesSortedByKey(s.SelectMapFrom(cs), nil)
}

func (s *IndexIntervalSelector[C]) SelectMapFrom(cs []C) map[int]C {
	if cs == nil {
		return nil
	}
	if len(cs) == 0 {
		return map[int]C{}
	}
	from, to, emptySelection := correctBoundariesForNbItems(s.From, s.To, len(cs))
	if emptySelection {
		return map[int]C{}
	}

	selectedActualToExtendedIndexes := map[int]int{}
	itv := Interval[int]{Lower: from, Upper: to}
	for extI := -len(cs); extI < len(cs); extI++ {
		if itv.Contains(extI, nil) {
			actI, _ := toActualIndex(extI, len(cs))
			// may overwrite smaller extI, but this is okay
			selectedActualToExtendedIndexes[actI] = extI
		}
	}

	out := map[int]C{}
	for actI := range selectedActualToExtendedIndexes {
		extI := selectedActualToExtendedIndexes[actI]
		out[extI] = cs[actI]
	}
	return out
}

// correctBoundariesForNbItems transforms the given boundaries to new boundaries that are enabled and within the extended index range.
// It also checks if the selection will be empty with the given boundaries, e.g. when selecting the second item
// when he input only contains one.
func correctBoundariesForNbItems(from, to Boundary[int], nbItems int) (corrFrom Boundary[int], corrTo Boundary[int], emptySelection bool) {
	if from.closed {
		if from.value < -nbItems {
			// from is set below the extended set of ids, so it doesn't constrain the selection
			from.closed = false
		}
		if from.value >= nbItems {
			// from is set above the range of the ids, so no id can be selected
			// e.g. from set to 3 when there are only two ids given
			return from, to, true
		}
	}
	if to.closed {
		if to.value < -nbItems {
			// to is set below the range of the ids, so no id can be selected
			// e.g. to set to -2 (= second to last) when there is only one id given
			return from, to, true
		}
		if to.value >= nbItems {
			// to is set above the range of ids, so it doesn't constrain the selection
			to.closed = false
		}
	}

	if !from.closed && !to.closed {
		from = Incl(0)
		to = Incl(nbItems - 1)
	}
	if from.closed && !to.closed { // only From applies
		if from.value < 0 {
			to = Incl(0)
		} else {
			to = Incl(nbItems - 1)
		}
	}
	if !from.closed && to.closed { // only To applies
		if to.value < 0 {
			from = Incl(-nbItems)
		} else {
			from = Incl(0)
		}
	}
	return from, to, false
}

func toActualIndex(extI int, nbItems int) (actI int, outOfBounds bool) {
	if extI < -nbItems || extI >= nbItems {
		return 0, true
	}
	actI = extI
	if extI < 0 {
		actI = extI + nbItems
	}
	return
}

func (s *IndexIntervalSelector[C]) String() string {
	if s == nil {
		return "<nil>"
	}
	return intervalSelectorToString(s.From, s.To, func(i int) string {
		return "#" + strconv.Itoa(i)
	})
}

func (s *IndexIntervalSelector[C]) MarshalJSON() ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("cannot marshal nil *IntervalSelector")
	}
	return []byte(`"` + s.String() + `"`), nil
}

func (s *IndexIntervalSelector[C]) UnmarshalJSON(bts []byte) (err error) {
	s.From, s.To, err = unmarshalIntervalSelectorJSON(bts, parseIndex)
	return
}

func (s *IndexIntervalSelector[C]) MarshalYAML() (any, error) {
	return s.String(), nil
}

func (s *IndexIntervalSelector[C]) UnmarshalYAML(n *yaml.Node) error {
	return s.UnmarshalJSON(yamlNodeToJson(n))
}

//#endregion

//#region UnmarshallableStringSelector

type unmarshallableSelector[C comparable] struct {
	MultiTypeUnmarshaler[Selector[C]]
}

func (u *unmarshallableSelector[C]) SelectFrom(cs []C) []C {
	panic("calling SelectFrom on not-unmarshalled unmarshallableSelector")
}

func UnmarshallableStringSelector(setResult func(s Selector[string])) Selector[string] {
	return &unmarshallableSelector[string]{
		MultiTypeUnmarshaler: MultiTypeUnmarshaler[Selector[string]]{
			GetCandidates: func() []Selector[string] {
				return []Selector[string]{
					&AllSelector[string]{},
					&NoneSelector[string]{},
					&IndexSetSelector[string]{},
					&SetSelector[string]{},
					&NotSelector[string]{},
					&IndexIntervalSelector[string]{},
					&IntervalSelector[string]{},
					&RegexSelector{},
					&OneIndexSelector[string]{},
					&OneSelector[string]{},
					//TODO UnionSelector, IntersectionSelector
				}
			},
			SetResult: setResult,
		},
	}
}

func unmarshalString(bts []byte) (string, error) {
	str := string(bts)
	if !strings.HasPrefix(str, `"`) {
		return "", fmt.Errorf(`cannot unmarshal bytes not starting with '"' to string`)
	}
	if !strings.HasSuffix(str, `"`) {
		return "", fmt.Errorf(`cannot unmarshal bytes not ending with '"' to string`)
	}
	return str[1 : len(str)-1], nil
}

type fromStringer interface {
	FromString(s string) error
}

func parseValue[C comparable](s string, c *C) (err error) {
	v := reflect.ValueOf(c)
	switch any(*c).(type) {
	case string:
		v.Elem().Set(reflect.ValueOf(s))
		return
	}
	if fs, castOk := any(c).(fromStringer); castOk {
		return fs.FromString(s)
	}
	return fmt.Errorf("cannot parse value of type %T", c)
}

func parseIndex(str string, i *int) (err error) {
	if !strings.HasPrefix(str, "#") {
		return fmt.Errorf("cannot parse as index string not starting with '#'")
	}
	str = str[1:]
	*i, err = strconv.Atoi(str)
	return
}

func intervalSelectorToString[C comparable](from, to Boundary[C], toString func(c C) string) string {
	fromStr := ""
	fromVal, fromIncl, fromClosed := from.Value()
	if fromClosed {
		if fromIncl {
			fromStr = "["
		} else {
			fromStr = "]"
		}
		fromStr = fromStr + toString(fromVal)
	} else {
		fromStr = "]"
	}
	toStr := ""
	toVal, toIncl, toClosed := to.Value()
	if toClosed {
		if toIncl {
			toStr = "]"
		} else {
			toStr = "["
		}
		toStr = toString(toVal) + toStr
	} else {
		toStr = "["
	}
	return fromStr + `,` + toStr
}

func unmarshalIntervalSelectorJSON[C comparable](bts []byte, parse func(s string, c *C) error) (from, to Boundary[C], err error) {
	str, err := unmarshalString(bts)
	if err != nil {
		return Boundary[C]{}, Boundary[C]{}, err
	}
	fromIncluded := false
	switch str[0] {
	case '[':
		fromIncluded = true
	case ']':
		fromIncluded = false
	default:
		return Boundary[C]{}, Boundary[C]{}, fmt.Errorf(`cannot unmarshal as interval selector string not starting with '[' or ']'`)
	}
	toIncluded := false
	switch str[len(str)-1] {
	case ']':
		toIncluded = true
	case '[':
		toIncluded = false
	default:
		return Boundary[C]{}, Boundary[C]{}, fmt.Errorf(`cannot unmarshal as interval selector string not ending with ']' or '['`)
	}
	str = str[1 : len(str)-1]
	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		return Boundary[C]{}, Boundary[C]{}, fmt.Errorf("cannot unmarshal as interval selector string not containing exactly one ','")
	}
	from, err = parseBoundary(parts[0], fromIncluded, parse)
	if err != nil {
		return Boundary[C]{}, Boundary[C]{}, err
	}
	to, err = parseBoundary(parts[1], toIncluded, parse)
	if err != nil {
		return Boundary[C]{}, Boundary[C]{}, err
	}
	return
}

func parseBoundary[C comparable](s string, valueIncluded bool, parse func(s string, c *C) error) (Boundary[C], error) {
	if len(s) == 0 {
		return Open[C](), nil
	}
	val := new(C)
	err := parse(strings.TrimSpace(s), val)
	if err != nil {
		return Boundary[C]{}, err
	}
	if valueIncluded {
		return Incl(*val), nil
	} else {
		return Excl(*val), nil
	}
}

func yamlNodeToJson(n *yaml.Node) []byte {
	//TODO support other types than string
	s := n.Value
	if len(s) == 0 {
		return []byte(`""`)
	}
	switch s[0] {
	case '"':
		return []byte(s)
	case '\'':
		return []byte("'" + s[1:len(s)-1] + "'")
	default:
		return []byte(`"` + s + `"`)
	}
}

//#endregion

//#region IntervalActivityRefSelector

func IntervalActivityRefSelector(from, to Boundary[ActivityRef]) IntervalSelector[ActivityRef] {
	return IntervalSelector[ActivityRef]{
		From: from,
		To:   to,
		Compare: func(left, right ActivityRef) (leftIsSmaller bool, equal bool) {
			l := left.ActivityId
			r := right.ActivityId
			return l < r, l == r
		},
	}
}

//#endregion

//#region SeparatedVersionRefSelector

type SeparatedVersionRefSelector struct {
	ActivityRefSelector Selector[ActivityRef]
	VersionIdSelector   Selector[string]
}

var _ Selector[VersionRef] = SeparatedVersionRefSelector{}

func (s SeparatedVersionRefSelector) SelectFrom(versionRefs []VersionRef) []VersionRef {
	activityRefToVersionRefs := Group(versionRefs, func(vr VersionRef) ActivityRef {
		return vr.ActivityRef
	})
	passedActivityRefs := s.ActivityRefSelector.SelectFrom(GetMapKeys(activityRefToVersionRefs))
	out := make([]VersionRef, 0, len(versionRefs))
	for i := range passedActivityRefs {
		activityRef := passedActivityRefs[i]
		versionRefs := activityRefToVersionRefs[activityRef]
		versionIds := TransformSlice(versionRefs, func(_ int, ref VersionRef) string { return ref.VersionId })
		slices.Sort(versionIds)
		passedVersionIds := s.VersionIdSelector.SelectFrom(versionIds)
		passedVersionRefs := TransformSlice(passedVersionIds, func(_ int, vid string) VersionRef {
			return VersionRef{
				ActivityRef: activityRef,
				VersionId:   vid,
			}
		})
		out = append(out, passedVersionRefs...)
	}
	return out
}

func (s *SeparatedVersionRefSelector) MarshalJSON() ([]byte, error) {
	panic("IndexIntervalSelector#MarshalJSON not implemented yet")
}

func (s *SeparatedVersionRefSelector) UnmarshalJSON(bts []byte) error {
	panic("IndexIntervalSelector#UnmarshalJSON not implemented yet")
}

//#endregion
