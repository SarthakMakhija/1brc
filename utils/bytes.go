package utils

import (
	"errors"
	"unsafe"
)

var separator = byte(';')

var ErrInvalidLineFormat = errors.New("invalid line format")

func SplitIntoStationNameAndTemperature(line []byte) (string, string, error) {
	for index := 0; index < len(line); index++ {
		if line[index] == separator {
			return zeroCopyString(line[:index]), zeroCopyString(line[index+1:]), nil
		}
	}
	return "", "", ErrInvalidLineFormat
}

func zeroCopyString(slice []byte) string {
	return *(*string)(unsafe.Pointer(&slice))
}
