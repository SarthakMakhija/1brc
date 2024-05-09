package parser

import (
	brc "1brc"
	"1brc/bytes"
	bytes2 "bytes"
	"io"
	"os"
	"runtime"
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
	statisticsByStationName := make(map[string]*StationTemperatureStatistics, 10_0000)
	buffer := make([]byte, brc.ReadSize)

	reader := io.NewSectionReader(file, chunk.StartOffset, chunk.Size)
	var err error
	var n int
	var offset int

	for {
		n, err = reader.Read(buffer[offset:])
		if n > 0 {
			n = n + offset
			var last int
			for index := range buffer[:n] {
				if buffer[index] == '\n' {
					stationName, temperature, err := bytes.SplitIntoStationNameAndTemperature(buffer[last:index])
					if err != nil {
						panic(err)
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
			if err == io.EOF {
				return NewStationTemperatureStatisticsChunkSummary(statisticsByStationName), nil
			}
			if err != io.EOF {
				return StationTemperatureStatisticsChunkSummary{}, err
			}
		}
	}
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
	entries  []Entry
	capacity int
	mask     uint64
}

func NewStatisticsByStationNameMap(capacity int) *StatisticsByStationNameMap {
	return &StatisticsByStationNameMap{
		entries:  make([]Entry, capacity),
		capacity: capacity,
		mask:     uint64(capacity - 1),
	}
}

func (statisticsByStationName StatisticsByStationNameMap) Get(hash uint64, stationName []byte) *StationTemperatureStatistics {
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
