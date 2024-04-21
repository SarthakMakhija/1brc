package bytes

import (
	"errors"
	"fmt"
)

var errFloatParse = errors.New("cannot parse float64")

var float64pow10 = [...]float64{
	1e0, 1e1, 1e2, 1e3, 1e4, 1e5, 1e6, 1e7, 1e8, 1e9, 1e10, 1e11, 1e12, 1e13, 1e14, 1e15, 1e16,
}

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
		dotIndex := currentIndex

		inputLength := uint(len(input))
		for currentIndex < inputLength {
			wholeValue = wholeValue*10 + uint64(input[currentIndex]-'0')
			currentIndex++
		}

		asFloat = float64(wholeValue) / float64pow10[currentIndex-dotIndex]
		if currentIndex >= inputLength {
			if minus {
				asFloat = -asFloat
			}
			return asFloat, nil
		}
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
