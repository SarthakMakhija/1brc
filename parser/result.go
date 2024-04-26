package parser

import (
	"1brc/bytes"
	bytes2 "bytes"
	"github.com/dolthub/swiss"
	"sort"
)

type StationTemperatureStatisticsResult struct {
	statisticsByStationName    *swiss.Map[string, *StationTemperatureStatistics]
	printableTemperatureBuffer *bytes2.Buffer
	printableBuffer            *bytes2.Buffer
}

func NewStationTemperatureStatisticsResult(statisticsByStationName *swiss.Map[string, *StationTemperatureStatistics]) StationTemperatureStatisticsResult {
	printableBuffer := &bytes2.Buffer{}
	printableBuffer.Grow(printableBufferSizePerStatistic)

	printableTemperatureBuffer := &bytes2.Buffer{}
	printableTemperatureBuffer.Grow(64)
	return StationTemperatureStatisticsResult{
		statisticsByStationName:    statisticsByStationName,
		printableTemperatureBuffer: printableTemperatureBuffer,
		printableBuffer:            printableBuffer,
	}
}

func (result StationTemperatureStatisticsResult) get(stationName string) (*StationTemperatureStatistics, bool) {
	return result.statisticsByStationName.Get(stationName)
}

func (result StationTemperatureStatisticsResult) minTemperatureOf(stationName string) float32 {
	statistic, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return bytes.TemperatureAsFloat32(statistic.minTemperature)
}

func (result StationTemperatureStatisticsResult) maxTemperatureOf(stationName string) float32 {
	statistic, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return bytes.TemperatureAsFloat32(statistic.maxTemperature)
}

func (result StationTemperatureStatisticsResult) averageTemperatureOf(stationName string) float32 {
	statistic, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	averageTemperature := float32(statistic.aggregateTemperature) / float32(statistic.totalEntries)
	return averageTemperature * 0.1
}

func (result StationTemperatureStatisticsResult) allStationsSorted() []string {
	stationNames := make([]string, result.statisticsByStationName.Count())
	index := 0
	result.statisticsByStationName.Iter(func(stationName string, _ *StationTemperatureStatistics) (stop bool) {
		stationNames[index] = stationName
		index++
		return false
	})
	sort.Strings(stationNames)
	return stationNames
}

func (result StationTemperatureStatisticsResult) PrintableResult() string {
	stationNames := result.allStationsSorted()
	stationCount := len(stationNames)

	output := &bytes2.Buffer{}
	output.Grow(printableBufferSizePerStatistic*stationCount + 2 + stationCount)
	output.WriteByte('{')

	for _, stationName := range stationNames {
		statistic, _ := result.get(stationName)
		output.Write(statistic.stringify(stationName, result.printableTemperatureBuffer, result.printableBuffer))
		output.WriteByte(';')
	}
	output.WriteByte('}')
	return output.String()
}
