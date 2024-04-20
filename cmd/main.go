package main

import (
	brc "1brc"
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var fileName = flag.String("f", "", "file name")
var outputDevice io.Writer = os.Stdout

func main() {
	flag.Parse()
	if *fileName == "" {
		_, _ = fmt.Fprintln(os.Stderr, "-f flag is required")
		return
	}

	print1brcStatistics(*fileName)
}

func print1brcStatistics(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("error opening the file %v, %v", fileName, err))
	}
	defer func() {
		_ = file.Close()
	}()
	_, _ = io.WriteString(outputDevice, printableResult(parse(file)))
}

func parse(file *os.File) brc.StationTemperatureStatisticsResult {
	temperatureStatisticsResult, err := brc.Parse(bufio.NewReader(file))
	if err != nil {
		panic(fmt.Errorf("error parsing the file %v, %v", *fileName, err))
	}
	return temperatureStatisticsResult
}

func printableResult(result brc.StationTemperatureStatisticsResult) string {
	output := strings.Builder{}
	output.WriteString("{")

	statisticsByStationName := result.Iterator()
	for statisticsByStationName.Next() {
		statistic := statisticsByStationName.Value().(brc.StationTemperatureStatistics)
		output.WriteString(statistic.Stringify(statisticsByStationName.Key().(string)))
		output.WriteString(";")
	}
	output.WriteString("}")
	return output.String()
}
