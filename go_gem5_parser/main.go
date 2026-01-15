
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
		fmt.Println("Found ", i, " interested items.")
		return InterestMap, i
	}
}

func parseLine(line string, InterestMap *map[string]bool) (*Entry, bool) {
	line = strings.TrimSpace(line)
	fields := strings.Fields(line)

	// TODO: More complicated lines
	if len(fields) < 3 || len(fields) > 5 {
		return nil, false
	}

	name := fields[0]
	if !(*InterestMap)[name] {
		return nil, false
	}

	value, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		log.Fatal(err)
	}

	entry := &Entry{Name: name, Value: value}

	thirdfield := fields[2]

	if strings.HasPrefix(thirdfield, "#") {
		// Three-Field line parse
		entry.Description = thirdfield[2:]
		entry.HasPercentage = false
	} else if strings.Compare(thirdfield, "(Unspecified)") == 0 {
		entry.Description = "Unspecified"
		entry.HasPercentage = false
	} else {
		if len(fields) != 5 {
			// TODO: Better err handling
			return nil, false
		}
		// Five-Field line parse
		entry.HasPercentage = true
		p1, _ := strconv.ParseFloat(strings.TrimSuffix(thirdfield, "%"), 64)
		if err != nil {
			log.Fatal(err)
		}
		entry.Percentage1 = p1

		p2Str := fields[3]
		p2, _ := strconv.ParseFloat(strings.TrimSuffix(p2Str, "%"), 64)
		if err != nil {
			log.Fatal(err)
		}
		entry.Percentage2 = p2
		entry.Description = fields[4]
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
		entry, exist := parseLine(strings.TrimSpace(scanner.Text()), InterestMap)
		if exist {
			AllEntries = append(AllEntries, *entry)
		}
	}
	return AllEntries
}

func GetEntry(decode []string) (Entry, error) {
	var ent Entry
	ent.Name = "a"
	ent.Value = 1
	ent.Percentage1 = 1
	ent.Percentage2 = 2
	ent.Description = "A"
	ent.HasPercentage = false

	return ent, nil
}

func main() {
	var InterestFile = flag.String("interest", "", "The (relative path to) file that contain interested data")
	var StatsFile = flag.String("stats", "", "The (relative path to) file that contain stats.txt")
	// var OutFile = flag.String("out", "", "The (relative path to) the output file")
	flag.Parse()

	InterestMap, len := GetInterest(InterestFile)
	if len == 0 {
		log.Fatal("No interested items found in the interest file!")
	}

	Parselines(&InterestMap, StatsFile, len)
}
