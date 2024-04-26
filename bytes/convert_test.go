package bytes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToToTemperature(t *testing.T) {
	inputTable := []struct {
		input    []byte
		expected Temperature
	}{
		{input: []byte("10.4"), expected: 104},
		{input: []byte("19.7"), expected: 197},
		{input: []byte("-34.1"), expected: -341},
		{input: []byte("40.6"), expected: 406},
		{input: []byte("40.9"), expected: 409},
		{input: []byte("-40.1"), expected: -401},
	}

	for _, row := range inputTable {
		result, err := ToTemperature(row.input)
		assert.Nil(t, err)
		assert.Equal(t, row.expected, result)
	}
}

func TestFormat(t *testing.T) {
	inputTable := []struct {
		input    Temperature
		expected string
	}{
		{input: 104, expected: "10.4"},
		{input: 197, expected: "19.7"},
		{input: -341, expected: "-34.1"},
		{input: 406, expected: "40.6"},
		{input: 409, expected: "40.9"},
		{input: -401, expected: "-40.1"},
	}

	slice := make([]byte, 0, 64)
	for _, row := range inputTable {
		result := Format(row.input, slice)
		assert.Equal(t, row.expected, string(result))
	}
}
