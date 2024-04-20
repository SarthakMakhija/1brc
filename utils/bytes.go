package utils

import "errors"

var separator = byte(';')

var ErrInvalidLineFormat = errors.New("invalid line format")

func SplitIntoStationNameAndTemperature(line []byte) ([]byte, []byte, error) {
	for index := 0; index < len(line); index++ {
		if line[index] == separator {
			return line[:index], line[index+1:], nil
		}
	}
	return nil, nil, ErrInvalidLineFormat
}
