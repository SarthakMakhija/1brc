package parser

import (
	brc "1brc"
	"1brc/bytes"
	"github.com/dolthub/swiss"
	"io"
)

func ParseV2(reader io.Reader) (StationTemperatureStatisticsResult, error) {
	statisticsByStationName := swiss.NewMap[string, *StationTemperatureStatistics](10_000)
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
						return StationTemperatureStatisticsResult{}, err
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
				return StationTemperatureStatisticsResult{}, err
			}
		}
	}
	return NewStationTemperatureStatisticsResult(statisticsByStationName), nil
}

func updateStatistics(stationName []byte, temperature float64, statisticsByStationName *swiss.Map[string, *StationTemperatureStatistics]) {
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
		existingStatistics.averageTemperature = (existingStatistics.aggregateTemperature) / float64(existingStatistics.totalEntries)
	}
}
