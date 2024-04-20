package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitIntoStationNameAndTemperature(t *testing.T) {
	line := []byte("Odesa;10")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "Odesa", stationName)
	assert.Equal(t, "10", temperature)
}

func TestSplitIntoStationNameAndNegativeTemperature(t *testing.T) {
	line := []byte("Odesa;-10.45")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "Odesa", stationName)
	assert.Equal(t, "-10.45", temperature)
}

func TestSplitAnInvalidLine(t *testing.T) {
	line := []byte("Odesa:10")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Error(t, err)
	assert.Equal(t, "", stationName)
	assert.Equal(t, "", temperature)
}
