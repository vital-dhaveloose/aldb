package aldb

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestOneSelector(t *testing.T) {
	s := OneSelector[string]{Value: "foo-bar"}

	assert.Contains(t, s.SelectFrom([]string{
		"foo-bar", "bar-foo",
	}), "foo-bar")
}

func TestOneByIndexSelector(t *testing.T) {
	s := OneIndexSelector[string]{Index: -1}

	selection := s.SelectMapFrom([]string{"foo", "bar", "zol"})

	assert.Len(t, selection, 1)
	assert.Equal(t, selection[-1], "zol")

	//TODO test cases with selected index from < -len(ids) to > len(ids)-1
}

var months = []string{"jan", "feb", "mar", "may", "apr", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}

func TestSetByIndexSelector(t *testing.T) {
	s := IndexSetSelector[string]{IndexSet: []int{-3, -2, -1}}

	selection := s.SelectMapFrom(months)

	assert.Len(t, selection, 3)
	assert.Equal(t, selection[-3], "oct")
	assert.Equal(t, selection[-2], "nov")
	assert.Equal(t, selection[-1], "dec")

	//TODO test cases with selected index from < -len(ids) to > len(ids)-1
}

func TestUnionSelector(t *testing.T) {
	s := UnionSelector[string]{
		SubSelectors: []Selector[string]{
			&NoneSelector[string]{},
			&OneSelector[string]{Value: "foo"},
			&SetSelector[string]{Set: []string{"bar", "zol"}},
			&IntervalSelector[string]{From: Incl("x")},
		},
	}

	selected := s.SelectFrom([]string{
		"foo", "bar", "zol",
		"x", "xylophone", "y", "yoga", "zebra",
		"car", "bird", "word", "", "fish",
	})

	assert.ElementsMatch(t, selected, []string{
		"foo", "bar", "zol",
		"x", "xylophone", "y", "yoga", "zebra",
	})
}

