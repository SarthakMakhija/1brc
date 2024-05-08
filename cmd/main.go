package main

import (
	"1brc/parser"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
)

func main() {
	flag.Parse()
	assertFileNameV3()
	if mayBeStartCpuProfileV3() {
		defer pprof.StopCPUProfile()
	}
	print1brcStatisticsV3(*fileName)
}

func print1brcStatisticsV3(fileName string) {
	_, _ = io.WriteString(outputDevice, printableResultV3(parseV3(fileName)))
}

func printableResultV3(result parser.StationTemperatureStatisticsSummary) string {
	return result.PrintableResult()
}

func parseV3(fileName string) parser.StationTemperatureStatisticsSummary {
	temperatureStatisticsResult, err := parser.ParseV3(fileName)

	if err != nil {
		panic(fmt.Errorf("error parsing the file %v, %v", fileName, err))
	}
	return temperatureStatisticsResult
}

func assertFileNameV3() {
	if *fileName == "" {
		panic("-f flag is required")
	}
}

func mayBeStartCpuProfileV3() bool {
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
