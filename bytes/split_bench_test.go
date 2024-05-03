package bytes

import (
	"fmt"
	"testing"
)

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

/*
After changing SplitIntoStationNameAndTemperature to iterate from len-1 to 0.

BenchmarkSplitIntoStationNameAndTemperature-16                               	1000000000	         5.586 ns/op
BenchmarkSplitIntoStationNameAndTemperature-16                               	1000000000	         5.561 ns/op
BenchmarkSplitIntoStationNameAndTemperature-16                               	1000000000	         5.543 ns/op
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

/*
Baseline with SplitIntoStationNameAndTemperature iterating from index 0 to len-1.
go test -run none -bench SplitIntoStationNameAndTemperature -benchtime 20s -count 6

goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	262566970	        44.85 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	267653738	        45.28 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	267477218	        45.41 ns/op
*/

/*
After changing SplitIntoStationNameAndTemperature to iterate from len-1 to 0.

BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         5.543 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         5.656 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         5.549 ns/op
*/

/*
After changing SplitIntoStationNameAndTemperature to handle conversion to temperature also.

BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         3.465 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         3.480 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         3.551 ns/op
*/

/*
After changing SplitIntoStationNameAndTemperature to reduce the cost of multiplying the
variable `multiplier` by 10.

BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         3.110 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         2.964 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         3.002 ns/op
*/

/*
After changing SplitIntoStationNameAndTemperature to have zero bound check.

go test -run none -bench BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B -benchtime 20s -count 6 .
goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         2.680 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         2.662 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         2.646 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         2.653 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         2.644 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B-16    	1000000000	         2.673 ns/op
*/
func BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100B(b *testing.B) {
	line := []byte(fmt.Sprintf("%v;%v", stationName(100), -99.9))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := SplitIntoStationNameAndTemperature(line)
		if err != nil {
			panic(err)
		}
	}
}

/*
go test -run none -bench BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndPositiveTemperatureWithSingleDigitBeforeDot -benchtime 20s -count 6 .
goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndPositiveTemperatureWithSingleDigitBeforeDot-16    	1000000000	         1.724 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndPositiveTemperatureWithSingleDigitBeforeDot-16    	1000000000	         1.718 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndPositiveTemperatureWithSingleDigitBeforeDot-16    	1000000000	         1.771 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndPositiveTemperatureWithSingleDigitBeforeDot-16    	1000000000	         1.758 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndPositiveTemperatureWithSingleDigitBeforeDot-16    	1000000000	         1.717 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndPositiveTemperatureWithSingleDigitBeforeDot-16    	1000000000	         1.732 ns/op
*/
func BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndPositiveTemperatureWithSingleDigitBeforeDot(b *testing.B) {
	line := []byte(fmt.Sprintf("%v;%v", stationName(100), 9.9))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := SplitIntoStationNameAndTemperature(line)
		if err != nil {
			panic(err)
		}
	}
}

/*
go test -run none -bench BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndNegativeTemperatureWithSingleDigitBeforeDot -benchtime 20s -count 6 .
goos: linux
goarch: amd64
pkg: 1brc/bytes
cpu: 13th Gen Intel(R) Core(TM) i7-1360P

BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndNegativeTemperatureWithSingleDigitBeforeDot-16    	1000000000	         2.055 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndNegativeTemperatureWithSingleDigitBeforeDot-16    	1000000000	         2.084 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndNegativeTemperatureWithSingleDigitBeforeDot-16    	1000000000	         2.062 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndNegativeTemperatureWithSingleDigitBeforeDot-16    	1000000000	         2.107 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndNegativeTemperatureWithSingleDigitBeforeDot-16    	1000000000	         2.057 ns/op
BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndNegativeTemperatureWithSingleDigitBeforeDot-16    	1000000000	         2.084 ns/op
*/
func BenchmarkSplitIntoStationNameAndTemperatureWithStationNameAsLongAs100BAndNegativeTemperatureWithSingleDigitBeforeDot(b *testing.B) {
	line := []byte(fmt.Sprintf("%v;%v", stationName(100), -9.9))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, err := SplitIntoStationNameAndTemperature(line)
		if err != nil {
			panic(err)
		}
	}
}

func stationName(length int) string {
	station := make([]byte, 0, length)
	for index := 0; index < length; index++ {
		station = append(station, 'a')
	}
	return string(station)
}
