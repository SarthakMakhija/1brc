package bench

import (
	brc "1brc"
	"bufio"
	"os"
	"testing"
)

/*
**
go test -run none -bench . -benchtime 10s -count 7 | tee bench_44k.txt
benchstat bench_44k.txt

goos: linux
goarch: amd64
pkg: 1brc/bench
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkParseWeatherStations44K-16    	 6535921	      1743 ns/op
BenchmarkParseWeatherStations44K-16    	 6775304	      1698 ns/op
BenchmarkParseWeatherStations44K-16    	 6924301	      1688 ns/op
BenchmarkParseWeatherStations44K-16    	 6917014	      1685 ns/op
BenchmarkParseWeatherStations44K-16    	 6888340	      1692 ns/op
BenchmarkParseWeatherStations44K-16    	 6947594	      1689 ns/op
BenchmarkParseWeatherStations44K-16    	 6912704	      1697 ns/op

benchstat bench_44k.txt
goos: linux
goarch: amd64
pkg: 1brc/bench
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

	│ bench_44k.txt │
	│    sec/op     │

ParseWeatherStations44K-16     1.692µ ± 3%
**
*/
func BenchmarkParseWeatherStations44K(b *testing.B) {
	file, err := os.Open("../fixture/44K_weather_stations.csv")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = file.Close()
	}()
	reader := bufio.NewReader(file)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := brc.Parse(reader)
		if err != nil {
			panic(err)
		}
	}
}
