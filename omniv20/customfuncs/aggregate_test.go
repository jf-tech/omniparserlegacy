package customfuncs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvg(t *testing.T) {
	for _, test := range []struct {
		name     string
		values   []string
		err      string
		expected string
	}{
		{
			name:     "nil",
			values:   nil,
			err:      "",
			expected: "0",
		},
		{
			name:     "empty",
			values:   []string{},
			err:      "",
			expected: "0",
		},
		{
			name:     "single",
			values:   []string{"3.14159265358"},
			err:      "",
			expected: "3.14159265358",
		},
		{
			name:     "multiple small ones",
			values:   []string{"3.45", "5.38"},
			err:      "",
			expected: "4.415",
		},
		{
			name:     "multiple big ones",
			values:   []string{"1.23e+9", "0.34E+10"},
			err:      "",
			expected: "2.315e+09",
		},
		{
			name:     "invalid value",
			values:   []string{"1", "two"},
			err:      `strconv.ParseFloat: parsing "two": invalid syntax`,
			expected: "",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := Avg(nil, test.values...)
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

func TestSum(t *testing.T) {
	for _, test := range []struct {
		name     string
		values   []string
		err      string
		expected string
	}{
		{
			name:     "nil",
			values:   nil,
			err:      "",
			expected: "0",
		},
		{
			name:     "empty",
			values:   []string{},
			err:      "",
			expected: "0",
		},
		{
			name:     "single",
			values:   []string{"3.14159265358"},
			err:      "",
			expected: "3.14159265358",
		},
		{
			name:     "multiple small ones",
			values:   []string{"3.45", "5.38"},
			err:      "",
			expected: "8.83",
		},
		{
			name:     "multiple big ones",
			values:   []string{"1.23e+9", "0.34E+10"},
			err:      "",
			expected: "4.63e+09",
		},
		{
			name:     "invalid value",
			values:   []string{"1", "two"},
			err:      `strconv.ParseFloat: parsing "two": invalid syntax`,
			expected: "",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			result, err := Sum(nil, test.values...)
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
