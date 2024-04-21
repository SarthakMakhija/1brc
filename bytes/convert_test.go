package bytes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToFloat64(t *testing.T) {
	inputTable := []struct {
		input    string
		expected float64
	}{
		{input: "10.43", expected: 10.43},
		{input: "19.4333", expected: 19.4333},
		{input: "-34.5997", expected: -34.5997},
		{input: "40.6943", expected: 40.6943},
		{input: "40.6", expected: 40.6},
		{input: "-40.1", expected: -40.1},
	}

	for _, row := range inputTable {
		result, err := toFloat64(row.input)
		assert.Nil(t, err)
		assert.Equal(t, row.expected, result)
	}
}