func TestFromToIndexSelector(t *testing.T) {
	s := IndexIntervalSelector[string]{
		From: Incl(-4),
		To:   Excl(-2),
	}

	selection := s.SelectMapFrom([]string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"})

	assert.Len(t, selection, 2)
	assert.Equal(t, selection[-4], "thu")
	assert.Equal(t, selection[-3], "fri")
}

func TestSeparatedSelector(t *testing.T) {
	s := SeparatedVersionRefSelector{
		ActivityRefSelector: &SetSelector[ActivityRef]{Set: []ActivityRef{
			{ActivityId: "foo"},
			{ActivityId: "bar"},
		}},
		VersionIdSelector: &IndexIntervalSelector[string]{
			From: Incl(-2), // previous
			To:   Incl(-1), // latest
		},
	}

	selection := s.SelectFrom([]VersionRef{
		Vref("foo", "12:00"),
		Vref("foo", "13:00"),
		Vref("foo", "14:00"),
		Vref("bar", "14:00"),
		Vref("zol", "13:00"),
		Vref("zol", "14:00"),
	})

	assert.Len(t, selection, 3)
	assert.ElementsMatch(t, selection, []VersionRef{
		Vref("foo", "13:00"),
		Vref("foo", "14:00"),
		Vref("bar", "14:00"),
	})

}

func Vref(activityId string, versionId string) VersionRef {
	return VersionRef{
		ActivityRef: ActivityRef{
			ActivityId: activityId,
		},
		VersionId: versionId,
	}
}

func TestParseValue(t *testing.T) {
	v := &VersionRef{}

	assert.Equal(t, "|", fmt.Sprintf("%v", v))

	err := parseValue("one|two", v)
	assert.NoError(t, err)
	assert.Equal(t, v.ActivityRef.ActivityId, "foo")
	assert.Equal(t, v.VersionId, "bar")
}

type Request struct {
	Selector Selector[string]
}

func TestUnmarshalJSON(t *testing.T) {
	type testCase struct {
		name           string
		stringRep      string
		assertSelector func(t *testing.T, s Selector[string])
	}

	cases := []testCase{
		{
			name:      "all",
			stringRep: `*`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				_, castOk := s.(*AllSelector[string])
				assert.True(t, castOk)
			},
		},
		{
			name:      "none",
			stringRep: `{}`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				_, castOk := s.(*NoneSelector[string])
				assert.True(t, castOk)
			},
		},
		{
			name:      "one",
			stringRep: `foobar`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				os, castOk := s.(*OneSelector[string])
				assert.True(t, castOk)
				assert.Equal(t, "foobar", os.Value)
			},
		},
		{
			name:      "set",
			stringRep: `{foo, bar}`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				ss, castOk := s.(*SetSelector[string])
				assert.True(t, castOk)
				if castOk {
					assert.ElementsMatch(t, []string{"foo", "bar"}, ss.Set)
				}
			},
		},
		{
			name:      "not",
			stringRep: `!{foo, bar}`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				ns, castOk := s.(*NotSelector[string])
				assert.True(t, castOk)
				if castOk {
					assert.ElementsMatch(t, []string{"foo", "bar"}, ns.BlackList)
				}
			},
		},
		{
			name:      "closed-interval",
			stringRep: `[foo,bar[`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				is, castOk := s.(*IntervalSelector[string])
				assert.True(t, castOk)
				if castOk {
					fromVal, fromInclusive, fromClosed := is.From.Value()
					assert.True(t, fromClosed)
					assert.True(t, fromInclusive)
					assert.Equal(t, fromVal, "foo")

					toVal, toInclusive, toClosed := is.To.Value()
					assert.True(t, toClosed)
					assert.False(t, toInclusive)
					assert.Equal(t, toVal, "bar")
				}
			},
		},
		{
			name:      "half-open-interval",
			stringRep: `[foo,[`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				is, castOk := s.(*IntervalSelector[string])
				assert.True(t, castOk)
				if castOk {
					fromVal, fromInclusive, fromClosed := is.From.Value()
					assert.True(t, fromClosed)
					assert.True(t, fromInclusive)
					assert.Equal(t, fromVal, "foo")

					_, _, toClosed := is.To.Value()
					assert.False(t, toClosed)
				}
			},
		},
		{
			name:      "regex",
			stringRep: `^foobar$`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				rs, castOk := s.(*RegexSelector)
				assert.True(t, castOk)
				if castOk {
					assert.Equal(t, "^foobar$", rs.Regex.String())
				}
			},
		},
		{
			name:      "one-index",
			stringRep: `#-1`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				os, castOk := s.(*OneIndexSelector[string])
				assert.True(t, castOk)
				assert.Equal(t, -1, os.Index)
			},
		},
		{
			name:      "set-index",
			stringRep: `{#-1, #0}`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				ss, castOk := s.(*IndexSetSelector[string])
				assert.True(t, castOk)
				if castOk {
					assert.ElementsMatch(t, []int{-1, 0}, ss.IndexSet)
				}
			},
		},
		{
			name:      "closed-interval-index",
			stringRep: `[#-2,#1[`,
			assertSelector: func(t *testing.T, s Selector[string]) {
				is, castOk := s.(*IndexIntervalSelector[string])
				assert.True(t, castOk)
				if castOk {
					fromVal, fromInclusive, fromClosed := is.From.Value()
					assert.True(t, fromClosed)
					assert.True(t, fromInclusive)
					assert.Equal(t, fromVal, -2)

					toVal, toInclusive, toClosed := is.To.Value()
					assert.True(t, toClosed)
					assert.False(t, toInclusive)
					assert.Equal(t, toVal, 1)
				}
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := &Request{}
			req.Selector = UnmarshallableStringSelector(func(s Selector[string]) {
				req.Selector = s
			})

			err := json.Unmarshal([]byte(
				fmt.Sprintf(`{"selector": "%s"}`, tc.stringRep),
			), req)

			assert.NoError(t, err)
			tc.assertSelector(t, req.Selector)

			jso, err := json.Marshal(req.Selector)
			assert.NoError(t, err)
			assert.Equal(t, fmt.Sprintf(`"%s"`, tc.stringRep), string(jso))
		})
	}
}

