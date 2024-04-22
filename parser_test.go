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
	assert.Equal(t, []string{"Canberra", "Mogadishu", "Odesa", "Tirana"}, result.allStationsSorted())
}

func TestParseWithTemperaturesForMinTemperature(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.7\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6.1, result.minTemperatureOf("Canberra"))
	assert.Equal(t, 5.9, result.minTemperatureOf("Mogadishu"))
	assert.Equal(t, 10.2, result.minTemperatureOf("Odesa"))
	assert.Equal(t, 9.7, result.minTemperatureOf("Tirana"))
}

func TestParseWithTemperaturesForMaxTemperature(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.7\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6.1, result.maxTemperatureOf("Canberra"))
	assert.Equal(t, 6.8, result.maxTemperatureOf("Mogadishu"))
	assert.Equal(t, 10.2, result.maxTemperatureOf("Odesa"))
	assert.Equal(t, 15.1, result.maxTemperatureOf("Tirana"))
}

func TestParseWithTemperaturesForAverageTemperature(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.7\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	assert.Equal(t, 6.1, result.averageTemperatureOf("Canberra"))
	assert.Equal(t, 6.35, result.averageTemperatureOf("Mogadishu"))
	assert.Equal(t, 10.2, result.averageTemperatureOf("Odesa"))
	assert.InDelta(t, 12.4, result.averageTemperatureOf("Tirana"), 0.01)
}

func TestPrintableResult(t *testing.T) {
	input := bytes.NewReader([]byte("Odesa;10.2\nMogadishu;5.9\nTirana;15.1\nTirana;9.3\nMogadishu;6.8\nCanberra;6.1\n"))
	result, err := Parse(input)

	assert.Nil(t, err)
	printableResult := result.PrintableResult()

	assert.Equal(t,
		"{Canberra:6.1/6.1/6.1;Mogadishu:5.9/6.35/6.8;Odesa:10.2/10.2/10.2;Tirana:9.3/12.2/15.1;}",
		printableResult,
	)
}
