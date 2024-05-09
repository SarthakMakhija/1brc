package parser

type StationTemperatureStatisticsChunkSummary struct {
	statisticsByStationName *StatisticsByStationNameMap
}

func NewStationTemperatureStatisticsChunkSummary(statisticsByStationName *StatisticsByStationNameMap) StationTemperatureStatisticsChunkSummary {
	return StationTemperatureStatisticsChunkSummary{
		statisticsByStationName: statisticsByStationName,
	}
}
