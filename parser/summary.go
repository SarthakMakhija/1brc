package parser

import (
	"1brc/bytes"
	bytes2 "bytes"
	"sort"
)

type StationTemperatureStatisticsSummary struct {
	statisticsByStationName    map[string]*StationTemperatureStatistics
	printableTemperatureBuffer *bytes2.Buffer
	printableBuffer            *bytes2.Buffer
}

func NewStationTemperatureStatisticsResult(statisticsByStationName map[string]*StationTemperatureStatistics) StationTemperatureStatisticsSummary {
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
	statistics, ok := summary.statisticsByStationName[(stationName)]
	return statistics, ok
}

func (summary StationTemperatureStatisticsSummary) minTemperatureOf(stationName string) float32 {
	statistic, ok := summary.get(stationName)
	if !ok {
		return 0.0
	}
	return bytes.TemperatureAsFloat32(statistic.minTemperature)
}

func (summary StationTemperatureStatisticsSummary) maxTemperatureOf(stationName string) float32 {
	statistic, ok := summary.get(stationName)
	if !ok {
		return 0.0
	}
	return bytes.TemperatureAsFloat32(statistic.maxTemperature)
}

func (summary StationTemperatureStatisticsSummary) averageTemperatureOf(stationName string) float32 {
	statistic, ok := summary.get(stationName)
	if !ok {
		return 0.0
	}
	averageTemperature := float32(statistic.aggregateTemperature) / float32(statistic.totalEntries)
	return averageTemperature * 0.1
}

func (summary StationTemperatureStatisticsSummary) allStationsSorted() []string {
	stationNames := make([]string, len(summary.statisticsByStationName))
	index := 0
	for stationName, _ := range summary.statisticsByStationName {
		stationNames[index] = stationName
		index++
	}
	sort.Strings(stationNames)
	return stationNames
}

const (
	unrollFactor = 4
	mask         = unrollFactor - 1
)

func (summary StationTemperatureStatisticsSummary) PrintableResult() string {
	stationNames := summary.allStationsSorted()
	stationCount := len(stationNames)

	output := &bytes2.Buffer{}
	output.Grow(printableBufferSizePerStatistic*stationCount + 2 + stationCount)
	output.WriteByte('{')

	index := 0
	unrolledIterations, pendingIterations := stationCount/unrollFactor, stationCount&mask
	for iteration := 1; iteration <= unrolledIterations; iteration++ {
		stationNamesLocal := stationNames[index : index+unrollFactor : index+unrollFactor]

		stationName := stationNamesLocal[0]
		statistic, _ := summary.get(stationName)
		output.Write(statistic.stringify(stationName, summary.printableTemperatureBuffer, summary.printableBuffer))
		output.WriteByte(';')

		stationName = stationNamesLocal[1]
		statistic, _ = summary.get(stationName)
		output.Write(statistic.stringify(stationName, summary.printableTemperatureBuffer, summary.printableBuffer))
		output.WriteByte(';')

		stationName = stationNamesLocal[2]
		statistic, _ = summary.get(stationName)
		output.Write(statistic.stringify(stationName, summary.printableTemperatureBuffer, summary.printableBuffer))
		output.WriteByte(';')

		stationName = stationNamesLocal[3]
		statistic, _ = summary.get(stationName)
		output.Write(statistic.stringify(stationName, summary.printableTemperatureBuffer, summary.printableBuffer))
		output.WriteByte(';')

		index += unrollFactor
	}

	for iteration := 1; iteration <= pendingIterations; iteration++ {
		stationName := stationNames[index]
		statistic, _ := summary.get(stationName)
		output.Write(statistic.stringify(stationName, summary.printableTemperatureBuffer, summary.printableBuffer))
		output.WriteByte(';')
		index += 1
	}
	output.WriteByte('}')
	return output.String()
}
