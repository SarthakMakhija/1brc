package bytes

import (
	"errors"
)

const (
	Separator  = byte(';')
	minusSign  = byte('-')
	multiplier = int16(10)
)

var ErrInvalidLineFormat = errors.New("invalid line format")

// minLineLengthRequired is 5 as we will have:
// a minimum of 3 bytes for the temperature,
// 1 byte for the delimiter,
// 1 byte for the station name
const minLineLengthRequired = 5

// SplitIntoStationNameAndTemperature expects a valid line of the format:
// StationName;Temperature.
// Temperature must have a single digit after . (dot).
// It does not handle any other case.
func SplitIntoStationNameAndTemperature(line []byte) ([]byte, Temperature, error) {
	lineLength := len(line)
	//prevent bound checks
	if lineLength < minLineLengthRequired {
		return nil, 0, ErrInvalidLineFormat
	}

	fractionalValue := int16(line[lineLength-1] - '0') //lineLength-1 represents the fractional digit.
	integerValue := int16(line[lineLength-3] - '0')    //lineLength-3 represents the lowest position temperature digit.

	previousIndex := lineLength - 4
	ch := line[previousIndex]
	switch ch {
	case minusSign:
		eligibleTemperature := integerValue*multiplier + (fractionalValue)
		eligibleTemperature = ^eligibleTemperature + 1
		if lineLength-5 >= 0 { //prevent IsSliceInBounds checks
			return line[:previousIndex-1], eligibleTemperature, nil
		}
		return nil, -1, ErrInvalidLineFormat
	case Separator:
		eligibleTemperature := integerValue*multiplier + (fractionalValue)
		return line[:previousIndex], eligibleTemperature, nil
	default:
		integerValue = integerValue + int16(ch-('0'))*multiplier
	}

	lastIndex := lineLength - 5
	switch line[lastIndex] {
	case minusSign:
		eligibleTemperature := integerValue*multiplier + (fractionalValue)
		eligibleTemperature = ^eligibleTemperature + 1
		if lineLength-6 >= 0 { //prevent IsSliceInBounds checks
			return line[:lastIndex-1], eligibleTemperature, nil
		}
		return nil, -1, ErrInvalidLineFormat
	case Separator:
		eligibleTemperature := integerValue*multiplier + (fractionalValue)
		return line[:lastIndex], eligibleTemperature, nil
	default:
		return nil, -1, ErrInvalidLineFormat
	}
}

// ToTemperature assumes valid input which starts after the Separator (;).
func ToTemperature(slice []byte) Temperature {
	negative := slice[0] == minusSign
	if negative {
		slice = slice[1:]
	}

	var eligibleTemperature Temperature
	switch len(slice) {
	case 3:
		temperature := (slice[0]-'0')*10 + slice[2] - '0'
		eligibleTemperature = Temperature(temperature)
	case 4:
		eligibleTemperature = Temperature(slice[0])*100 + Temperature(slice[1])*10 + Temperature(slice[3]) - '0'*(100+10+1)
	}
	if negative {
		eligibleTemperature = -eligibleTemperature
	}
	return eligibleTemperature
}

func ToTemperatureWithNewLine(slice []byte) (Temperature, int) {
	negative := slice[0] == minusSign
	if negative {
		slice = slice[1:]
	}

	temperature := Temperature(0)
	numberOfBytesRead := 0
	_ = slice[3]

	if slice[1] == '.' {
		temperatureValue := (slice[0]-'0')*10 + (slice[2] - '0')
		temperature = Temperature(temperatureValue)
		numberOfBytesRead = 3
	} else {
		_ = slice[4]
		temperature = Temperature(slice[0])*100 + Temperature(slice[1])*10 + Temperature(slice[3]) - '0'*(100+10+1)
		numberOfBytesRead = 4
	}

	if negative {
		numberOfBytesRead += 1
		temperature = -temperature
	}
	return temperature, numberOfBytesRead
}
