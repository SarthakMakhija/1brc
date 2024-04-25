package parser

import (
	"1brc/bytes"
	"bufio"
	bytes2 "bytes"
	"github.com/dolthub/swiss"
	"io"
)

type StationTemperatureStatistics struct {
	minTemperature       float64
	maxTemperature       float64
	aggregateTemperature float64
	totalEntries         uint64
	averageTemperature   float64
}

func (statistic StationTemperatureStatistics) stringify(
	stationName string,
	temperatureBuffer *bytes2.Buffer,
	resultBuffer *bytes2.Buffer,
) string {
	temperatureBuffer.Reset()
	resultBuffer.Reset()

	resultBuffer.WriteString(stationName)
	resultBuffer.WriteByte(':')
	resultBuffer.WriteString(bytes.Format(
		statistic.minTemperature,
		temperatureBuffer.Bytes(),
	))
	resultBuffer.WriteByte('/')
	resultBuffer.WriteString(bytes.Format(
		statistic.aggregateTemperature/float64(statistic.totalEntries),
		temperatureBuffer.Bytes(),
	))
	resultBuffer.WriteByte('/')
	resultBuffer.WriteString(bytes.Format(
		statistic.maxTemperature,
		temperatureBuffer.Bytes(),
	))

	return resultBuffer.String()
}

const (
	maxSizeOfStationName            = 100
	numberOfSeparators              = 3
	maxSizeOfTemperature            = 4
	printableBufferSizePerStatistic = maxSizeOfStationName + numberOfSeparators + maxSizeOfTemperature*3
)

// Parse
// TODO: rounding
func Parse(reader io.Reader) (StationTemperatureStatisticsResult, error) {
	scanner := bufio.NewScanner(reader)
	statisticsByStationName := swiss.NewMap[string, *StationTemperatureStatistics](10_000)

	for scanner.Scan() {
		line := scanner.Bytes()
		stationName, temperature, err := bytes.SplitIntoStationNameAndTemperature(line)
		if err != nil {
			if err == io.EOF {
				return NewStationTemperatureStatisticsResult(statisticsByStationName), nil
			}
			return StationTemperatureStatisticsResult{}, err
		}
		existingStatistics, ok := statisticsByStationName.Get(string(stationName))
		if !ok {
			statisticsByStationName.Put(string(stationName), &StationTemperatureStatistics{
				minTemperature:       temperature,
				maxTemperature:       temperature,
				aggregateTemperature: temperature,
				totalEntries:         1,
				averageTemperature:   temperature,
			})
		} else {
			minTemperature, maxTemperature := existingStatistics.minTemperature, existingStatistics.maxTemperature
			if temperature < existingStatistics.minTemperature {
				minTemperature = temperature
			}
			if temperature > existingStatistics.maxTemperature {
				maxTemperature = temperature
			}
			existingStatistics.minTemperature = minTemperature
			existingStatistics.maxTemperature = maxTemperature
			existingStatistics.aggregateTemperature = temperature + existingStatistics.aggregateTemperature
			existingStatistics.totalEntries = existingStatistics.totalEntries + 1
		}
	}
	return NewStationTemperatureStatisticsResult(statisticsByStationName), nil
}
