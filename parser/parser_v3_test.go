package parser

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseWithTemperaturesForSortedStationsNamesV3(t *testing.T) {
	result, err := ParseV3("../fixture/10_weather_stations.txt")

	assert.Nil(t, err)
	assert.Equal(t, []string{"Canberra", "Halifax", "Mogadishu", "Odesa", "Tirana"}, result.allStationsSorted())
}

func TestParseWithTemperaturesForMinTemperatureV3(t *testing.T) {
	result, err := ParseV3("../fixture/10_weather_stations.txt")

	assert.Nil(t, err)
	assert.Equal(t, float32(7.0), result.minTemperatureOf("Canberra"))
	assert.Equal(t, float32(6.0), result.minTemperatureOf("Mogadishu"))
	assert.Equal(t, float32(6.0), result.minTemperatureOf("Odesa"))
	assert.Equal(t, float32(3.0), result.minTemperatureOf("Tirana"))
	assert.Equal(t, float32(10.0), result.minTemperatureOf("Halifax"))
}

func TestParseWithTemperaturesForMaxTemperatureV3(t *testing.T) {
	result, err := ParseV3("../fixture/10_weather_stations.txt")

	assert.Nil(t, err)
	assert.Equal(t, float32(7.0), result.maxTemperatureOf("Canberra"))
	assert.Equal(t, float32(90.0), result.maxTemperatureOf("Mogadishu"))
	assert.Equal(t, float32(66.0), result.maxTemperatureOf("Odesa"))
	assert.Equal(t, float32(9.0), result.maxTemperatureOf("Tirana"))
	assert.Equal(t, float32(10.0), result.maxTemperatureOf("Halifax"))
}

func TestParseWithTemperaturesForAverageTemperatureV3(t *testing.T) {
	result, err := ParseV3("../fixture/10_weather_stations.txt")

	assert.Nil(t, err)
	assert.Equal(t, float32(7.0), result.averageTemperatureOf("Canberra"))
	assert.Equal(t, float32(48.0), result.averageTemperatureOf("Mogadishu"))
	assert.Equal(t, float32(27.0), result.averageTemperatureOf("Odesa"))
	assert.Equal(t, float32(6.0), result.averageTemperatureOf("Tirana"))
	assert.Equal(t, float32(10.0), result.averageTemperatureOf("Halifax"))
}
