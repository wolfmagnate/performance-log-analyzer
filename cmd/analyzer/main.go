package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wolfmagnate/performance-log-analyzer/internal/event"
	"github.com/wolfmagnate/performance-log-analyzer/internal/output"
	"github.com/wolfmagnate/performance-log-analyzer/internal/stats"
)

func main() {
	inputFile := flag.String("input", "performance_log.txt", "Input log file path")
	flag.Parse()

	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	events, err := event.ProcessLogFile(file)
	if err != nil {
		fmt.Printf("Error processing log file: %v\n", err)
		os.Exit(1)
	}

	statistics := stats.CalculateStats(events)
	output.PrintEventTreeWithStats(events, statistics)
}