func TestUnmarshalYAML(t *testing.T) {
	type testCase struct {
		name           string
		stringRep      string
		assertSelector func(t *testing.T, s Selector[string]) bool
	}

	cases := []testCase{
		{
			name:      "all",
			stringRep: `'*'`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				_, castOk := s.(*AllSelector[string])
				return assert.True(t, castOk)
			},
		},
		{
			name:      "none",
			stringRep: `'{}'`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				_, castOk := s.(*NoneSelector[string])
				return assert.True(t, castOk)
			},
		},
		{
			name:      "one",
			stringRep: `foobar`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				os, castOk := s.(*OneSelector[string])
				return assert.True(t, castOk) &&
					assert.Equal(t, "foobar", os.Value)
			},
		},
		{
			name:      "set",
			stringRep: `'{foo, bar}'`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				ss, castOk := s.(*SetSelector[string])
				return assert.True(t, castOk) &&
					assert.ElementsMatch(t, []string{"foo", "bar"}, ss.Set)
			},
		},
		{
			name:      "not",
			stringRep: `'!{foo, bar}'`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				ns, castOk := s.(*NotSelector[string])
				return assert.True(t, castOk) &&
					assert.ElementsMatch(t, []string{"foo", "bar"}, ns.BlackList)
			},
		},
		{
			name:      "closed-interval",
			stringRep: `'[foo,bar['`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				is, castOk := s.(*IntervalSelector[string])
				if assert.True(t, castOk) {
					fromVal, fromInclusive, fromClosed := is.From.Value()
					toVal, toInclusive, toClosed := is.To.Value()

					return assert.True(t, fromClosed) &&
						assert.True(t, fromInclusive) &&
						assert.Equal(t, fromVal, "foo") &&
						assert.True(t, toClosed) &&
						assert.False(t, toInclusive) &&
						assert.Equal(t, toVal, "bar")
				}
				return false
			},
		},
		{
			name:      "half-open-interval",
			stringRep: `'[foo,['`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				is, castOk := s.(*IntervalSelector[string])
				if assert.True(t, castOk) {
					fromVal, fromInclusive, fromClosed := is.From.Value()
					_, _, toClosed := is.To.Value()

					return assert.True(t, fromClosed) &&
						assert.True(t, fromInclusive) &&
						assert.Equal(t, fromVal, "foo") &&
						assert.False(t, toClosed)
				}
				return false
			},
		},
		{
			name:      "regex",
			stringRep: `^foobar$`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				rs, castOk := s.(*RegexSelector)
				return assert.True(t, castOk) &&
					assert.Equal(t, "^foobar$", rs.Regex.String())
			},
		},
		{
			name:      "one-index",
			stringRep: `'#-1'`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				os, castOk := s.(*OneIndexSelector[string])
				return assert.True(t, castOk) &&
					assert.Equal(t, -1, os.Index)
			},
		},
		{
			name:      "set-index",
			stringRep: `'{#-1, #0}'`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				ss, castOk := s.(*IndexSetSelector[string])
				return assert.True(t, castOk) &&
					assert.ElementsMatch(t, []int{-1, 0}, ss.IndexSet)
			},
		},
		{
			name:      "closed-interval-index",
			stringRep: `'[#-2,#1['`,
			assertSelector: func(t *testing.T, s Selector[string]) bool {
				is, castOk := s.(*IndexIntervalSelector[string])
				if assert.True(t, castOk) {
					fromVal, fromInclusive, fromClosed := is.From.Value()
					toVal, toInclusive, toClosed := is.To.Value()

					return assert.True(t, fromClosed) &&
						assert.True(t, fromInclusive) &&
						assert.Equal(t, fromVal, -2) &&
						assert.True(t, toClosed) &&
						assert.False(t, toInclusive) &&
						assert.Equal(t, toVal, 1)
				}
				return false
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var slr Selector[string]
			slr = UnmarshallableStringSelector(func(s Selector[string]) {
				slr = s
			})

			err := yaml.Unmarshal([]byte(tc.stringRep), slr)

			if assert.NoError(t, err) {
				tc.assertSelector(t, slr)

				yml, err := yaml.Marshal(slr)
				if assert.NoError(t, err) {
					assert.Equal(t, tc.stringRep, strings.TrimSpace(string(yml)))
				}
			}
		})
	}
}

func TestUnmarshalJSONIntervalActivityRefSelector(t *testing.T) {
	s := IntervalActivityRefSelector(Open[ActivityRef](), Open[ActivityRef]())

	err := json.Unmarshal([]byte(`"[foo,bar["`), &s)

	assert.NoError(t, err)

	fromVal, fromInclusive, fromClosed := s.From.Value()
	assert.True(t, fromClosed)
	assert.True(t, fromInclusive)
	assert.Equal(t, fromVal.ActivityId, "foo")

	toVal, toInclusive, toClosed := s.To.Value()
	assert.True(t, toClosed)
	assert.False(t, toInclusive)
	assert.Equal(t, toVal.ActivityId, "bar")
}
