package stats

import (
	"math"

	"github.com/wolfmagnate/performance-log-analyzer/internal/event"
)

type EventStats struct {
	TotalTime  int64
	Count      int
	SquaredSum int64
}

func CalculateStats(events []*event.Event) map[string]*EventStats {
	statsMap := make(map[string]*EventStats)
	for _, e := range events {
		calculateEventStats(e, statsMap)
	}
	return statsMap
}

func calculateEventStats(e *event.Event, statsMap map[string]*EventStats) {
	stats, exists := statsMap[e.Name]
	if !exists {
		stats = &EventStats{}
		statsMap[e.Name] = stats
	}

	stats.TotalTime += e.Duration
	stats.Count++
	stats.SquaredSum += e.Duration * e.Duration

	for _, child := range e.Children {
		calculateEventStats(child, statsMap)
	}
}

func CalculateMean(stats *EventStats) float64 {
	if stats.Count == 0 {
		return 0
	}
	return float64(stats.TotalTime) / float64(stats.Count)
}

func CalculateStdDev(stats *EventStats) float64 {
	if stats.Count == 0 {
		return 0
	}
	mean := CalculateMean(stats)
	variance := float64(stats.SquaredSum)/float64(stats.Count) - mean*mean
	return math.Sqrt(variance)
}
