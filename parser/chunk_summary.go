package parser

type StationTemperatureStatisticsChunkSummary struct {
	stationEntries []Entry
	totalEntries   int
}

func NewStationTemperatureStatisticsChunkSummary(stationEntries []Entry, totalEntries int) StationTemperatureStatisticsChunkSummary {
	return StationTemperatureStatisticsChunkSummary{
		stationEntries: stationEntries,
		totalEntries:   totalEntries,
	}
}
