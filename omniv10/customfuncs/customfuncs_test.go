package customfuncs

import (
	"sort"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/jf-tech/go-corelib/jsons"
	"github.com/jf-tech/go-corelib/strs"
	"github.com/jf-tech/omniparser/transformctx"
	"github.com/stretchr/testify/assert"

	"github.com/jf-tech/omniparserlegacy/omniv10/errs"
)

func TestDumpOmniV10OnlyCustomFuncNames(t *testing.T) {
	var names []string
	for name := range OmniV10OnlyCustomFuncs {
		names = append(names, name)
	}
	sort.Strings(names)
	cupaloy.SnapshotT(t, jsons.BPM(names))
}

func TestContainsPattern(t *testing.T) {
	for _, test := range []struct {
		name     string
		pattern  string
		strs     []string
		err      string
		expected string
	}{
		{
			name:    "invalid pattern",
			pattern: "[",
			err:     "error parsing regexp: missing closing ]: `[`",
		},
		{
			name:     "empty strs",
			pattern:  ".*",
			strs:     nil,
			err:      "",
			expected: "false",
		},
		{
			name:     "contains",
			pattern:  "^[0-9]+$",
			strs:     []string{"abc", "efg", "123"},
			err:      "",
			expected: "true",
		},
		{
			name:     "not contains",
			pattern:  "^[0-9]+$",
			strs:     []string{"abc", "efg", "x123"},
			err:      "",
			expected: "false",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := ContainsPattern(nil, test.pattern, test.strs...)
			if test.err != "" {
				assert.Error(t, err)
				assert.Equal(t, test.err, err.Error())
				assert.Equal(t, "", result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestExternal(t *testing.T) {
	for _, test := range []struct {
		name      string
		externals map[string]string
		lookup    string
		err       string
		expected  string
	}{
		{
			name:      "externals nil",
			externals: nil,
			lookup:    "abc",
			err:       "cannot find external property 'abc'",
			expected:  "",
		},
		{
			name:      "externals empty",
			externals: map[string]string{},
			lookup:    "efg",
			err:       "cannot find external property 'efg'",
			expected:  "",
		},
		{
			name:      "not found",
			externals: map[string]string{"abc": "abc"},
			lookup:    "efg",
			err:       "cannot find external property 'efg'",
			expected:  "",
		},
		{
			name:      "found",
			externals: map[string]string{"abc": "123"},
			lookup:    "abc",
			err:       "",
			expected:  "123",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			v, err := External(
				&transformctx.Ctx{ExternalProperties: test.externals},
				test.lookup,
			)
			switch {
			case strs.IsStrNonBlank(test.err):
				assert.Error(t, err)
				assert.Equal(t, test.err, err.Error())
				assert.Equal(t, "", v)
			default:
				assert.NoError(t, err)
				assert.Equal(t, test.expected, v)
			}
		})
	}
}

func TestFloor(t *testing.T) {
	for _, test := range []struct {
		name          string
		value         string
		decimalPlaces string
		err           string
		expected      string
	}{
		{
			name:          "invalid value",
			value:         "??",
			decimalPlaces: "2",
			err:           `unable to parse value '??' to float64: strconv.ParseFloat: parsing "??": invalid syntax`,
			expected:      "",
		},
		{
			name:          "invalid decimal place value",
			value:         "3.1415926",
			decimalPlaces: "??",
			err:           `unable to parse decimal place value '??' to int: strconv.Atoi: parsing "??": invalid syntax`,
			expected:      "",
		},
		{
			name:          "decimal places less than 0",
			value:         "3.1415926",
			decimalPlaces: "-1",
			err:           `decimal place value must be an integer with range of [0,100], instead, got -1`,
			expected:      "",
		},
		{
			name:          "decimal places > 100",
			value:         "3.1415926",
			decimalPlaces: "101",
			err:           `decimal place value must be an integer with range of [0,100], instead, got 101`,
			expected:      "",
		},
		{
			name:          "decimal places less than available digits",
			value:         "3.1415926",
			decimalPlaces: "2",
			err:           "",
			expected:      "3.14",
		},
		{
			name:          "decimal places 0",
			value:         "3.1415926",
			decimalPlaces: "0",
			err:           "",
			expected:      "3",
		},
		{
			name:          "decimal places more than available digits",
			value:         "3.1415926",
			decimalPlaces: "20",
			err:           "",
			expected:      "3.1415926",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := Floor(nil, test.value, test.decimalPlaces)
			if test.err != "" {
				assert.Error(t, err)
				assert.Equal(t, test.err, err.Error())
				assert.Equal(t, "", result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestIfElse(t *testing.T) {
	for _, test := range []struct {
		name     string
		cv       []string
		err      string
		expected string
	}{
		{
			name: "wrong arg number - 0",
			cv:   nil,
			err:  "arg number must be odd, but got: 0",
		},
		{
			name: "wrong arg number - even",
			cv:   []string{"a", "b"},
			err:  "arg number must be odd, but got: 2",
		},
		{
			name: "condition not bool string",
			cv:   []string{"not bool string", "b", "c"},
			err:  "condition argument must be a boolean string, but got: not bool string",
		},
		{
			name:     "one of the conditions met",
			cv:       []string{"false", "abc", "true", "123", "rest"},
			err:      "",
			expected: "123",
		},
		{
			name:     "none of the conditions met",
			cv:       []string{"false", "abc", "false", "123", "rest"},
			err:      "",
			expected: "rest",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := IfElse(nil, test.cv...)
			if test.err != "" {
				assert.Error(t, err)
				assert.Equal(t, test.err, err.Error())
				assert.Equal(t, "", result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestIfEmpty(t *testing.T) {
	r, err := IsEmpty(nil, "")
	assert.NoError(t, err)
	assert.Equal(t, "true", r)
	r, err = IsEmpty(nil, "   ")
	assert.NoError(t, err)
	assert.Equal(t, "false", r)
	r, err = IsEmpty(nil, "abc")
	assert.NoError(t, err)
	assert.Equal(t, "false", r)
}

func TestReplace(t *testing.T) {
	regexStr := "[[:punct:]]|\\s+"
	replaceStr := ""
	// removed special char
	s, err := Replace(nil, "%%AE", regexStr, replaceStr)
	assert.NoError(t, err)
	assert.Equal(t, "AE", s)
	// removed whitespace
	s, err = Replace(nil, "%   %AE", regexStr, replaceStr)
	assert.NoError(t, err)
	assert.Equal(t, "AE", s)
	// nil/ or empty string
	s, err = Replace(nil, "  ", regexStr, replaceStr)
	assert.NoError(t, err)
	assert.Equal(t, "", s)
	// invalid regex
	s, err = Replace(nil, "abc", "[", replaceStr)
	assert.Error(t, err)
}

func TestRetrieveBySplit(t *testing.T) {
	for _, test := range []struct {
		name     string
		str      string
		sep      string
		index    string
		err      string
		expected string
	}{
		{
			name:     "invalid index",
			str:      "",
			sep:      "",
			index:    "invalid",
			err:      `strconv.Atoi: parsing "invalid": invalid syntax`,
			expected: "",
		},
		{
			name:     "index out of range: lower bound",
			str:      "123",
			sep:      "",
			index:    "-1",
			err:      `'index' is out of range. original string: 123, separator: , index: -1`,
			expected: "",
		},
		{
			name:     "index out of range: upper bound",
			str:      "123",
			sep:      "",
			index:    "12345",
			err:      `'index' is out of range. original string: 123, separator: , index: 12345`,
			expected: "",
		},
		{
			name:     "str non-empty, sep empty -> utf-8 explosion",
			str:      "123",
			sep:      "",
			index:    "1",
			err:      ``,
			expected: "2",
		},
		{
			name:     "str empty, sep empty",
			str:      "",
			sep:      "",
			index:    "0",
			err:      `'index' is out of range. original string: , separator: , index: 0`,
			expected: "",
		},
		{
			name:     "str non-empty, sep non-empty",
			str:      "1,2,3",
			sep:      ",",
			index:    "2",
			err:      ``,
			expected: "3",
		},
		{
			name:     "sep not found",
			str:      "123",
			sep:      ",",
			index:    "0",
			err:      ``,
			expected: "123",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r, err := RetrieveBySplit(nil, test.str, test.sep, test.index)
			if test.err == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, r)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.err, err.Error())
				assert.Equal(t, "", r)
			}
		})
	}
}

func TestRowSkip(t *testing.T) {
	// loose regex
	s, err := RowSkip(nil, "abcdefg", "abc")
	assert.True(t, errs.IsNonErrorRecordSkipped(err))
	assert.Equal(t, "", s)

	// strict regex - no match, i.e. no row skipping
	s, err = RowSkip(nil, "abcdefg", "^abc$")
	assert.NoError(t, err)
	assert.Equal(t, "", s)

	// strict regex - match, row skipping.
	s, err = RowSkip(nil, "abc", "^abc$")
	assert.True(t, errs.IsNonErrorRecordSkipped(err))
	assert.Equal(t, "", s)

	// invalid regex
	s, err = RowSkip(nil, "abc", "[")
	assert.Error(t, err)
	// but not the row skip err
	assert.False(t, errs.IsNonErrorRecordSkipped(err))
}

func TestSplitIntoJSONArray(t *testing.T) {
	for _, test := range []struct {
		name     string
		str      string
		sep      string
		trim     []string
		err      string
		expected string
	}{
		{
			name:     "len(trim) > 1",
			str:      "",
			sep:      "",
			trim:     []string{"true", "true"},
			err:      `cannot specify 'trim' argument more than once`,
			expected: "",
		},
		{
			name:     "invalid trim value",
			str:      "",
			sep:      "",
			trim:     []string{"invalid"},
			err:      `expect 'trim' to be of boolean string, but got: 'invalid'`,
			expected: "",
		},
		{
			name:     "sep empty -> utf-8 explosion",
			str:      "12   34",
			sep:      "",
			trim:     nil,
			err:      ``,
			expected: `["1","2"," "," "," ","3","4"]`,
		},
		{
			name:     "no sep found",
			str:      "  12   34  ",
			sep:      ",",
			trim:     []string{"true"},
			err:      ``,
			expected: `["12   34"]`,
		},
		{
			name:     "str empty",
			str:      "",
			sep:      ",",
			trim:     []string{"false"},
			err:      ``,
			expected: `[]`,
		},
		{
			name:     "sep found",
			str:      " 1 | 2 | 3 ",
			sep:      "|",
			trim:     []string{"false"},
			err:      ``,
			expected: `[" 1 "," 2 "," 3 "]`,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r, err := SplitIntoJSONArray(nil, test.str, test.sep, test.trim...)
			if test.err == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, r)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.err, err.Error())
				assert.Equal(t, "", r)
			}
		})
	}
}

func TestStrEqualAny(t *testing.T) {
	for _, test := range []struct {
		name     string
		s        string
		others   []string
		expected string
	}{
		{
			name:     "empty others",
			s:        "abc",
			others:   nil,
			expected: "false",
		},
		{
			name:     "no match in others",
			s:        "abc",
			others:   []string{"", "ab", "a b c"},
			expected: "false",
		},
		{
			name:     "found",
			s:        "abc",
			others:   []string{"", "ab", "a b c", "abc"},
			expected: "true",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r, err := StrEqualAny(nil, test.s, test.others...)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, r)
		})
	}
}

func TestSubstring(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		startIndex string
		lengthStr  string
		expected   string
		err        string
	}{
		{
			name:       "invalid startIndex",
			str:        "123456",
			startIndex: "abc",
			lengthStr:  "5",
			expected:   "",
			err:        `unable to convert start index 'abc' into int, err: strconv.Atoi: parsing "abc": invalid syntax`,
		},
		{
			name:       "invalid lengthStr",
			str:        "123456",
			startIndex: "5",
			lengthStr:  "abc",
			expected:   "",
			err:        `unable to convert length 'abc' into int, err: strconv.Atoi: parsing "abc": invalid syntax`,
		},
		{
			name:       "empty startIndex",
			str:        "123456",
			startIndex: "",
			lengthStr:  "5",
			expected:   "",
			err:        `unable to convert start index '' into int, err: strconv.Atoi: parsing "": invalid syntax`,
		},
		{
			name:       "empty lengthStr",
			str:        "123456",
			startIndex: "5",
			lengthStr:  "",
			expected:   "",
			err:        `unable to convert length '' into int, err: strconv.Atoi: parsing "": invalid syntax`,
		},
		{
			name:       "empty str",
			str:        "",
			startIndex: "0",
			lengthStr:  "0",
			expected:   "",
			err:        "",
		},
		{
			name:       "empty str with non-0 startIndex",
			str:        "",
			startIndex: "1",
			lengthStr:  "0",
			expected:   "",
			err:        `start index 1 is out of bounds (string length is 0)`,
		},
		{
			name:       "empty str with non-0 lengthStr",
			str:        "",
			startIndex: "0",
			lengthStr:  "1",
			expected:   "",
			err:        `start 0 + length 1 is out of bounds (string length is 0)`,
		},
		{
			name:       "0 startIndex",
			str:        "123456",
			startIndex: "0",
			lengthStr:  "4",
			expected:   "1234",
			err:        "",
		},
		{
			name:       "lengthStr is 1",
			str:        "123456",
			startIndex: "4",
			lengthStr:  "1",
			expected:   "5",
			err:        "",
		},
		{
			name:       "lengthStr is 0",
			str:        "123456",
			startIndex: "1",
			lengthStr:  "0",
			expected:   "",
			err:        "",
		},
		{
			name:       "lengthStr is -1",
			str:        "123456",
			startIndex: "3",
			lengthStr:  "-1",
			expected:   "456",
			err:        "",
		},
		{
			name:       "negative startIndex",
			str:        "123456",
			startIndex: "-4",
			lengthStr:  "4",
			expected:   "",
			err:        `start index -4 is out of bounds (string length is 6)`,
		},
		{
			name:       "negative lengthStr other than -1",
			str:        "123456",
			startIndex: "4",
			lengthStr:  "-2",
			expected:   "",
			err:        `length must be >= -1, but got -2`,
		},
		{
			name:       "out-of-bounds startIndex",
			str:        "123456",
			startIndex: "9",
			lengthStr:  "2",
			expected:   "",
			err:        `start index 9 is out of bounds (string length is 6)`,
		},
		{
			name:       "out-of-bounds lengthStr",
			str:        "123456",
			startIndex: "2",
			lengthStr:  "7",
			expected:   "",
			err:        `start 2 + length 7 is out of bounds (string length is 6)`,
		},
		{
			name:       "out-of-bounds startIndex and lengthStr",
			str:        "123456",
			startIndex: "10",
			lengthStr:  "9",
			expected:   "",
			err:        `start index 10 is out of bounds (string length is 6)`,
		},
		{
			name:       "substring starts at the beginning",
			str:        "123456",
			startIndex: "0",
			lengthStr:  "4",
			expected:   "1234",
			err:        "",
		},
		{
			name:       "substring ends at the end",
			str:        "123456",
			startIndex: "2",
			lengthStr:  "4",
			expected:   "3456",
			err:        "",
		},
		{
			name:       "substring starts at the end",
			str:        "123456",
			startIndex: "6",
			lengthStr:  "0",
			expected:   "",
			err:        "",
		},
		{
			name:       "substring ends at the beginning",
			str:        "123456",
			startIndex: "0",
			lengthStr:  "0",
			expected:   "",
			err:        "",
		},
		{
			name:       "substring is the whole string",
			str:        "123456",
			startIndex: "0",
			lengthStr:  "6",
			expected:   "123456",
			err:        "",
		},
		{
			name:       "non-ASCII string",
			str:        "ü:ü",
			startIndex: "1",
			lengthStr:  "2",
			expected:   ":ü",
			err:        "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Substring(nil, test.str, test.startIndex, test.lengthStr)
			if test.err == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			} else {
				assert.Error(t, err)
				assert.Equal(t, test.err, err.Error())
				assert.Equal(t, "", result)
			}
		})
	}
}

func TestSwitchFunc(t *testing.T) {
	for _, test := range []struct {
		name         string
		expr         string
		casesReturns []string
		err          string
		expected     string
	}{
		{
			name:         "empty casesReturns",
			expr:         "abc",
			casesReturns: nil,
			err:          "length of 'casesAndValues' must be odd, but got: 0",
		},
		{
			name:         "even casesReturns length",
			expr:         "abc",
			casesReturns: []string{"1", "2", "3", "4"},
			err:          "length of 'casesAndValues' must be odd, but got: 4",
		},
		{
			name:         "no case, just default",
			expr:         "abc",
			casesReturns: []string{"default"},
			expected:     "default",
		},
		{
			name: "case string contains special characters",
			expr: "How do you do",
			casesReturns: []string{
				"How do you do?", "Wrong",
				"How do you do", "Correct",
				"Huh"},
			expected: "Correct",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := SwitchFunc(nil, test.expr, test.casesReturns...)
			if test.err != "" {
				assert.Error(t, err)
				assert.Equal(t, test.err, err.Error())
				assert.Equal(t, "", result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestSwitchByPattern(t *testing.T) {
	for _, test := range []struct {
		name            string
		expr            string
		patternsReturns []string
		err             string
		expected        string
	}{
		{
			name:            "empty patternsReturns",
			expr:            "abc",
			patternsReturns: nil,
			err:             "length of 'patternsAndValues' must be odd, but got: 0",
		},
		{
			name:            "even patternsReturns length",
			expr:            "abc",
			patternsReturns: []string{"1", "2", "3", "4"},
			err:             "length of 'patternsAndValues' must be odd, but got: 4",
		},
		{
			name:            "regex invalid",
			expr:            "abc",
			patternsReturns: []string{"[", "2", "3"},
			err:             "invalid pattern '[', err: error parsing regexp: missing closing ]: `[`",
		},
		{
			name:            "no pattern, only default",
			expr:            "abc",
			patternsReturns: []string{"default"},
			expected:        "default",
		},
		{
			name: "case string contains special characters",
			expr: "2019/02/23",
			patternsReturns: []string{
				"^[0-9]{2}/[0-9]{2}/[0-9]{4}$", "Wrong",
				"^[0-9]{4}/[0-9]{2}/[0-9]{2}$", "Correct",
				"Huh"},
			expected: "Correct",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := SwitchByPattern(nil, test.expr, test.patternsReturns...)
			if test.err != "" {
				assert.Error(t, err)
				assert.Equal(t, test.err, err.Error())
				assert.Equal(t, "", result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
