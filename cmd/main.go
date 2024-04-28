package main

import (
	brc "1brc"
	"1brc/parser"
	"bufio"
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

func printableResult(result parser.StationTemperatureStatisticsSummary) string {
	return result.PrintableResult()
}

func parse(file *os.File) parser.StationTemperatureStatisticsSummary {
	temperatureStatisticsResult, err :=
		parser.ParseV2(bufio.NewReaderSize(file, brc.BufferSize))

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
