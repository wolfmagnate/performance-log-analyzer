package output

import (
	"strings"

	"github.com/fatih/color"
	"github.com/wolfmagnate/performance-log-analyzer/internal/event"
	"github.com/wolfmagnate/performance-log-analyzer/internal/stats"
)

var (
	eventNameColor  = color.New(color.FgCyan, color.Bold)
	durationColor   = color.New(color.FgYellow)
	statsLabelColor = color.New(color.FgGreen)
	statsValueColor = color.New(color.FgWhite)
)

func PrintEventTreeWithStats(events []*event.Event, statsMap map[string]*stats.EventStats) {
	for _, e := range events {
		printEventWithStats(e, statsMap, 0)
	}
}

func printEventWithStats(e *event.Event, statsMap map[string]*stats.EventStats, depth int) {
	indent := strings.Repeat("  ", depth)
	eventStats := statsMap[e.Name]

	eventNameColor.Printf("%s%s ", indent, e.Name)
	durationColor.Printf("(%d ms)\n", e.Duration)

	mean := stats.CalculateMean(eventStats)
	stdDev := stats.CalculateStdDev(eventStats)

	statsLabelColor.Printf("%s  Avg: ", indent)
	statsValueColor.Printf("%.2f ms, ", mean)
	statsLabelColor.Printf("StdDev: ")
	statsValueColor.Printf("%.2f ms, ", stdDev)
	statsLabelColor.Printf("Count: ")
	statsValueColor.Printf("%d\n", eventStats.Count)

	for _, child := range e.Children {
		printEventWithStats(child, statsMap, depth+1)
	}
}
