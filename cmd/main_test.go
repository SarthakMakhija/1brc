package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test1BrcStatistics(t *testing.T) {
	outputBuffer := &bytes.Buffer{}
	outputDevice = outputBuffer

	print1brcStatisticsV3("../fixture/10_weather_stations.txt")

	assert.Equal(t,
		"{Canberra:7.0/7.0/7.0;Halifax:10.0/10.0/10.0;Mogadishu:6.0/48.0/90.0;Odesa:6.0/27.0/66.0;Tirana:3.0/6.0/9.0;}",
		string(outputBuffer.Bytes()),
	)
}
