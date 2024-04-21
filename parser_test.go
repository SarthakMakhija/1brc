package brc

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseWithTemperaturesForSortedStationsNames(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.23\nMogadishu;5.97\nTirana;15.12\nTirana;9.79\nMogadishu;6.89\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, []string{"Canberra", "Mogadishu", "Odesa", "Tirana"}, result.AllStationsSorted())
}

func TestParseWithTemperaturesForMinTemperature(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.23\nMogadishu;5.97\nTirana;15.12\nTirana;9.79\nMogadishu;6.89\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6.1, result.MinTemperatureOf("Canberra"))
	assert.Equal(t, 5.97, result.MinTemperatureOf("Mogadishu"))
	assert.Equal(t, 10.23, result.MinTemperatureOf("Odesa"))
	assert.Equal(t, 9.79, result.MinTemperatureOf("Tirana"))
}

func TestParseWithTemperaturesForMaxTemperature(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.23\nMogadishu;5.97\nTirana;15.12\nTirana;9.79\nMogadishu;6.89\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6.1, result.MaxTemperatureOf("Canberra"))
	assert.Equal(t, 6.89, result.MaxTemperatureOf("Mogadishu"))
	assert.Equal(t, 10.23, result.MaxTemperatureOf("Odesa"))
	assert.Equal(t, 15.12, result.MaxTemperatureOf("Tirana"))
}

func TestParseWithTemperaturesForAverageTemperature(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.23\nMogadishu;5.97\nTirana;15.12\nTirana;9.79\nMogadishu;6.89\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6.1, result.AverageTemperatureOf("Canberra"))
	assert.Equal(t, 6.43, result.AverageTemperatureOf("Mogadishu"))
	assert.Equal(t, 10.23, result.AverageTemperatureOf("Odesa"))
	assert.InDelta(t, 12.45, result.AverageTemperatureOf("Tirana"), 0.01)
}
