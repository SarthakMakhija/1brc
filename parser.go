package brc

import (
	"bufio"
	"github.com/emirpasic/gods/maps/treemap"
	"io"
	"strconv"
	"strings"
)

type StationTemperatureStatistics struct {
	minTemperature       float64
	maxTemperature       float64
	aggregateTemperature float64
	totalEntries         uint64
	averageTemperature   float64
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
	return v.(StationTemperatureStatistics).minTemperature
}

func (result StationTemperatureStatisticsResult) MaxTemperatureOf(stationName string) float64 {
	v, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return v.(StationTemperatureStatistics).maxTemperature
}

func (result StationTemperatureStatisticsResult) AverageTemperatureOf(stationName string) float64 {
	v, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return v.(StationTemperatureStatistics).averageTemperature
}

func (result StationTemperatureStatisticsResult) AllStationsSorted() []interface{} {
	return result.statisticsByStationName.Keys()
}

func Parse(reader io.Reader) (StationTemperatureStatisticsResult, error) {
	scanner := bufio.NewScanner(reader)
	statisticsByStationName := treemap.NewWithStringComparator()

	for scanner.Scan() {
		line := scanner.Text()
		stationName, temperature, err := temperatureByStationName(line)
		if err != nil {
			if err == io.EOF {
				return NewStationTemperatureStatisticsResult(statisticsByStationName), nil
			}
			return StationTemperatureStatisticsResult{}, err
		}
		stats, ok := statisticsByStationName.Get(stationName)
		if !ok {
			statisticsByStationName.Put(stationName, StationTemperatureStatistics{
				minTemperature:       temperature,
				maxTemperature:       temperature,
				aggregateTemperature: temperature,
				totalEntries:         1,
				averageTemperature:   temperature,
			})
		} else {
			existingStatistics := stats.(StationTemperatureStatistics)
			minTemperature, maxTemperature := existingStatistics.minTemperature, existingStatistics.maxTemperature
			if temperature < existingStatistics.minTemperature {
				minTemperature = temperature
			}
			if temperature > existingStatistics.maxTemperature {
				maxTemperature = temperature
			}
			statisticsByStationName.Put(stationName, StationTemperatureStatistics{
				minTemperature:       minTemperature,
				maxTemperature:       maxTemperature,
				aggregateTemperature: temperature + existingStatistics.aggregateTemperature,
				totalEntries:         existingStatistics.totalEntries + 1,
				averageTemperature:   (temperature + existingStatistics.aggregateTemperature) / float64(existingStatistics.totalEntries+1),
			})
		}
	}
	return NewStationTemperatureStatisticsResult(statisticsByStationName), nil
}

func temperatureByStationName(line string) (string, float64, error) {
	parts := strings.Split(line, ";")
	temperature, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return "", 0, err
	}
	return parts[0], temperature, nil
}
