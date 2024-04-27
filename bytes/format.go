package bytes

import (
	"strconv"
)

type Temperature = int16

func Format(temperature Temperature, slice []byte) []byte {
	return strconv.AppendFloat(slice[:], float64(TemperatureAsFloat32(temperature)), 'f', 1, 32)
}

func FormatTemperatureAsFloat32(temperature float32, slice []byte) []byte {
	return strconv.AppendFloat(slice[:], float64(temperature), 'f', 1, 32)
}

func TemperatureAsFloat32(temperature Temperature) float32 {
	return float32(temperature) * 0.1
}
