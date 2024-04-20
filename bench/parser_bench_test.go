package bench

import (
	brc "1brc"
	"bufio"
	"os"
	"testing"
)

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
