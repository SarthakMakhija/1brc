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
	line := []byte("Odesa;-10.4")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "Odesa", string(stationName))
	assert.Equal(t, Temperature(-104), temperature)
}

func TestSplitIntoStationNameAndNegativeTemperatureWithStationNameHavingSpecialCharacters(t *testing.T) {
	line := []byte("São Paulo;-10.4")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "São Paulo", string(stationName))
	assert.Equal(t, Temperature(-104), temperature)
}
