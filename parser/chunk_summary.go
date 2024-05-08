package parser

type StationTemperatureStatisticsChunkSummary struct {
	statisticsByStationName map[string]*StationTemperatureStatistics
}

func NewStationTemperatureStatisticsChunkSummary(statisticsByStationName map[string]*StationTemperatureStatistics) StationTemperatureStatisticsChunkSummary {
	return StationTemperatureStatisticsChunkSummary{
		statisticsByStationName: statisticsByStationName,
	}
}
