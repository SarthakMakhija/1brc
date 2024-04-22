package brc

import (
	"1brc/bytes"
	"bufio"
	bytes2 "bytes"
	"github.com/dolthub/swiss"
	"io"
	"sort"
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

func (statistic StationTemperatureStatistics) stringify(stationName string, buffer *bytes2.Buffer) string {
	buffer.Reset()

	buffer.WriteString(stationName)
	buffer.WriteByte(':')
	buffer.WriteString(strconv.FormatFloat(statistic.minTemperature, 'f', -1, 64))
	buffer.WriteByte('/')
	buffer.WriteString(strconv.FormatFloat(statistic.averageTemperature, 'f', -1, 64))
	buffer.WriteByte('/')
	buffer.WriteString(strconv.FormatFloat(statistic.maxTemperature, 'f', -1, 64))

	return buffer.String()
}

const (
	maxSizeOfStationName            = 100
	numberOfSeparators              = 3
	maxSizeOfTemperature            = 4
	printableBufferSizePerStatistic = maxSizeOfStationName + numberOfSeparators + maxSizeOfTemperature*3
)

type StationTemperatureStatisticsResult struct {
	statisticsByStationName *swiss.Map[string, *StationTemperatureStatistics]
	printableBuffer         *bytes2.Buffer
}

func NewStationTemperatureStatisticsResult(statisticsByStationName *swiss.Map[string, *StationTemperatureStatistics]) StationTemperatureStatisticsResult {
	printableBuffer := &bytes2.Buffer{}
	printableBuffer.Grow(printableBufferSizePerStatistic)

	return StationTemperatureStatisticsResult{
		statisticsByStationName: statisticsByStationName,
		printableBuffer:         printableBuffer,
	}
}

func (result StationTemperatureStatisticsResult) get(stationName string) (*StationTemperatureStatistics, bool) {
	return result.statisticsByStationName.Get(stationName)
}

func (result StationTemperatureStatisticsResult) minTemperatureOf(stationName string) float64 {
	statistic, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return statistic.minTemperature
}

func (result StationTemperatureStatisticsResult) maxTemperatureOf(stationName string) float64 {
	statistic, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return statistic.maxTemperature
}

func (result StationTemperatureStatisticsResult) averageTemperatureOf(stationName string) float64 {
	statistic, ok := result.statisticsByStationName.Get(stationName)
	if !ok {
		return 0.0
	}
	return statistic.averageTemperature
}

func (result StationTemperatureStatisticsResult) allStationsSorted() []string {
	stationNames := make([]string, 0, result.statisticsByStationName.Count())
	result.statisticsByStationName.Iter(func(k string, _ *StationTemperatureStatistics) (stop bool) {
		stationNames = append(stationNames, k)
		return false
	})
	sort.Strings(stationNames)
	return stationNames
}

func (result StationTemperatureStatisticsResult) PrintableResult() string {
	output := strings.Builder{}
	output.WriteString("{")

	stationNames := result.allStationsSorted()
	for _, stationName := range stationNames {
		statistic, _ := result.get(stationName)
		output.WriteString(statistic.stringify(stationName, result.printableBuffer))
		output.WriteString(";")
	}
	output.WriteString("}")
	return output.String()
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
