package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitIntoStationNameAndTemperature(t *testing.T) {
	line := []byte("Odesa;10")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "Odesa", string(stationName))
	assert.Equal(t, "10", string(temperature))
}

func TestSplitAnInvalidLine(t *testing.T) {
	line := []byte("Odesa:10")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Error(t, err)
	assert.Equal(t, []byte(nil), stationName)
	assert.Equal(t, []byte(nil), temperature)
}
