package brc

import (
	"1brc/bytes"
	"bufio"
	"fmt"
	"github.com/dolthub/swiss"
	"io"
	"sort"
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
	statisticsByStationName *swiss.Map[string, *StationTemperatureStatistics]
}

func NewStationTemperatureStatisticsResult(statisticsByStationName *swiss.Map[string, *StationTemperatureStatistics]) StationTemperatureStatisticsResult {
	return StationTemperatureStatisticsResult{
		statisticsByStationName: statisticsByStationName,
	}
}

func (result StationTemperatureStatisticsResult) Get(stationName string) (*StationTemperatureStatistics, bool) {
	return result.statisticsByStationName.Get(stationName)
}

func (result StationTemperatureStatisticsResult) MinTemperatureOf(stationName string) float64 {
	statistic, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return statistic.minTemperature
}

func (result StationTemperatureStatisticsResult) MaxTemperatureOf(stationName string) float64 {
	statistic, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return statistic.maxTemperature
}

func (result StationTemperatureStatisticsResult) AverageTemperatureOf(stationName string) float64 {
	statistic, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return statistic.averageTemperature
}

func (result StationTemperatureStatisticsResult) AllStationsSorted() []string {
	stationNames := make([]string, 0, result.statisticsByStationName.Count())
	result.statisticsByStationName.Iter(func(k string, _ *StationTemperatureStatistics) (stop bool) {
		stationNames = append(stationNames, k)
		return false
	})
	sort.Strings(stationNames)
	return stationNames
}

// Parse
// TODO: rounding
func Parse(reader io.Reader) (StationTemperatureStatisticsResult, error) {
	scanner := bufio.NewScanner(reader)
	statisticsByStationName := swiss.NewMap[string, *StationTemperatureStatistics](10_000)

	for scanner.Scan() {
		line := scanner.Bytes()
		stationName, temperature, err := temperatureByStationName(line)
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
			existingStatistics.averageTemperature = (existingStatistics.aggregateTemperature) / float64(existingStatistics.totalEntries)
		}
	}
	return NewStationTemperatureStatisticsResult(statisticsByStationName), nil
}

func temperatureByStationName(line []byte) ([]byte, float64, error) {
	stationName, temperature, err := bytes.SplitIntoStationNameAndTemperature(line)
	if err != nil {
		return nil, 0, err
	}
	return stationName, temperature, nil
}
