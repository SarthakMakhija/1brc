package bytes

import (
	"errors"
)

const (
	separator = byte(';')
	minusSign = byte('-')
)

var ErrInvalidLineFormat = errors.New("invalid line format")

var temperatureMultiplier = [2]int16{1, 10}

// SplitIntoStationNameAndTemperature expects a valid line of the format:
// StationName;Temperature.
// Temperature must have a single digit after . (dot).
// It does not handle any other case.
func SplitIntoStationNameAndTemperature(line []byte) ([]byte, Temperature, error) {
	lineLength := len(line)

	fractionalValue := int16(line[lineLength-1] - '0')
	integerValue := int16(0)
	temperatureDigitIndex := 0

	minus := false
	for index := lineLength - 3; index >= 0; index-- {
		switch ch := line[index]; ch {
		case separator:
			eligibleTemperature := integerValue*10 + (fractionalValue)
			if minus {
				eligibleTemperature = -eligibleTemperature
			}
			return line[:index], eligibleTemperature, nil
		case minusSign:
			minus = true
		default:
			multiplier := temperatureMultiplier[temperatureDigitIndex]
			integerValue = integerValue + int16(ch-'0')*multiplier
			temperatureDigitIndex++
		}
	}
	return nil, -1, ErrInvalidLineFormat
}
