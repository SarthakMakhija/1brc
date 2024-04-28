package parser

import (
	"1brc/bytes"
	bytes2 "bytes"
	"github.com/dolthub/swiss"
	"sort"
)

type StationTemperatureStatisticsSummary struct {
	statisticsByStationName    *swiss.Map[string, *StationTemperatureStatistics]
	printableTemperatureBuffer *bytes2.Buffer
	printableBuffer            *bytes2.Buffer
}

func NewStationTemperatureStatisticsResult(statisticsByStationName *swiss.Map[string, *StationTemperatureStatistics]) StationTemperatureStatisticsSummary {
	printableBuffer := &bytes2.Buffer{}
	printableBuffer.Grow(printableBufferSizePerStatistic)

	printableTemperatureBuffer := &bytes2.Buffer{}
	printableTemperatureBuffer.Grow(32)
	return StationTemperatureStatisticsSummary{
		statisticsByStationName:    statisticsByStationName,
		printableTemperatureBuffer: printableTemperatureBuffer,
		printableBuffer:            printableBuffer,
	}
}

func (summary StationTemperatureStatisticsSummary) get(stationName string) (*StationTemperatureStatistics, bool) {
	return summary.statisticsByStationName.Get(stationName)
}

func (summary StationTemperatureStatisticsSummary) minTemperatureOf(stationName string) float32 {
	statistic, ok := summary.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return bytes.TemperatureAsFloat32(statistic.minTemperature)
}

func (summary StationTemperatureStatisticsSummary) maxTemperatureOf(stationName string) float32 {
	statistic, ok := summary.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return bytes.TemperatureAsFloat32(statistic.maxTemperature)
}

func (summary StationTemperatureStatisticsSummary) averageTemperatureOf(stationName string) float32 {
	statistic, ok := summary.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	averageTemperature := float32(statistic.aggregateTemperature) / float32(statistic.totalEntries)
	return averageTemperature * 0.1
}

func (summary StationTemperatureStatisticsSummary) allStationsSorted() []string {
	stationNames := make([]string, summary.statisticsByStationName.Count())
	index := 0
	summary.statisticsByStationName.Iter(func(stationName string, _ *StationTemperatureStatistics) (stop bool) {
		stationNames[index] = stationName
		index++
		return false
	})
	sort.Strings(stationNames)
	return stationNames
}

func (summary StationTemperatureStatisticsSummary) PrintableResult() string {
	stationNames := summary.allStationsSorted()
	stationCount := len(stationNames)

	output := &bytes2.Buffer{}
	output.Grow(printableBufferSizePerStatistic*stationCount + 2 + stationCount)
	output.WriteByte('{')

	for _, stationName := range stationNames {
		statistic, _ := summary.get(stationName)
		output.Write(statistic.stringify(stationName, summary.printableTemperatureBuffer, summary.printableBuffer))
		output.WriteByte(';')
	}
	output.WriteByte('}')
	return output.String()
}
