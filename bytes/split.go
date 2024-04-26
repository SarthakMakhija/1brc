package bytes

import (
	"errors"
)

const separator = byte(';')

var ErrInvalidLineFormat = errors.New("invalid line format")

func SplitIntoStationNameAndTemperature(line []byte) ([]byte, float64, error) {
	lineLength := len(line)
	for index := lineLength - 1; index >= 0; index-- {
		if line[index] == separator {
			temperature, err := ToTemperature(line[index+1:])
			return line[:index], temperature, err
		}
	}
	return nil, 0, ErrInvalidLineFormat
}
