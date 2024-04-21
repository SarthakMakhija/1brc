package bytes

import (
	"testing"
)

var GlobalSink float64

/*
Benchtime had to be increased and the benchmark had to be changed from a single toFloat64 to a loop, to
get the benchstat variance in range.

go test -run none -bench . -benchtime 20s -count 6 -cpuprofile convert_to_float_64_cpu_baseline.out | tee convert_to_float_64_basline.txt

goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkConvertToFloat64-16    	22524248	      2079 ns/op
BenchmarkConvertToFloat64-16    	11464390	      2052 ns/op
BenchmarkConvertToFloat64-16    	11422123	      2116 ns/op
BenchmarkConvertToFloat64-16    	10681839	      2065 ns/op
BenchmarkConvertToFloat64-16    	11560833	      2122 ns/op
BenchmarkConvertToFloat64-16    	11359008	      2092 ns/op

These results indicate, 2100ns for 200 string to float64 conversions. This means almost 10-11ns per string to
float64 conversion.

benchstat convert_to_float_64_basline.txt

goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

	│ convert_to_float_64_basline.txt │
	│             sec/op              │

ConvertToFloat64-16                       2.086µ ± 2%
*/
func BenchmarkConvertToFloat64(b *testing.B) {
	var localSink float64
	for i := 0; i < b.N; i++ {
		for count := 1; count <= 200; count++ {
			result, err := toFloat64("-10.443")
			if err != nil {
				panic(err)
			}
			localSink = result
		}
	}
	GlobalSink = localSink
}
