package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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
		fmt.Println()
		return InterestMap, i
	}
}

func parseLine(line string, InterestMap *map[string]bool) (*Entry, bool) {
	var dataPart, commentPart string
	hashIdx := strings.Index(line, "#")
	if hashIdx != -1 {
		dataPart = line[:hashIdx]
		commentPart = strings.TrimSpace(line[hashIdx+1:])
	} else if hashIdx = strings.Index(line, "(Unspecified)"); hashIdx != -1 {
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

func Parselines(InterestMap *map[string]bool, StatsFile *string, len int) map[string]Entry {
	file, err := os.Open(*StatsFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	AllEntries := make(map[string]Entry)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println(strings.TrimSpace(scanner.Text()))
		entry, exist := parseLine(strings.TrimSpace(scanner.Text()), InterestMap)
		if exist {
			AllEntries[(*entry).Name] = *entry
		}
	}
	return AllEntries
}
