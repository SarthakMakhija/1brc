package bytes

import "testing"

/*
After optimizations in convert (till Round8).
go test -run none -bench SplitIntoStationNameAndTemperature -benchtime 20s -count 6

goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkSplitIntoStationNameAndTemperature-16    	1000000000	         5.832 ns/op
BenchmarkSplitIntoStationNameAndTemperature-16    	1000000000	         5.812 ns/op
BenchmarkSplitIntoStationNameAndTemperature-16    	1000000000	         5.811 ns/op
BenchmarkSplitIntoStationNameAndTemperature-16    	1000000000	         5.810 ns/op
BenchmarkSplitIntoStationNameAndTemperature-16    	1000000000	         5.811 ns/op
BenchmarkSplitIntoStationNameAndTemperature-16    	1000000000	         5.877 ns/op
*/
func BenchmarkSplitIntoStationNameAndTemperature(b *testing.B) {
	line := []byte("Odesa;-10.3")
	for i := 0; i < b.N; i++ {
		_, _, err := SplitIntoStationNameAndTemperature(line)
		if err != nil {
			panic(err)
		}
	}
}
