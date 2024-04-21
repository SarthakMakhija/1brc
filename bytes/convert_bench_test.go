package bytes

import (
	"testing"
)

var GlobalSink float64

/*
Benchtime had to be increased and the benchmark had to be changed from a single toFloat64 to a loop, to
get the benchstat variance in range.

Originally toFloat64 was taking string input.

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

/*
After removing redundant conditions and changing toFloat64 to operate on byte slice.

go test -run none -bench BenchmarkConvertToFloat64 -benchtime 20s -count 6 -cpuprofile convert_to_float_64_cpu_operate_on_byte_slice.out
goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkConvertToFloat64-16    	26882529	      1795 ns/op
BenchmarkConvertToFloat64-16    	14356762	      1727 ns/op
BenchmarkConvertToFloat64-16    	13436760	      1743 ns/op
BenchmarkConvertToFloat64-16    	13323096	      1758 ns/op
BenchmarkConvertToFloat64-16    	13637918	      1749 ns/op
BenchmarkConvertToFloat64-16    	13335798	      1757 ns/op
*/

/*
After changing convert to handle only one fractional digit.

go test -run none -bench BenchmarkConvertToFloat64 -benchtime 20s -count 6 -cpuprofile convert_to_float_64_cpu_handle_one_fractional_digit.out
goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P
BenchmarkConvertToFloat64-16    	36951853	      1159 ns/op
BenchmarkConvertToFloat64-16    	15636702	      1539 ns/op
BenchmarkConvertToFloat64-16    	15605932	      1509 ns/op
BenchmarkConvertToFloat64-16    	15333734	      1556 ns/op
BenchmarkConvertToFloat64-16    	15502268	      1550 ns/op
BenchmarkConvertToFloat64-16    	15492415	      1552 ns/op
*/
func BenchmarkConvertToFloat64(b *testing.B) {
	var localSink float64
	input := []byte("-10.443")
	for i := 0; i < b.N; i++ {
		for count := 1; count <= 200; count++ {
			result, err := toFloat64(input)
			if err != nil {
				panic(err)
			}
			localSink = result
		}
	}
	GlobalSink = localSink
}
