package customfuncs

import (
	"fmt"
	"strconv"

	"github.com/jf-tech/omniparser/transformctx"
)

// Avg computes the average of the specified 'values'. If no 'values' are given, it returns '0'.
func Avg(_ *transformctx.Ctx, values ...string) (string, error) {
	if len(values) == 0 {
		return "0", nil
	}
	s := float64(0)
	for _, v := range values {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return "", err
		}
		s += f
	}
	return fmt.Sprintf("%v", s/float64(len(values))), nil
}

// Sum computes the sum of the specified 'values'. If no 'values' are given, it returns '0'.
func Sum(_ *transformctx.Ctx, values ...string) (string, error) {
	s := float64(0)
	for _, v := range values {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return "", err
		}
		s += f
	}
	return fmt.Sprintf("%v", s), nil
}
