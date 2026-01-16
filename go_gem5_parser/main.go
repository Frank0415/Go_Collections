package main

import (
	"flag"
	"log"
	"os"
)

type Entry struct {
	Name          string
	Value         float64
	Percentage1   float64
	Percentage2   float64
	Description   string
	HasPercentage bool
}

type Stats struct {
}

type DataWriter interface {
	Write(entries []Entry, filepath *os.File) error
}

func main() {
	var InterestFile = flag.String("interest", "interests.txt", "The (relative path to) file that contain interested data")
	var StatsFile = flag.String("stats", "m5out/stats.txt", "The (relative path to) file that contain stats.txt")
	var OutFile = flag.String("out", "out.md", "The (relative path to) the output file")
	var Format = flag.String("format", "Markdown", "The default output type of the file")
	flag.Parse()

	InterestMap, len := GetInterest(InterestFile)
	if len == 0 {
		log.Fatal("No interested items found in the interest file!")
	}

	AllEntries := Parselines(&InterestMap, StatsFile, len)

	WriteData(OutFile, AllEntries, Format, GetStats(&AllEntries))
}