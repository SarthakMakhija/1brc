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
	//prevent bound checks
	if lineLength < 5 {
		return nil, 0, ErrInvalidLineFormat
	}

	fractionalValue := int16(line[lineLength-1] - '0') //lineLength-1 represents the fractional digit.
	integerValue := int16(line[lineLength-3] - '0')    //lineLength-3 represents the lowest position temperature digit.

	ch := line[lineLength-4]
	switch ch {
	case minusSign:
		eligibleTemperature := integerValue*multiplier + (fractionalValue)
		eligibleTemperature = ^eligibleTemperature + 1
		if lineLength-5 >= 0 { //prevent IsSliceInBounds checks
			return line[:lineLength-5], eligibleTemperature, nil
		}
		return nil, -1, ErrInvalidLineFormat
	case separator:
		eligibleTemperature := integerValue*multiplier + (fractionalValue)
		return line[:lineLength-4], eligibleTemperature, nil
	default:
		integerValue = integerValue + int16(ch-('0'))*multiplier
	}

	ch = line[lineLength-5]
	switch ch {
	case minusSign:
		eligibleTemperature := integerValue*multiplier + (fractionalValue)
		eligibleTemperature = ^eligibleTemperature + 1
		if lineLength-6 >= 0 { //prevent IsSliceInBounds checks
			return line[:lineLength-6], eligibleTemperature, nil
		}
		return nil, -1, ErrInvalidLineFormat
	case separator:
		eligibleTemperature := integerValue*multiplier + (fractionalValue)
		return line[:lineLength-5], eligibleTemperature, nil
	default:
		return nil, -1, ErrInvalidLineFormat
	}
}
