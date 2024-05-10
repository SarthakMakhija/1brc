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

func TestSplitIntoStationNameAndTemperatureWithoutDelimiter(t *testing.T) {
	line := []byte("São Paulo:1.4")
	_, _, err := SplitIntoStationNameAndTemperature(line)

	assert.Error(t, err)
}

func TestSplitIntoStationNameAndTemperatureWithLessThanMinLineLengthRequired(t *testing.T) {
	line := []byte("S")
	_, _, err := SplitIntoStationNameAndTemperature(line)

	assert.Error(t, err)
}

func TestSplitIntoStationNameAndTemperatureWithTemperatureNotInRange(t *testing.T) {
	line := []byte("São Paulo;111.4")
	_, _, err := SplitIntoStationNameAndTemperature(line)

	assert.Error(t, err)
}

func TestConvertsTheSliceToTemperature(t *testing.T) {
	line := []byte("10.1")
	temperature := ToTemperature(line)

	assert.Equal(t, Temperature(101), temperature)
}

func TestConvertsTheSliceToNegativeTemperature(t *testing.T) {
	line := []byte("-99.9")
	temperature := ToTemperature(line)

	assert.Equal(t, Temperature(-999), temperature)
}

func TestConvertsTheSliceToTemperatureLessThan10(t *testing.T) {
	line := []byte("1.4")
	temperature := ToTemperature(line)

	assert.Equal(t, Temperature(14), temperature)
}

func TestConvertsTheSliceToNegativeTemperatureLessThan10(t *testing.T) {
	line := []byte("-1.4")
	temperature := ToTemperature(line)

	assert.Equal(t, Temperature(-14), temperature)
}

func TestConvertsTheSliceWithNewLineToTemperature(t *testing.T) {
	line := []byte("10.1\n")
	temperature, numberOfBytesRead := ToTemperatureWithNewLine(line)

	assert.Equal(t, 4, numberOfBytesRead)
	assert.Equal(t, Temperature(101), temperature)
}

func TestConvertsTheSliceWithNewLineToNegativeTemperature(t *testing.T) {
	line := []byte("-99.9\n")
	temperature, numberOfBytesRead := ToTemperatureWithNewLine(line)

	assert.Equal(t, 5, numberOfBytesRead)
	assert.Equal(t, Temperature(-999), temperature)
}

func TestConvertsTheSliceWithNewLineToTemperatureLessThan10(t *testing.T) {
	line := []byte("1.4\n")
	temperature, numberOfBytesRead := ToTemperatureWithNewLine(line)

	assert.Equal(t, 3, numberOfBytesRead)
	assert.Equal(t, Temperature(14), temperature)
}

func TestConvertsTheSliceWithNewLineToNegativeTemperatureLessThan10(t *testing.T) {
	line := []byte("-1.4\n")
	temperature, numberOfBytesRead := ToTemperatureWithNewLine(line)

	assert.Equal(t, 4, numberOfBytesRead)
	assert.Equal(t, Temperature(-14), temperature)
}
