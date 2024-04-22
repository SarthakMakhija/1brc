package brc

import (
	"testing"
)

/*
go test -run none -bench Stringify -benchtime 10s -count 6 | tee stringify_baseline.txt
goos: linux
goarch: amd64
pkg: 1brc
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkStringify-16    	34292252	       325.0 ns/op
BenchmarkStringify-16    	36423590	       328.8 ns/op
BenchmarkStringify-16    	36095724	       325.2 ns/op
BenchmarkStringify-16    	37095597	       324.5 ns/op
BenchmarkStringify-16    	37372545	       324.4 ns/op
BenchmarkStringify-16    	36647731	       323.9 ns/op

benchstat stringify_baseline.txt

goos: linux
goarch: amd64
pkg: 1brc
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

	│ stringify_baseline.txt │
	│      sec/op      │

Stringify-16        324.8n ± 1%
*/
func BenchmarkStringify(b *testing.B) {
	statistic := StationTemperatureStatistics{
		minTemperature:     -10.3,
		maxTemperature:     10.8,
		averageTemperature: 5.6,
	}
	stationName := "New Mexico"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = statistic.stringify(stationName)
	}
}
