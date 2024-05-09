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
	capacity      = 1 << 16
	fnv1aOffset64 = 14695981039346656037
	fnv1aPrime64  = 1099511628211
)

func ParseV3(filePath string) (StationTemperatureStatisticsSummary, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return StationTemperatureStatisticsSummary{}, err
	}

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

	statisticsByStationName := make(map[string]*StationTemperatureStatistics) //TODO: allocate size?
	for i := 0; i < len(chunks); i++ {
		statisticsChunkSummary := <-chunkSummaries
		for stationName, statistics := range statisticsChunkSummary.statisticsByStationName {
			update(stationName, statistics, statisticsByStationName)
		}
	}
	return NewStationTemperatureStatisticsSummary(statisticsByStationName), nil
}

func readChunk(file *os.File, chunk bytes.Chunk) (StationTemperatureStatisticsChunkSummary, error) {
	reader := io.NewSectionReader(file, chunk.StartOffset, chunk.Size)
	buffer := make([]byte, brc.ReadSize)
	statisticsByStationName := NewStatisticsByStationNameMap(capacity)

	offset := 0
	for {
		n, err := reader.Read(buffer[offset:])
		if n > 0 {
			n = n + offset
			leftOver := updateStatisticsIn(buffer[:n], statisticsByStationName)
			if leftOverLength := len(leftOver); leftOverLength > 0 {
				copy(buffer, leftOver)
				offset += leftOverLength
			}
		}
		if err != nil {
			if err == io.EOF {
				return NewStationTemperatureStatisticsChunkSummary(statisticsByStationName.ToGoMap()), nil
			}
			if err != io.EOF {
				return StationTemperatureStatisticsChunkSummary{}, err
			}
		}
	}
}

func updateStatisticsIn(currentBuffer []byte, statisticsByStationName *StatisticsByStationNameMap) []byte {
	for len(currentBuffer) > 0 {
		newLineIndex := bytes2.IndexByte(currentBuffer, '\n')
		if newLineIndex == -1 {
			return currentBuffer
		}
		separatorIndex := -1
		hash := uint64(fnv1aOffset64)
		for index, ch := range currentBuffer[:newLineIndex] {
			if ch == bytes.Separator {
				separatorIndex = index
				break
			}
			hash ^= uint64(ch)
			hash *= fnv1aPrime64
		}

		temperature := bytes.ToTemperature(currentBuffer[separatorIndex+1 : newLineIndex])
		statistics := statisticsByStationName.Get(hash, currentBuffer[:separatorIndex])

		if statistics.totalEntries == 0 {
			statistics.minTemperature = temperature
			statistics.maxTemperature = temperature
			statistics.aggregateTemperature = int64(temperature)
			statistics.totalEntries = 1
		} else {
			statistics.minTemperature = min(statistics.minTemperature, temperature)
			statistics.maxTemperature = max(statistics.maxTemperature, temperature)
			statistics.aggregateTemperature = statistics.aggregateTemperature + int64(temperature)
			statistics.totalEntries += 1
			statisticsByStationName.entryCount += 1
		}
		if newLineIndex+1 < len(currentBuffer) {
			currentBuffer = currentBuffer[newLineIndex+1:]
		} else {
			break
		}
	}
	return currentBuffer
}

func update(stationName string, summary *StationTemperatureStatistics, statisticsByStationName map[string]*StationTemperatureStatistics) {
	existingStatistics, ok := statisticsByStationName[stationName]
	if !ok {
		statisticsByStationName[stationName] = summary
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
	station       [100]byte
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

func (statisticsByStationName *StatisticsByStationNameMap) ToGoMap() map[string]*StationTemperatureStatistics {
	statistics := make(map[string]*StationTemperatureStatistics, statisticsByStationName.entryCount)
	for index := range statisticsByStationName.entries {
		entry := &statisticsByStationName.entries[index]
		if entry.statistics.totalEntries > 0 {
			statistics[string(entry.station[:entry.stationLength])] = &entry.statistics
		}
	}
	return statistics
}
