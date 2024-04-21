package bytes

import (
	"errors"
	"fmt"
)

var errFloatParse = errors.New("cannot parse float64")

const fractionPartRepresentation = float64(10)

// toFloat64 simple to float64 conversion.
// It requires . to be present, will fail if . is missing.
// It does not consider + or infinity symbol.
func toFloat64(input []byte) (float64, error) {
	return convert(input)
}

func convert(input []byte) (float64, error) {
	var asFloat float64
	minus := input[0] == '-'

	wholeValue, currentIndex := integerPart(input)
	if input[currentIndex] == '.' {
		currentIndex++
		wholeValue = wholeValue*10 + uint64(input[currentIndex]-'0')

		asFloat = float64(wholeValue) / fractionPartRepresentation
		if minus {
			asFloat = -asFloat
		}
		return asFloat, nil
	}
	return 0, fmt.Errorf("%v, input %s", errFloatParse, input)
}

func integerPart(input []byte) (uint64, uint) {
	currentIndex := uint(0)
	minus := input[0] == '-'
	if minus {
		currentIndex++
	}

	wholeValue := uint64(0)
	for input[currentIndex] != '.' {
		wholeValue = wholeValue*10 + uint64(input[currentIndex]-'0')
		currentIndex++
	}
	return wholeValue, currentIndex
}
