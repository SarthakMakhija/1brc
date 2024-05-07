package parser

import (
	brc "1brc"
	"1brc/bytes"
	"io"
)

func ParseV2(reader io.Reader) (StationTemperatureStatisticsSummary, error) {
	statisticsByStationName := make(map[string]*StationTemperatureStatistics, 10_0000)
	buffer := make([]byte, brc.ReadSize)

	var err error
	var n int
	var offset int

	for err != io.EOF {
		n, err = reader.Read(buffer[offset:])
		if n > 0 {
			n = n + offset
			var last int
			for index := range buffer[:n] {
				if buffer[index] == '\n' {
					stationName, temperature, err := bytes.SplitIntoStationNameAndTemperature(buffer[last:index])
					if err != nil {
						return StationTemperatureStatisticsSummary{}, err
					}
					updateStatistics(stationName, temperature, statisticsByStationName)
					last = index + 1
				}
			}
			offset = n - last
			if offset > 0 {
				copy(buffer, buffer[last:n])
			}
		}
		if err != nil {
			if err != io.EOF {
				return StationTemperatureStatisticsSummary{}, err
			}
		}
	}
	return NewStationTemperatureStatisticsResult(statisticsByStationName), nil
}

func updateStatistics(stationName []byte, temperature bytes.Temperature, statisticsByStationName map[string]*StationTemperatureStatistics) {
	existingStatistics, ok := statisticsByStationName[string(stationName)]
	if !ok {
		statisticsByStationName[string(stationName)] = &StationTemperatureStatistics{
			minTemperature:       temperature,
			maxTemperature:       temperature,
			aggregateTemperature: temperature,
			totalEntries:         1,
		}
	} else {
		if temperature < existingStatistics.minTemperature {
			existingStatistics.minTemperature = temperature
		} else if temperature > existingStatistics.maxTemperature {
			existingStatistics.maxTemperature = temperature
		}
		existingStatistics.aggregateTemperature = temperature + existingStatistics.aggregateTemperature
		existingStatistics.totalEntries = existingStatistics.totalEntries + 1
	}
}
