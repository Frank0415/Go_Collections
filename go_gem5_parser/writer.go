package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type CsvWriter struct {
	FilePath string
}

type MdWriter struct {
	FilePath string
}

func (w CsvWriter) Write(entries []Entry, filepath *os.File) {
}

func (w MdWriter) Write(entries []Entry, filepath *os.File) {

}

func PrintAllStats(stats *TMAStats) {
	if stats == nil {
		fmt.Println("Stats is nil!")
		return
	}

	p := stats.pmu
	t := stats.mytma

	fmt.Println("==================== Raw TMA Metrics Collected from GEM5 ====================")
	fmt.Printf("  L1_retire:           %.4f\n", t.L1_retire)
	fmt.Printf("  L1_badspec:          %.4f\n", t.L1_badspec)
	fmt.Printf("  L1_frontend:         %.4f\n", t.L1_frontend)
	fmt.Printf("  L1_backend:          %.4f\n", t.L1_backend)
	fmt.Printf("  L0_fullfrontend:     %.4f\n", t.L0_fullfrontend)
	fmt.Printf("  L0_frontendutil:     %.4f\n", t.L0_frontendutil)

	fmt.Println("\n==================== Base Pipeline Stats ====================")
	fmt.Printf("  Cycles:              %d\n", p.Cycles)
	fmt.Printf("  Simticks:            %d\n", p.Simticks)
	fmt.Printf("  InstFetched:         %d\n", p.InstFetched)
	fmt.Printf("  SlotsIssued:         %d\n", p.SlotsIssued)
	fmt.Printf("  SlotsRetired:        %d\n", p.SlotsRetired)

	fmt.Println("\n==================== Execution & Retire Stats ====================")
	fmt.Printf("  OpsExecuted:         %d\n", p.OpsExecuted)
	fmt.Printf("  MispredRetired:      %d\n", p.MispredRetired)

	fmt.Println("\n==================== Pipeline Bubbles (Stalls) ====================")
	fmt.Printf("  FetchBubbles:        %d (I-Cache Wait)\n", p.FetchBubbles)
	fmt.Printf("  RecoveryBubbles:     %d (Squash)\n", p.RecoveryBubbles)
	fmt.Printf("  MachineClears:       %d (Nukes)\n", p.MachineClears)

	fmt.Println("\n==================== Structural Stalls ====================")
	fmt.Printf("  LoadQueueFull:       %d\n", p.LoadQueueFull)
	fmt.Printf("  StoreQueueFull:      %d\n", p.StoreQueueFull)
	fmt.Printf("  InstQueueFull:       %d\n", p.InstQueueFull)
	fmt.Printf("  LSQBlockedByCache:   %d\n", p.LSQBlockedByCache)

	fmt.Println("\n==================== Memory Hierarchy (L1) ====================")
	fmt.Printf("  L1D Access:          %d\n", p.L1D.Access)
	fmt.Printf("  L1D Hits:            %d\n", p.L1D.Hits)
	fmt.Printf("  L1D Misses:          %d\n", p.L1D.Misses)
	fmt.Printf("  ----------------------------\n")
	fmt.Printf("  L1I Access:          %d\n", p.L1I.Access)
	fmt.Printf("  L1I Hits:            %d\n", p.L1I.Hits)
	fmt.Printf("  L1I Misses:          %d\n", p.L1I.Misses)

	fmt.Println("\n==================== Memory Hierarchy (L2 & DRAM) ====================")
	fmt.Printf("  L2 Access:           %d\n", p.L2.Access)
	fmt.Printf("  L2 Hits:             %d\n", p.L2.Hits)
	fmt.Printf("  L2 Misses:           %d\n", p.L2.Misses)
	fmt.Printf("  ----------------------------\n")
	fmt.Printf("  MemReadReqs:         %d\n", p.MemReadReqs)
	fmt.Printf("  MemQueueStallCount:  %d\n", p.MemQueueStallCount)
	fmt.Printf("  MeanLoadAccessTime:  %.2f cycles\n", p.MeanLoadAccessTime)

	fmt.Println("\n==================== Thread 0 Stats ====================")
	fmt.Printf("  Num Insts:           %d\n", p.Thread0.numInsts)
	fmt.Printf("  Num Ops:             %d\n", p.Thread0.numOps)

	fmt.Println("=====================================================================")
}

func WriteData(OutFile *string, entries map[string]Entry, format *string, stats *TMAStats) {
	file, err := os.Create(*OutFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	PrintAllStats(stats)

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