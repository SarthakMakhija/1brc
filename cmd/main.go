package main

import (
	brc "1brc"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strings"
)

var fileName = flag.String("f", "", "file name")
var cpuProfileFileName = flag.String("cpuprofile", "", "write cpu profile to file")
var outputDevice io.Writer = os.Stdout

func main() {
	flag.Parse()
	assertFileName()
	if mayBeStartCpuProfile() {
		defer pprof.StopCPUProfile()
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

func printableResult(result brc.StationTemperatureStatisticsResult) string {
	output := strings.Builder{}
	output.WriteString("{")

	stationNames := result.AllStationsSorted()
	for _, stationName := range stationNames {
		statistic, _ := result.Get(stationName)
		output.WriteString(statistic.Stringify(stationName))
		output.WriteString(";")
	}
	output.WriteString("}")
	return output.String()
}

func parse(file *os.File) brc.StationTemperatureStatisticsResult {
	temperatureStatisticsResult, err := brc.Parse(bufio.NewReader(file))
	if err != nil {
		panic(fmt.Errorf("error parsing the file %v, %v", *fileName, err))
	}
	return temperatureStatisticsResult
}

func assertFileName() {
	if *fileName == "" {
		panic("-f flag is required")
	}
}

func mayBeStartCpuProfile() bool {
	if *cpuProfileFileName != "" {
		cpuProfileFile, err := os.Create(*cpuProfileFileName)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(cpuProfileFile)
		return true
	}
	return false
}
