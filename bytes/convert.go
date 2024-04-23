package bytes

import (
	"errors"
	"fmt"
	"strconv"
)

var errFloatParse = errors.New("cannot parse float64")

const (
	fractionPartRepresentation = float64(10)
	minusSign                  = byte('-')
)

// ToTemperature converts the input to a float64 representation.
// It requires . to be present, will fail if . is missing.
// It does not consider + or infinity symbol.
func ToTemperature(input []byte) (float64, error) {
	return convert(input)
}

func Format(temperature float64, slice []byte) string {
	appended := strconv.AppendFloat(slice[:], temperature, 'f', -1, 64)
	return string(appended)
}

func convert(input []byte) (float64, error) {
	var asFloat float64

	minus := input[0] == minusSign
	integerValue, currentIndex := integerPart(input, minus)
	if input[currentIndex] == '.' {
		currentIndex++
		fractionalValue := input[currentIndex] - '0'
		floatValue := uint16(integerValue)*10 + uint16(fractionalValue)

		asFloat = float64(floatValue) / fractionPartRepresentation
		if minus {
			asFloat = -asFloat
		}
		return asFloat, nil
	}
	return 0, fmt.Errorf("%v, input %s", errFloatParse, input)
}

func integerPart(input []byte, minus bool) (uint8, uint) {
	currentIndex := uint(0)
	if minus {
		currentIndex++
	}
	integerValue := uint8(0)
	for input[currentIndex] != '.' {
		integerValue = integerValue*10 + (input[currentIndex] - '0')
		currentIndex++
	}
	return integerValue, currentIndex
}
