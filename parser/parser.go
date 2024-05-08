package parser

import (
	"1brc/bytes"
	"bufio"
	bytes2 "bytes"
	"io"
)

type StationTemperatureStatistics struct {
	minTemperature       bytes.Temperature
	maxTemperature       bytes.Temperature
	aggregateTemperature int64 //TODO: decide if we need int64 for this.
	totalEntries         int64
}

func (statistic StationTemperatureStatistics) stringify(
	stationName string,
	temperatureBuffer *bytes2.Buffer,
	resultBuffer *bytes2.Buffer,
) []byte {
	temperatureBuffer.Reset()
	resultBuffer.Reset()

	resultBuffer.WriteString(stationName)
	resultBuffer.WriteByte(':')
	resultBuffer.Write(bytes.Format(
		statistic.minTemperature,
		temperatureBuffer.Bytes(),
	))
	resultBuffer.WriteByte('/')
	averageTemperature := float32(statistic.aggregateTemperature) / float32(statistic.totalEntries)
	resultBuffer.Write(bytes.FormatTemperatureAsFloat32(
		averageTemperature*0.1,
		temperatureBuffer.Bytes(),
	))
	resultBuffer.WriteByte('/')
	resultBuffer.Write(bytes.Format(
		statistic.maxTemperature,
		temperatureBuffer.Bytes(),
	))

	return resultBuffer.Bytes()
}

const (
	maxSizeOfStationName            = 100
	numberOfSeparators              = 3
	maxSizeOfTemperature            = 4
	printableBufferSizePerStatistic = maxSizeOfStationName + numberOfSeparators + maxSizeOfTemperature*3
)

// Parse
// TODO: rounding
func Parse(reader io.Reader) (StationTemperatureStatisticsSummary, error) {
	scanner := bufio.NewScanner(reader)
	statisticsByStationName := make(map[string]*StationTemperatureStatistics, 10_0000)

	for scanner.Scan() {
		line := scanner.Bytes()
		stationName, temperature, err := bytes.SplitIntoStationNameAndTemperature(line)
		if err != nil {
			if err == io.EOF {
				return NewStationTemperatureStatisticsSummary(statisticsByStationName), nil
			}
			return StationTemperatureStatisticsSummary{}, err
		}
		existingStatistics, ok := statisticsByStationName[(string(stationName))]
		if !ok {
			statisticsByStationName[string(stationName)] = &StationTemperatureStatistics{
				minTemperature:       temperature,
				maxTemperature:       temperature,
				aggregateTemperature: int64(temperature),
				totalEntries:         1,
			}
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
			existingStatistics.aggregateTemperature = int64(temperature) + existingStatistics.aggregateTemperature
			existingStatistics.totalEntries = existingStatistics.totalEntries + 1
		}
	}
	return NewStationTemperatureStatisticsSummary(statisticsByStationName), nil
}
