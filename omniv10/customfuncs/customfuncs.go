package customfuncs

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/jf-tech/go-corelib/caches"
	"github.com/jf-tech/go-corelib/strs"
	"github.com/jf-tech/omniparser/customfuncs"
	"github.com/jf-tech/omniparser/transformctx"

	"github.com/jf-tech/omniparserlegacy/omniv10/errs"
)

// OmniV10OnlyCustomFuncs contains omniparser '1.0' specific custom funcs.
var OmniV10OnlyCustomFuncs = map[string]customfuncs.CustomFuncType{
	// keep these custom funcs lexically sorted
	"containsPattern":             ContainsPattern,
	"dateTimeToRfc3339":           customfuncs.DateTimeToRFC3339,
	"dateTimeWithLayoutToRfc3339": customfuncs.DateTimeLayoutToRFC3339,
	"epochToDateTimeRfc3339":      customfuncs.EpochToDateTimeRFC3339,
	"eval":                        Eval,
	"external":                    External,
	"floor":                       Floor,
	"ifElse":                      IfElse,
	"isEmpty":                     IsEmpty,
	"replace":                     Replace,
	"retrieveBySplit":             RetrieveBySplit,
	"rowSkip":                     RowSkip,
	"rsubstring":                  RSubstring,
	"splitIntoJsonArray":          SplitIntoJSONArray,
	"strEqualAny":                 StrEqualAny,
	"substring":                   Substring,
	"switch":                      SwitchFunc,
	"switchByPattern":             SwitchByPattern,
}

// OmniV10CustomFuncs contains all custom funcs supported by omniparser '1.0'.
var OmniV10CustomFuncs = customfuncs.Merge(
	customfuncs.CommonCustomFuncs,
	OmniV10OnlyCustomFuncs)

// ContainsPattern checks if any of 'strs' matches the given regex pattern, returns "true" if yes, and "false" if no.
func ContainsPattern(_ *transformctx.Ctx, regexPattern string, strs ...string) (string, error) {
	r, err := caches.GetRegex(regexPattern)
	if err != nil {
		return "", err
	}
	for _, str := range strs {
		if r.MatchString(str) {
			return "true", nil
		}
	}
	return "false", nil
}

// External looks up an external property value by name and fails if not found.
func External(ctx *transformctx.Ctx, name string) (string, error) {
	if v, found := ctx.ExternalProperties[name]; found {
		return v, nil
	}
	return "", fmt.Errorf("cannot find external property '%s'", name)
}

// Floor parses 'valueStr' in as a float64 and then floors/truncates it to the specified decimal places.
func Floor(_ *transformctx.Ctx, valueStr, decimalPlaces string) (string, error) {
	v, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return "", fmt.Errorf("unable to parse value '%s' to float64: %s", valueStr, err.Error())
	}
	dp, err := strconv.Atoi(decimalPlaces)
	if err != nil {
		return "", fmt.Errorf("unable to parse decimal place value '%s' to int: %s", decimalPlaces, err.Error())
	}
	if dp < 0 || dp > 100 {
		return "", fmt.Errorf("decimal place value must be an integer with range of [0,100], instead, got %d", dp)
	}
	p10 := math.Pow10(dp)
	return fmt.Sprintf("%v", math.Floor(v*p10)/p10), nil
}

// IfElse does a 'if-elseif-elseif-...-else' pattern check and returns the value that meets the condition.
// Note 'conditionsAndValues' must be of odd length, each of the even values is a condition value and it should
// be a boolean string and each of the odd values is a value to be returned, should its corresponding condition
// value is "true".
func IfElse(_ *transformctx.Ctx, conditionsAndValues ...string) (string, error) {
	if len(conditionsAndValues)%2 != 1 {
		return "", fmt.Errorf("arg number must be odd, but got: %d", len(conditionsAndValues))
	}
	for i := 0; i < len(conditionsAndValues)/2; i++ {
		condition, err := strconv.ParseBool(conditionsAndValues[2*i])
		if err != nil {
			return "", fmt.Errorf(
				`condition argument must be a boolean string, but got: %s`, conditionsAndValues[2*i])
		}
		if condition {
			return conditionsAndValues[(2*i)+1], nil
		}
	}
	return conditionsAndValues[len(conditionsAndValues)-1], nil
}

