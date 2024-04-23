package bytes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToFloat64(t *testing.T) {
	inputTable := []struct {
		input    []byte
		expected float64
	}{
		{input: []byte("10.4"), expected: 10.4},
		{input: []byte("19.7"), expected: 19.7},
		{input: []byte("-34.1"), expected: -34.1},
		{input: []byte("40.6"), expected: 40.6},
		{input: []byte("40.9"), expected: 40.9},
		{input: []byte("-40.1"), expected: -40.1},
	}

	for _, row := range inputTable {
		result, err := toFloat64(row.input)
		assert.Nil(t, err)
		assert.Equal(t, row.expected, result)
	}
}

func TestFormat(t *testing.T) {
	inputTable := []struct {
		input    float64
		expected string
	}{
		{input: 10.4, expected: "10.4"},
		{input: 19.7, expected: "19.7"},
		{input: -34.1, expected: "-34.1"},
		{input: 40.6, expected: "40.6"},
		{input: 40.9, expected: "40.9"},
		{input: -40.1, expected: "-40.1"},
	}

	temperatureBytes := make([]byte, 0, 64)
	for _, row := range inputTable {
		result := Format(row.input, temperatureBytes)
		assert.Equal(t, row.expected, result)
	}
}
