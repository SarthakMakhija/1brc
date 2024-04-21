package bytes

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSplitIntoStationNameAndTemperature(t *testing.T) {
	line := []byte("Odesa;10.1")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "Odesa", stationName)
	assert.Equal(t, 10.1, temperature)
}

func TestSplitIntoStationNameAndNegativeTemperature(t *testing.T) {
	line := []byte("Odesa;-10.45")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Nil(t, err)
	assert.Equal(t, "Odesa", stationName)
	assert.Equal(t, -10.45, temperature)
}

func TestSplitAnInvalidLineBecauseOfInvalidSeparator(t *testing.T) {
	line := []byte("Odesa:10.2")
	stationName, temperature, err := SplitIntoStationNameAndTemperature(line)

	assert.Error(t, err)
	assert.Equal(t, "", stationName)
	assert.Equal(t, float64(0), temperature)
}
