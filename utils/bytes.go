package utils

import (
	"errors"
	"unsafe"
)

var separator = byte(';')

var ErrInvalidLineFormat = errors.New("invalid line format")

func SplitIntoStationNameAndTemperature(line []byte) (string, string, error) {
	stationName := func(endIndex int) []byte {
		station := make([]byte, endIndex)
		copy(station, line[:endIndex])
		return station
	}

	temperature := func(startIndex int) []byte {
		temperature := make([]byte, len(line)-startIndex)
		copy(temperature, line[startIndex:])
		return temperature
	}

	for index := 0; index < len(line); index++ {
		if line[index] == separator {
			return zeroCopyString(stationName(index)), zeroCopyString(temperature(index + 1)), nil
		}
	}
	return "", "", ErrInvalidLineFormat
}

func zeroCopyString(slice []byte) string {
	return *(*string)(unsafe.Pointer(&slice))
}
