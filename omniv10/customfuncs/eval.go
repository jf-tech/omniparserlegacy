package customfuncs

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/jf-tech/go-corelib/caches"
	"github.com/jf-tech/go-corelib/strs"
	"github.com/jf-tech/omniparser/transformctx"
)

var exprCache = caches.NewLoadingCache()

const (
	argTypeString  = "string"
	argTypeInt     = "int"
	argTypeFloat   = "float"
	argTypeBoolean = "boolean"
)

func constructEvalParam(argDecl, argValue string) (name string, value interface{}, err error) {
	declParts := strings.Split(argDecl, ":")
	if len(declParts) != 2 {
		return "", nil, errors.New("arg decl must be in format of '<arg_name>:<arg_type>'")
	}
	name = declParts[0]
	if !strs.IsStrNonBlank(name) {
		return "", nil, errors.New("arg_name in '<arg_name>:<arg_type>' cannot be empty/whitespace string")
	}
	switch declParts[1] {
	case argTypeString:
		return name, argValue, nil
	case argTypeInt:
		f, err := strconv.ParseFloat(argValue, 64)
		if err != nil {
			return "", nil, err
		}
		return name, int64(f), nil
	case argTypeFloat:
		f, err := strconv.ParseFloat(argValue, 64)
		if err != nil {
			return "", nil, err
		}
		return name, f, nil
	case argTypeBoolean:
		b, err := strconv.ParseBool(argValue)
		if err != nil {
			return "", nil, err
		}
		return name, b, nil
	default:
		return "", nil, fmt.Errorf("arg_type '%s' in '<arg_name>:<arg_type>' is not supported", declParts[1])
	}
}

// Eval evaluates a given expression with input params and return the result in string.
// For supported expression formats and operators, check:
// https://github.com/Knetic/govaluate/blob/master/MANUAL.md
func Eval(_ *transformctx.Ctx, exprStr string, args ...string) (string, error) {
	if len(args)%2 != 0 {
		return "", errors.New("invalid number of args to 'eval'")
	}
	params := make(map[string]interface{}, len(args)/2)
	for i := 0; i < len(args)/2; i++ {
		n, v, err := constructEvalParam(args[i*2], args[i*2+1])
		if err != nil {
			return "", err
		}
		params[n] = v
	}
	expr, err := exprCache.Get(exprStr, func(key interface{}) (interface{}, error) {
		return govaluate.NewEvaluableExpression(exprStr)
	})
	if err != nil {
		return "", err
	}
	result, err := expr.(*govaluate.EvaluableExpression).Evaluate(params)
	if err != nil {
		return "", err
	}
	switch {
	case result == nil:
		return "", nil
	default:
		return fmt.Sprint(result), nil
	}
}
