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
		for index := range statisticsChunkSummary.statisticsByStationName.entries {
			entry := statisticsChunkSummary.statisticsByStationName.entries[index]
			if entry.statistics.totalEntries > 0 {
				update(entry.station[:entry.stationLength], &entry.statistics, statisticsByStationName)
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
	buffer := currentBuffer[:]
	for len(buffer) > 0 {
		newLineIndex := bytes2.IndexByte(buffer, '\n')
		if newLineIndex == -1 {
			return buffer
		}
		separatorIndex := -1
		hash := uint64(fnv1aOffset64)
		for index, ch := range buffer[:newLineIndex] {
			if ch == bytes.Separator {
				separatorIndex = index
				break
			}
			hash ^= uint64(ch)
			hash *= fnv1aPrime64
		}

		temperature := bytes.ToTemperature(buffer[separatorIndex+1 : newLineIndex])
		statistics := statisticsByStationName.Get(hash, buffer[:separatorIndex])

		if statistics.totalEntries == 0 {
			statistics.minTemperature = temperature
			statistics.maxTemperature = temperature
			statistics.aggregateTemperature = int64(temperature)
			statistics.totalEntries = 1
		} else {
			if temperature < statistics.minTemperature {
				statistics.minTemperature = temperature
			} else if temperature > statistics.maxTemperature {
				statistics.maxTemperature = temperature
			}
			statistics.aggregateTemperature = statistics.aggregateTemperature + int64(temperature)
			statistics.totalEntries += 1
			statisticsByStationName.entryCount += 1
		}
		if newLineIndex+1 < len(buffer) {
			buffer = buffer[newLineIndex+1:]
		} else {
			return nil
		}
	}
	return buffer
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
	station       [128]byte
	statistics    StationTemperatureStatistics
	hash          uint64
	stationLength int
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

func (statisticsByStationName *StatisticsByStationNameMap) Get(hash uint64, stationName []byte) *StationTemperatureStatistics {
	index := hash & statisticsByStationName.mask
	entry := &statisticsByStationName.entries[index]

	for entry.stationLength > 0 && !(entry.hash == hash && bytes2.Equal(entry.station[:entry.stationLength], stationName)) {
		index = (index + 1) & statisticsByStationName.mask
		entry = &statisticsByStationName.entries[index]
	}
	if entry.stationLength == 0 {
		entry.hash = hash
		entry.stationLength = copy(entry.station[:], stationName)
	}
	return &entry.statistics
}
