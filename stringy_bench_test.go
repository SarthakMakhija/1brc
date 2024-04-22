package brc

import (
	bytes2 "bytes"
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

/*
After using:
buffer := &bytes2.Buffer{}
buffer.Grow(len(stationName) + 3 + 12)

go test -run none -bench Stringify -benchtime 10s -count 6
goos: linux
goarch: amd64
pkg: 1brc
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkStringify-16    	36611558	       287.6 ns/op
BenchmarkStringify-16    	41225746	       290.2 ns/op
BenchmarkStringify-16    	44013812	       288.9 ns/op
BenchmarkStringify-16    	47142757	       281.4 ns/op
BenchmarkStringify-16    	39721114	       286.5 ns/op
BenchmarkStringify-16    	43570278	       294.2 ns/op
*/

/*
After using a common buffer.

go test -run none -bench Stringify -benchtime 10s -count 6
goos: linux
goarch: amd64
pkg: 1brc
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkStringify-16    	43127408	       254.6 ns/op
BenchmarkStringify-16    	50992744	       243.9 ns/op
BenchmarkStringify-16    	47075221	       257.5 ns/op
BenchmarkStringify-16    	45149272	       261.8 ns/op
BenchmarkStringify-16    	48105145	       258.2 ns/op
BenchmarkStringify-16    	49907391	       252.6 ns/op
*/

/*
After using a common pre-sized buffer.

go test -run none -bench Stringify -benchtime 10s -count 6
goos: linux
goarch: amd64
pkg: 1brc
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkStringify-16    	45297440	       246.1 ns/op
BenchmarkStringify-16    	49040347	       247.9 ns/op
BenchmarkStringify-16    	51685412	       246.9 ns/op
BenchmarkStringify-16    	48939277	       249.1 ns/op
BenchmarkStringify-16    	48832756	       250.5 ns/op
BenchmarkStringify-16    	49280673	       241.1 ns/op
*/
func BenchmarkStringify(b *testing.B) {
	statistic := StationTemperatureStatistics{
		minTemperature:     -10.3,
		maxTemperature:     10.8,
		averageTemperature: 5.6,
	}
	stationName := "New Mexico"

	buffer := &bytes2.Buffer{}
	buffer.Grow(printableBufferSize)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = statistic.stringify(stationName, buffer)
	}
}
