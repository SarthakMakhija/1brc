package bytes

import (
	"errors"
)

const (
	separator  = byte(';')
	minusSign  = byte('-')
	multiplier = int16(10)
)

var ErrInvalidLineFormat = errors.New("invalid line format")

// SplitIntoStationNameAndTemperature expects a valid line of the format:
// StationName;Temperature.
// Temperature must have a single digit after . (dot).
// It does not handle any other case.
func SplitIntoStationNameAndTemperature(line []byte) ([]byte, Temperature, error) {
	lineLength := len(line)

	fractionalValue := int16(line[lineLength-1] - '0') //lineLength-1 represents the fractional digit.
	integerValue := int16(line[lineLength-3] - '0')    //lineLength-3 represents the lowest position temperature digit.

	minus := false
	for index := lineLength - 4; index >= 0; index-- {
		switch ch := line[index]; ch {
		case minusSign:
			minus = true
		case separator:
			eligibleTemperature := integerValue*10 + (fractionalValue)
			if minus {
				eligibleTemperature = -eligibleTemperature
			}
			return line[:index], eligibleTemperature, nil
		default:
			integerValue = integerValue + int16(ch-'0')*multiplier
		}
	}
	return nil, -1, ErrInvalidLineFormat
}
