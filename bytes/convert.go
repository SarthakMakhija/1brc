package bytes

import (
	"errors"
	"strconv"
)

var errParseTemperature = errors.New("cannot parse temperature")

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
	fractionalValue := input[len(input)-1] //bound check eliminated further in the code

	minus := input[0] == minusSign
	inputSlice := input
	if minus {
		inputSlice = input[1:]
	}

	integerValue := integerPart(inputSlice)
	fractionalValue = fractionalValue - '0'
	eligibleFloat := uint16(integerValue)*10 + uint16(fractionalValue)

	asTemperature := float64(eligibleFloat) / fractionPartRepresentation
	if minus {
		asTemperature = -asTemperature
	}
	return asTemperature, nil
}

func integerPart(input []byte) uint8 {
	integerValue := uint8(0)

	for _, ch := range input {
		if ch == '.' {
			return integerValue
		}
		integerValue = integerValue*10 + (ch - '0')
	}
	return integerValue
}
