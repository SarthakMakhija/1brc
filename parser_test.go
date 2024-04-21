package brc

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseWithTemperaturesForSortedStationsNames(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.7\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, []string{"Canberra", "Mogadishu", "Odesa", "Tirana"}, result.AllStationsSorted())
}

func TestParseWithTemperaturesForMinTemperature(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.7\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6.1, result.MinTemperatureOf("Canberra"))
	assert.Equal(t, 5.9, result.MinTemperatureOf("Mogadishu"))
	assert.Equal(t, 10.2, result.MinTemperatureOf("Odesa"))
	assert.Equal(t, 9.7, result.MinTemperatureOf("Tirana"))
}

func TestParseWithTemperaturesForMaxTemperature(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.7\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6.1, result.MaxTemperatureOf("Canberra"))
	assert.Equal(t, 6.8, result.MaxTemperatureOf("Mogadishu"))
	assert.Equal(t, 10.2, result.MaxTemperatureOf("Odesa"))
	assert.Equal(t, 15.1, result.MaxTemperatureOf("Tirana"))
}

func TestParseWithTemperaturesForAverageTemperature(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.7\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6.1, result.AverageTemperatureOf("Canberra"))
	assert.Equal(t, 6.35, result.AverageTemperatureOf("Mogadishu"))
	assert.Equal(t, 10.2, result.AverageTemperatureOf("Odesa"))
	assert.InDelta(t, 12.4, result.AverageTemperatureOf("Tirana"), 0.01)
}
