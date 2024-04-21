package brc

import (
	"1brc/bytes"
	"bufio"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"io"
	"strconv"
)

type StationTemperatureStatistics struct {
	minTemperature       float64
	maxTemperature       float64
	aggregateTemperature float64
	totalEntries         uint64
	averageTemperature   float64
}

func (statistic StationTemperatureStatistics) Stringify(stationName string) string {
	return fmt.Sprintf("%s:%v/%v/%v", stationName, statistic.minTemperature, statistic.averageTemperature, statistic.maxTemperature)
}

type StationTemperatureStatisticsResult struct {
	statisticsByStationName *treemap.Map
}

func NewStationTemperatureStatisticsResult(statisticsByStationName *treemap.Map) StationTemperatureStatisticsResult {
	return StationTemperatureStatisticsResult{
		statisticsByStationName: statisticsByStationName,
	}
}

func (result StationTemperatureStatisticsResult) MinTemperatureOf(stationName string) float64 {
	v, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return v.(*StationTemperatureStatistics).minTemperature
}

func (result StationTemperatureStatisticsResult) MaxTemperatureOf(stationName string) float64 {
	v, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return v.(*StationTemperatureStatistics).maxTemperature
}

func (result StationTemperatureStatisticsResult) AverageTemperatureOf(stationName string) float64 {
	v, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return v.(*StationTemperatureStatistics).averageTemperature
}

func (result StationTemperatureStatisticsResult) AllStationsSorted() []interface{} {
	return result.statisticsByStationName.Keys()
}

func (result StationTemperatureStatisticsResult) Iterator() treemap.Iterator {
	return result.statisticsByStationName.Iterator()
}

// Parse
// TODO: rounding
func Parse(reader io.Reader) (StationTemperatureStatisticsResult, error) {
	scanner := bufio.NewScanner(reader)
	statisticsByStationName := treemap.NewWithStringComparator()

	for scanner.Scan() {
		line := scanner.Bytes()
		stationName, temperature, err := temperatureByStationName(line)
		if err != nil {
			if err == io.EOF {
				return NewStationTemperatureStatisticsResult(statisticsByStationName), nil
			}
			return StationTemperatureStatisticsResult{}, err
		}
		stats, ok := statisticsByStationName.Get(stationName)
		if !ok {
			statisticsByStationName.Put(stationName, &StationTemperatureStatistics{
				minTemperature:       temperature,
				maxTemperature:       temperature,
				aggregateTemperature: temperature,
				totalEntries:         1,
				averageTemperature:   temperature,
			})
		} else {
			existingStatistics := stats.(*StationTemperatureStatistics)
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
	return NewStationTemperatureStatisticsResult(statisticsByStationName), nil
}

func temperatureByStationName(line []byte) (string, float64, error) {
	stationName, rawTemperature, err := bytes.SplitIntoStationNameAndTemperature(line)
	if err != nil {
		return "", 0, err
	}
	temperature, err := strconv.ParseFloat(rawTemperature, 64)
	if err != nil {
		return "", 0, err
	}
	return stationName, temperature, nil
}
