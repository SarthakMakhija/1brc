package utils

import (
	"errors"
)

var separator = byte(';')

var ErrInvalidLineFormat = errors.New("invalid line format")

func SplitIntoStationNameAndTemperature(line []byte) (string, string, error) {
	for index := 0; index < len(line); index++ {
		if line[index] == separator {
			return toString(line[:index]), toString(line[index+1:]), nil
		}
	}
	return "", "", ErrInvalidLineFormat
}

func toString(slice []byte) string {
	return string(slice)
}
