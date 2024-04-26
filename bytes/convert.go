package bytes

import (
	"strconv"
)

const (
	minusSign = byte('-')
)

type Temperature = int16

// ToTemperature converts the input to a Temperature representation.
// It requires . to be present, will fail if . is missing.
// It does not consider + or infinity symbol.
func ToTemperature(input []byte) (Temperature, error) {
	return convert(input)
}

func Format(temperature Temperature, slice []byte) []byte {
	return strconv.AppendFloat(slice[:], float64(TemperatureAsFloat32(temperature)), 'f', 1, 32)
}

func FormatTemperatureAsFloat32(temperature float32, slice []byte) []byte {
	return strconv.AppendFloat(slice[:], float64(temperature), 'f', 1, 32)
}

func TemperatureAsFloat32(temperature Temperature) float32 {
	return float32(temperature) * 0.1
}

func convert(input []byte) (Temperature, error) {
	fractionalValue := input[len(input)-1] //bound check eliminated further in the code

	minus := input[0] == minusSign
	inputSlice := input
	if minus {
		inputSlice = input[1:]
	}

	integerValue := integerPart(inputSlice)
	fractionalValue = fractionalValue - '0'
	eligibleTemperature := int16(integerValue)*10 + int16(fractionalValue)

	if minus {
		eligibleTemperature = -eligibleTemperature
	}
	return eligibleTemperature, nil
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
