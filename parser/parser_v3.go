package parser

import (
	brc "1brc"
	"1brc/bytes"
	bytes2 "bytes"
	"io"
	"os"
	"runtime"
)

const (
	capacity      = 1 << 14
	fnv1aOffset64 = 14695981039346656037
	fnv1aPrime64  = 1099511628211
	entryMask     = capacity - 1
)

type Entry struct {
	station    []byte
	statistics *StationTemperatureStatistics
	hash       uint64
}

func ParseV3(filePath string) (StationTemperatureStatisticsSummary, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return StationTemperatureStatisticsSummary{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	numberOfParts := runtime.NumCPU()
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
		entryCount := statisticsChunkSummary.totalEntries

		for index := range statisticsChunkSummary.stationEntries {
			entry := statisticsChunkSummary.stationEntries[index]
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
	buffer := make([]byte, brc.ReadSize)

	stationEntries := make([]Entry, capacity)
	totalEntries := 0
	offset := 0
	for {
		n, err := reader.Read(buffer[offset:])
		if n+offset > 0 {
			leftOver, entryCount := updateStatisticsIn(stationEntries, buffer[:n+offset])
			totalEntries += entryCount

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
	return NewStationTemperatureStatisticsChunkSummary(stationEntries, totalEntries), nil
}

func updateStatisticsIn(entries []Entry, currentBuffer []byte) ([]byte, int) {
	lastNewLineIndex := bytes2.LastIndexByte(currentBuffer, '\n')
	if lastNewLineIndex == -1 {
		return nil, 0
	}

	const advanceCursorPostTemperatureBy = 2 //one for newline and one to move to the next line

	leftOver, buffer := currentBuffer[lastNewLineIndex+1:], currentBuffer[:lastNewLineIndex+1]
	entryCount := 0
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

		isMinus := false
		temperatureSlice := buffer[separatorIndex+1:]
		if temperatureSlice[0] == bytes.MinusSign {
			temperatureSlice = temperatureSlice[1:]
			isMinus = true
		}

		temperature, numberOfBytesRead := bytes.ToTemperatureWithNewLine(temperatureSlice, isMinus)
		index := int(hash & entryMask)
		statistics := getStatistics(entries, hash, index, buffer[:separatorIndex])

		if statistics.totalEntries == 0 {
			statistics.minTemperature = temperature
			statistics.maxTemperature = temperature
			statistics.aggregateTemperature = int64(temperature)
			statistics.totalEntries = 1
			entryCount++
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
	return leftOver, entryCount
}

// getStatistics inlinable (cost 74)
func getStatistics(entries []Entry, hash uint64, index int, stationName []byte) *StationTemperatureStatistics {
	for {
		if entries[index].station == nil {
			key := make([]byte, len(stationName))
			copy(key, stationName)
			entry := Entry{
				hash:       hash,
				station:    key,
				statistics: &StationTemperatureStatistics{},
			}
			entries[index] = entry
			return entry.statistics
		}

		if entries[index].hash == hash && bytes2.Equal(entries[index].station, stationName) {
			return entries[index].statistics
		}
		index++
		if index >= capacity {
			index = 0
		}
	}
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
