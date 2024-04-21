package bytes

import (
	"errors"
	"unsafe"
)

var separator = byte(';')

var ErrInvalidLineFormat = errors.New("invalid line format")

func SplitIntoStationNameAndTemperature(line []byte) (string, float64, error) {
	stationName := func(endIndex int) []byte {
		station := make([]byte, endIndex)
		copy(station, line[:endIndex])
		return station
	}

	for index := 0; index < len(line); index++ {
		if line[index] == separator {
			temperature, err := toFloat64(line[index+1:])
			return zeroCopyString(stationName(index)), temperature, err
		}
	}
	return "", 0, ErrInvalidLineFormat
}

func zeroCopyString(slice []byte) string {
	return *(*string)(unsafe.Pointer(&slice))
}
