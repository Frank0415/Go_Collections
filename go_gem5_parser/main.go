package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Entry struct {
	Name          string
	Value         float64
	Percentage1   float64
	Percentage2   float64
	Description   string
	HasPercentage bool
}

func GetInterest(InterestFile *string) (map[string]bool, int) {
	file, err := os.Open(*InterestFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	InterestMap := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			InterestMap[line] = true
		}
	}

	if i := len(InterestMap); i == 0 {
		return InterestMap, 0
	} else {
		fmt.Println("Found", i, "interested items.")
		return InterestMap, i
	}
}

func parseLine(line string, InterestMap *map[string]bool) (*Entry, bool) {
	var dataPart, commentPart string
	hashIdx := strings.Index(line, "#")
	if hashIdx != -1 {
		dataPart = line[:hashIdx]
		commentPart = strings.TrimSpace(line[hashIdx+1:])
	} else {
		dataPart = line
	}

	fields := strings.Fields(dataPart)
	if len(fields) == 0 {
		return nil, false
	}

	entry := &Entry{
		Description: commentPart,
	}

	i := len(fields) - 1

	isPercentage := func(s string) bool {
		return strings.HasSuffix(s, "%")
	}

	if i >= 0 && isPercentage(fields[i]) {
		val, _ := strconv.ParseFloat(strings.TrimSuffix(fields[i], "%"), 64)
		entry.Percentage2 = val
		entry.HasPercentage = true
		i--
	}

	if i >= 0 && isPercentage(fields[i]) {
		val, _ := strconv.ParseFloat(strings.TrimSuffix(fields[i], "%"), 64)
		entry.Percentage1 = val
		entry.HasPercentage = true
		i--
	}

	if i >= 0 {
		val, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return nil, false
		}
		entry.Value = val
		i--
	} else {
		return nil, false 
	}

	if i >= 0 {
		entry.Name = strings.Join(fields[:i+1], " ")
	} else {
		return nil, false
	}

	if !(*InterestMap)[entry.Name] {
		return nil, false
	}

	return entry, true
}

func Parselines(InterestMap *map[string]bool, StatsFile *string, len int) []Entry {
	file, err := os.Open(*StatsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	AllEntries := make([]Entry, 0, len)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println(strings.TrimSpace(scanner.Text()))
		entry, exist := parseLine(strings.TrimSpace(scanner.Text()), InterestMap)
		if exist {
			AllEntries = append(AllEntries, *entry)
		}
	}
	return AllEntries
}

func DataWriter(OutFile *string, entries []Entry) {
	file, err := os.Create(*OutFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush() 
	for _, entry := range entries {
		if entry.HasPercentage {
			fmt.Fprintf(writer, "%s %f %f%% %f%% %s\n", entry.Name, entry.Value, entry.Percentage1, entry.Percentage2, entry.Description)
			// fmt.Printf("%s %f %f%% %f%% %s\n", entry.Name, entry.Value, entry.Percentage1, entry.Percentage2, entry.Description)
		} else {
			fmt.Fprintf(writer, "%s %f %s\n", entry.Name, entry.Value, entry.Description)
			// fmt.Printf("%s %f %s\n", entry.Name, entry.Value, entry.Description)
		}
	}
}

func main() {
	var InterestFile = flag.String("interest", "", "The (relative path to) file that contain interested data")
	var StatsFile = flag.String("stats", "", "The (relative path to) file that contain stats.txt")
	var OutFile = flag.String("out", "", "The (relative path to) the output file")
	flag.Parse()

	InterestMap, len := GetInterest(InterestFile)
	if len == 0 {
		log.Fatal("No interested items found in the interest file!")
	}

	AllEntries := Parselines(&InterestMap, StatsFile, len)

	DataWriter(OutFile, AllEntries)
}