// IsEmpty returns "true" if 'str' is empty; false if 'str' isn't. Note blank string (with whitespace
// only) is not considered empty.
func IsEmpty(_ *transformctx.Ctx, str string) (string, error) {
	if str == "" {
		return "true", nil
	}
	return "false", nil
}

// Replace replaces in 's' all the substrings that matches 'regexStr' with 'replStr'.
func Replace(_ *transformctx.Ctx, s, regexStr, replStr string) (string, error) {
	r, err := caches.GetRegex(regexStr)
	if err != nil {
		return "", err
	}
	return r.ReplaceAllString(s, replStr), nil
}

// RetrieveBySplit splits a 'str' by 'sep' and returns the 'indexStr'-th part. Note 'indexStr' is 0-based.
func RetrieveBySplit(_ *transformctx.Ctx, str, sep, indexStr string) (string, error) {
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return "", err
	}
	parts := strings.Split(str, sep)
	if index < 0 || index >= len(parts) {
		return "", fmt.Errorf(
			"'index' is out of range. original string: %s, separator: %s, index: %d", str, sep, index)
	}
	return parts[index], nil
}

// RowSkip checks if a 'value' matches 'skipRegex' pattern, and return a special error code NonErrorRecordSkipped
// to indicate to omniparser '1.0' to skip processing a record, if the pattern matches.
func RowSkip(_ *transformctx.Ctx, value, skipRegex string) (string, error) {
	r, err := caches.GetRegex(skipRegex)
	if err != nil {
		return "", err
	}
	if r.Match([]byte(value)) {
		return "", errs.NonErrorRecordSkipped(value)
	}
	return "", nil
}

// RSubstring extracts a part of a 'str' from 'startIndex' from the right end of the string with length equal
// to 'lengthStr'. 'startIndex' must be [0, len). 'lengthStr' can be '-1', in which case RSubstring will return
// all the remaining characters (towards the head) of 'str' starting from 'startIndex' from the right. Or
// 'lengthStr' can be >= 0, in which case RSubstring will only return 'lengthStr' number characters of 'str'
// starting from 'startIndex' from the right.
func RSubstring(ctx *transformctx.Ctx, str, startIndex, lengthStr string) (string, error) {
	// not the most efficient way of doing thing, but again, perf isn't the most critical here so keep it
	// simple.
	reverse := func(r []rune) []rune {
		l := len(r)
		for i := 0; i < l/2; i++ {
			r[i], r[l-i-1] = r[l-i-1], r[i]
		}
		return r
	}
	r := reverse([]rune(str))
	s, err := Substring(ctx, string(r), startIndex, lengthStr)
	if err != nil {
		return "", err
	}
	return string(reverse([]rune(s))), nil
}

// SplitIntoJSONArray splits an 's' by 'sep' and returns the resulting array of parts in a JSON string.
// An optional 'trim' argument can be specified to tell the function to space trim each part or not.
func SplitIntoJSONArray(_ *transformctx.Ctx, s, sep string, trim ...string) (string, error) {
	if len(trim) > 1 {
		return "", fmt.Errorf("cannot specify 'trim' argument more than once")
	}
	toTrim := false
	if len(trim) == 1 {
		b, err := strconv.ParseBool(trim[0])
		if err != nil {
			return "", fmt.Errorf("expect 'trim' to be of boolean string, but got: '%s'", trim[0])
		}
		toTrim = b
	}
	if s == "" {
		return "[]", nil
	}
	splits := strings.Split(s, sep)
	if toTrim {
		splits = strs.NoErrMapSlice(splits, func(s string) string {
			return strings.TrimSpace(s)
		})
	}
	bytes, _ := json.Marshal(splits)
	return string(bytes), nil
}

