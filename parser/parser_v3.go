package parser

import (
	brc "1brc"
	"1brc/bytes"
	bytes2 "bytes"
	"io"
	"os"
	"runtime"
)

/**
1. See if NewStatisticsByStationNameMap can be converted to a method
2. Change capacity (currently 1<<14)
*/

const (
	capacity      = 1 << 14
	fnv1aOffset64 = 14695981039346656037
	fnv1aPrime64  = 1099511628211
)

func ParseV3(filePath string) (StationTemperatureStatisticsSummary, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return StationTemperatureStatisticsSummary{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	numberOfParts := runtime.NumCPU() //TODO: adjust this?
	chunks, err := bytes.SplitFile(filePath, numberOfParts)
	if err != nil {
		return StationTemperatureStatisticsSummary{}, err
	}

	chunkSummaries := make(chan StationTemperatureStatisticsChunkSummary, len(chunks))
	for _, chunk := range chunks {
		go func(chunk bytes.Chunk) {
			statisticsChunkSummary, err := readChunk(file, chunk)
			if err != nil {
				panic(err)
			}
			chunkSummaries <- statisticsChunkSummary
		}(chunk)
	}

	statisticsByStationName := make(map[string]*StationTemperatureStatistics, 1<<16) //TODO: allocate size?
	for i := 0; i < len(chunks); i++ {
		statisticsChunkSummary := <-chunkSummaries
		entryCount := statisticsChunkSummary.statisticsByStationName.entryCount

		for index := range statisticsChunkSummary.statisticsByStationName.entries {
			entry := statisticsChunkSummary.statisticsByStationName.entries[index]
			if entry.station != nil {
				update(entry.station, entry.statistics, statisticsByStationName)
				entryCount--
			}
			if entryCount == 0 {
				break
			}
		}
	}
	return NewStationTemperatureStatisticsSummary(statisticsByStationName), nil
}

func readChunk(file *os.File, chunk bytes.Chunk) (StationTemperatureStatisticsChunkSummary, error) {
	reader := io.NewSectionReader(file, chunk.StartOffset, chunk.Size)
	statisticsByStationName := NewStatisticsByStationNameMap(capacity)
	buffer := make([]byte, brc.ReadSize)

	offset := 0
	for {
		n, err := reader.Read(buffer[offset:])
		if n+offset > 0 {
			leftOver := updateStatisticsIn(buffer[:n+offset], statisticsByStationName)
			if len(leftOver) > 0 {
				offset = copy(buffer[:], leftOver)
			} else {
				offset = 0
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			if err != io.EOF {
				return StationTemperatureStatisticsChunkSummary{}, err
			}
		}
	}
	return NewStationTemperatureStatisticsChunkSummary(statisticsByStationName), nil
}

func updateStatisticsIn(currentBuffer []byte, statisticsByStationName *StatisticsByStationNameMap) []byte {
	lastNewLineIndex := bytes2.LastIndexByte(currentBuffer, '\n')
	if lastNewLineIndex == -1 {
		return nil
	}

	const advanceCursorPostTemperatureBy = 2 //one for newline and one to move to the next line

	leftOver, buffer := currentBuffer[lastNewLineIndex+1:], currentBuffer[:lastNewLineIndex+1]
	for len(buffer) > 0 {
		separatorIndex := -1
		hash := uint64(fnv1aOffset64)

		for index, ch := range buffer[:] {
			if ch == bytes.Separator {
				separatorIndex = index
				break
			}
			hash ^= uint64(ch)
			hash *= fnv1aPrime64
		}
		temperature, numberOfBytesRead := bytes.ToTemperatureWithNewLine(buffer[separatorIndex+1:])
		statistics := statisticsByStationName.GetOrEmptyStatisticsFor(hash, buffer[:separatorIndex])

		if statistics.totalEntries == 0 {
			statistics.minTemperature = temperature
			statistics.maxTemperature = temperature
			statistics.aggregateTemperature = int64(temperature)
			statistics.totalEntries = 1
			statisticsByStationName.entryCount += 1
		} else {
			if temperature < statistics.minTemperature {
				statistics.minTemperature = temperature
			} else if temperature > statistics.maxTemperature {
				statistics.maxTemperature = temperature
			}
			statistics.aggregateTemperature = statistics.aggregateTemperature + int64(temperature)
			statistics.totalEntries += 1
		}
		buffer = buffer[separatorIndex+numberOfBytesRead+advanceCursorPostTemperatureBy:]
	}
	return leftOver
}

func update(stationName []byte, summary *StationTemperatureStatistics, statisticsByStationName map[string]*StationTemperatureStatistics) {
	existingStatistics, ok := statisticsByStationName[string(stationName)]
	if !ok {
		statisticsByStationName[string(stationName)] = summary
	} else {
		if summary.minTemperature < existingStatistics.minTemperature {
			existingStatistics.minTemperature = summary.minTemperature
		}
		if summary.maxTemperature > existingStatistics.maxTemperature {
			existingStatistics.maxTemperature = summary.maxTemperature
		}
		existingStatistics.aggregateTemperature = summary.aggregateTemperature + existingStatistics.aggregateTemperature
		existingStatistics.totalEntries = summary.totalEntries + existingStatistics.totalEntries
	}
}

type Entry struct {
	station    []byte
	statistics *StationTemperatureStatistics
	hash       uint64
}

type StatisticsByStationNameMap struct {
	entries    []Entry
	mask       uint64
	capacity   int
	entryCount int
}

func NewStatisticsByStationNameMap(capacity int) *StatisticsByStationNameMap {
	return &StatisticsByStationNameMap{
		entries:    make([]Entry, capacity),
		mask:       uint64(capacity - 1),
		capacity:   capacity,
		entryCount: 0,
	}
}

func (statisticsByStationName *StatisticsByStationNameMap) GetOrEmptyStatisticsFor(hash uint64, stationName []byte) *StationTemperatureStatistics {
	index := hash & statisticsByStationName.mask
	for {
		entry := &statisticsByStationName.entries[index]
		if entry.station == nil {
			key := make([]byte, len(stationName))
			copy(key, stationName)
			*entry = Entry{
				hash:       hash,
				station:    key,
				statistics: &StationTemperatureStatistics{},
			}
			return entry.statistics
		}

		if entry.hash == hash && bytes2.Equal(entry.station, stationName) {
			return entry.statistics
		}
		index = (index + 1) & statisticsByStationName.mask
	}
}
