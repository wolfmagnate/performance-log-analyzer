package stats

import (
	"math"
	"time"

	"github.com/wolfmagnate/performance-log-analyzer/internal/event"
)

type EventStats struct {
	TotalTime  time.Duration
	Count      int
	SquaredSum time.Duration
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
}

func CalculateMean(stats *EventStats) time.Duration {
	if stats.Count == 0 {
		return 0
	}
	return time.Duration(int64(stats.TotalTime) / int64(stats.Count))
}

func CalculateStdDev(stats *EventStats) time.Duration {
	if stats.Count == 0 {
		return 0
	}
	mean := CalculateMean(stats)
	variance := float64(stats.SquaredSum)/float64(stats.Count) - math.Pow(float64(mean), 2)
	return time.Duration(math.Sqrt(float64(variance)))
}
