package bytes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitIntoStationNameAndTemperature(t *testing.T) {
	line := []byte("Odesa;10.1")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "Odesa", string(stationName))
	assert.Equal(t, Temperature(101), temperature)
}

func TestSplitIntoStationNameAndNegativeTemperature(t *testing.T) {
	line := []byte("Odesa;-99.9")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "Odesa", string(stationName))
	assert.Equal(t, Temperature(-999), temperature)
}

func TestSplitIntoStationNameAndNegativeTemperatureWithStationNameHavingSpecialCharacters(t *testing.T) {
	line := []byte("São Paulo;-10.4")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", string(stationName))
	assert.Equal(t, Temperature(-104), temperature)
}

func TestSplitIntoStationNameAndPositiveTemperatureLessThan10(t *testing.T) {
	line := []byte("São Paulo;1.4")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", string(stationName))
	assert.Equal(t, Temperature(14), temperature)
}

func TestSplitIntoStationNameAndNegativeTemperatureLessThan10(t *testing.T) {
	line := []byte("São Paulo;-1.4")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", string(stationName))
	assert.Equal(t, Temperature(-14), temperature)
}
