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

func TestGetInStatisticsByStationNameMap(t *testing.T) {
	entry := Entry{}
	entry.statistics = &StationTemperatureStatistics{minTemperature: -990, maxTemperature: 10}
	entry.hash = 11
	entry.station = []byte("Odesa")

	statisticsByStationNameMap := NewStatisticsByStationNameMap(8)
	statisticsByStationNameMap.entries[3] = entry

	statistics := statisticsByStationNameMap.GetOrEmptyStatisticsFor(11, []byte("Odesa"))
	assert.Equal(t, int16(-990), statistics.minTemperature)
	assert.Equal(t, int16(10), statistics.maxTemperature)
}

func TestGetWithHashConflictInStatisticsByStationNameMap(t *testing.T) {
	entry := Entry{}
	entry.statistics = &StationTemperatureStatistics{minTemperature: -990, maxTemperature: 10}
	entry.hash = 11
	entry.station = []byte("Odesa")

	anotherEntry := Entry{}
	anotherEntry.statistics = &StationTemperatureStatistics{minTemperature: -890, maxTemperature: 20}
	anotherEntry.hash = 11
	anotherEntry.station = []byte("Delhi")

	statisticsByStationNameMap := NewStatisticsByStationNameMap(8)
	statisticsByStationNameMap.entries[3] = entry
	statisticsByStationNameMap.entries[4] = anotherEntry

	statistics := statisticsByStationNameMap.GetOrEmptyStatisticsFor(11, []byte("Delhi"))
	assert.Equal(t, int16(-890), statistics.minTemperature)
	assert.Equal(t, int16(20), statistics.maxTemperature)
}

func TestGetInStatisticsByStationNameMapWithEmptyStatistics(t *testing.T) {
	statisticsByStationNameMap := NewStatisticsByStationNameMap(8)

	statistics := statisticsByStationNameMap.GetOrEmptyStatisticsFor(11, []byte("Odesa"))
	assert.Equal(t, int64(0), statistics.totalEntries)
}
