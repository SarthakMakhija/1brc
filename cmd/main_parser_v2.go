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

var fileName = flag.String("f", "", "file name")
var cpuProfileFileName = flag.String("cpuprofile", "", "write cpu profile to file")
var outputDevice io.Writer = os.Stdout

func mainParserV2() {
	flag.Parse()
	assertFileNameV2()
	if mayBeStartCpuProfileV3() {
		defer pprof.StopCPUProfile()
	}
	print1brcStatisticsV2(*fileName)
}

func print1brcStatisticsV2(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("error opening the file %v, %v", fileName, err))
	}
	defer func() {
		_ = file.Close()
	}()
	_, _ = io.WriteString(outputDevice, printableResultV2(parseV2(file)))
}

func printableResultV2(result parser.StationTemperatureStatisticsSummary) string {
	return result.PrintableResult()
}

func parseV2(file *os.File) parser.StationTemperatureStatisticsSummary {
	temperatureStatisticsResult, err := parser.ParseV2(file)

	if err != nil {
		panic(fmt.Errorf("error parsing the file %v, %v", *fileName, err))
	}
	return temperatureStatisticsResult
}

func assertFileNameV2() {
	if *fileName == "" {
		panic("-f flag is required")
	}
}

func mayBeStartCpuProfileV2() bool {
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
