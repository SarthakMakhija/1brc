package bytes

import (
	"testing"
)

var GlobalSink float64

/*
Benchtime had to be increased and the benchmark had to be changed from a single ToTemperature to a loop, to
get the benchstat variance in range.

Originally ToTemperature was taking string input.

go test -run none -bench . -benchtime 20s -count 6 -cpuprofile convert_to_float_64_cpu_baseline.out | tee convert_to_float_64_basline.txt

goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkConvertToToTemperature-16    	22524248	      2079 ns/op
BenchmarkConvertToToTemperature-16    	11464390	      2052 ns/op
BenchmarkConvertToToTemperature-16    	11422123	      2116 ns/op
BenchmarkConvertToToTemperature-16    	10681839	      2065 ns/op
BenchmarkConvertToToTemperature-16    	11560833	      2122 ns/op
BenchmarkConvertToToTemperature-16    	11359008	      2092 ns/op

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

/*
After removing redundant conditions and changing ToTemperature to operate on byte slice.

go test -run none -bench BenchmarkConvertToToTemperature -benchtime 20s -count 6 -cpuprofile convert_to_float_64_cpu_operate_on_byte_slice.out
goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkConvertToToTemperature-16    	26882529	      1795 ns/op
BenchmarkConvertToToTemperature-16    	14356762	      1727 ns/op
BenchmarkConvertToToTemperature-16    	13436760	      1743 ns/op
BenchmarkConvertToToTemperature-16    	13323096	      1758 ns/op
BenchmarkConvertToToTemperature-16    	13637918	      1749 ns/op
BenchmarkConvertToToTemperature-16    	13335798	      1757 ns/op
*/

/*
After changing convert to handle only one fractional digit.

go test -run none -bench BenchmarkConvertToToTemperature -benchtime 20s -count 6 -cpuprofile convert_to_float_64_cpu_handle_one_fractional_digit.out
goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkConvertToToTemperature-16    	36951853	      1159 ns/op
BenchmarkConvertToToTemperature-16    	15636702	      1539 ns/op
BenchmarkConvertToToTemperature-16    	15605932	      1509 ns/op
BenchmarkConvertToToTemperature-16    	15333734	      1556 ns/op
BenchmarkConvertToToTemperature-16    	15502268	      1550 ns/op
BenchmarkConvertToToTemperature-16    	15492415	      1552 ns/op
*/

/*
After removing the cost of uint16 conversion.

go test -run none -bench BenchmarkConvertToToTemperature -benchtime 20s -count 6 -cpuprofile convert_to_float_64_cpu_handle_cost_of_uint16.out
goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkConvertToToTemperature-16    	23721129	       993.1 ns/op
BenchmarkConvertToToTemperature-16    	23805321	      1004 ns/op
BenchmarkConvertToToTemperature-16    	18496376	      1196 ns/op
BenchmarkConvertToToTemperature-16    	19726849	      1209 ns/op
BenchmarkConvertToToTemperature-16    	19353622	      1176 ns/op
BenchmarkConvertToToTemperature-16    	19573942	      1159 ns/op
*/

/*
After removing the redundant if check for minus in integerPart and presence of . in convert.

go test -run none -bench ConvertToFloat64  -benchtime 10s -count 6
goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkConvertToToTemperature-16    	18565965	       647.2 ns/op
BenchmarkConvertToToTemperature-16    	18382945	       646.0 ns/op
BenchmarkConvertToToTemperature-16    	18581326	       644.5 ns/op
BenchmarkConvertToToTemperature-16    	18651205	       645.1 ns/op
BenchmarkConvertToToTemperature-16    	18560593	       644.5 ns/op
BenchmarkConvertToToTemperature-16    	18618092	       645.4 ns/op
*/
func BenchmarkConvertToToTemperature(b *testing.B) {
	var localSink float64
	input := []byte("-10.443")
	for i := 0; i < b.N; i++ {
		for count := 1; count <= 200; count++ {
			result, err := ToTemperature(input)
			if err != nil {
				panic(err)
			}
			localSink = result
		}
	}
	GlobalSink = localSink
}