// StrEqualAny returns true if the 's' is equal to any of the strings specified by 'others'. If 'others' is empty
// or none of its strings is equal to 's', then "false" is returned.
func StrEqualAny(_ *transformctx.Ctx, s string, others ...string) (string, error) {
	for _, other := range others {
		if s == other {
			return "true", nil
		}
	}
	return "false", nil
}

// Substring extracts a part of a 'str' from 'startIndex' with 'lengthStr'. 'startIndex' must be [0, len).
// 'lengthStr' can be '-1', in which case Substring will return all the remaining characters of 'str' starting
// from 'startIndex'. Or 'lengthStr' can be >= 0, in which case Substring will only return 'lengthStr' number
// characters of 'str' starting from 'startIndex'.
func Substring(_ *transformctx.Ctx, str, startIndex, lengthStr string) (string, error) {
	start, err := strconv.Atoi(startIndex)
	if err != nil {
		return "", fmt.Errorf("unable to convert start index '%s' into int, err: %s", startIndex, err.Error())
	}
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return "", fmt.Errorf("unable to convert length '%s' into int, err: %s", lengthStr, err.Error())
	}
	if length < -1 {
		return "", fmt.Errorf("length must be >= -1, but got %d", length)
	}
	// We can/do deal with UTF-8 encoded strings. startIndex and length are all about
	// UTF-8 characters not just bytes.
	runes := []rune(str)
	runeLen := len(runes)
	if start < 0 || start > runeLen {
		return "", fmt.Errorf("start index %d is out of bounds (string length is %d)", start, runeLen)
	}
	if length == -1 {
		length = runeLen - start
	}
	if start+length > runeLen {
		return "", fmt.Errorf(
			"start %d + length %d is out of bounds (string length is %d)", start, length, runeLen)
	}
	return string(runes[start : start+length]), nil
}

// SwitchFunc tests a 'value' against a number of cases (by string literal comparison) and returns the
// matched value; if nothing matches, the last default value is returned.
func SwitchFunc(ctx *transformctx.Ctx, value string, casesAndValues ...string) (string, error) {
	if len(casesAndValues)%2 != 1 {
		return "", fmt.Errorf("length of 'casesAndValues' must be odd, but got: %d", len(casesAndValues))
	}
	patternsAndValues := strs.CopySlice(casesAndValues[0 : len(casesAndValues)-1])
	for i := 0; i < len(patternsAndValues)/2; i++ {
		patternsAndValues[2*i] = "^" + regexp.QuoteMeta(patternsAndValues[2*i]) + "$"
	}
	return SwitchByPattern(ctx, value, append(patternsAndValues, casesAndValues[len(casesAndValues)-1])...)
}

// SwitchByPattern tests a value against a number of patterns (by regex pattern match) and returns the
// matched value; if nothing matches, the last default value is returned.
func SwitchByPattern(_ *transformctx.Ctx, value string, patternsAndValues ...string) (string, error) {
	if len(patternsAndValues)%2 != 1 {
		return "", fmt.Errorf("length of 'patternsAndValues' must be odd, but got: %d", len(patternsAndValues))
	}
	patternValuePairs := patternsAndValues[0 : len(patternsAndValues)-1]
	for i := 0; i < len(patternValuePairs)/2; i++ {
		pattern, err := caches.GetRegex(patternValuePairs[2*i])
		if err != nil {
			return "", fmt.Errorf(`invalid pattern '%s', err: %s`, patternValuePairs[2*i], err.Error())
		}
		if pattern.MatchString(value) {
			return patternValuePairs[(2*i)+1], nil
		}
	}
	// Use default value when all previous conditions are not met.
	return patternsAndValues[len(patternsAndValues)-1], nil
}
