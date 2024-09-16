package event

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const timeFormat = "2006-01-02T15:04:05.999Z"

type LogEntry struct {
	EventName string `json:"eventName"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`
	ID        int    `json:"id"`
}

type Event struct {
	Name      string
	StartTime time.Time
	Duration  time.Duration
	Children  []*Event
}

func ProcessLogFile(reader io.Reader) ([]*Event, error) {
	scanner := bufio.NewScanner(reader)
	eventMap := make(map[int][]*Event)
	for scanner.Scan() {
		var entry LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, fmt.Errorf("failed to unmarshal log entry: %v", err)
		}
		if err := processLogEntry(&entry, eventMap); err != nil {
			return nil, fmt.Errorf("failed to process log entry: %v", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading log file: %v", err)
	}

	// Return the root events (those without parents)
	var rootEvents []*Event
	for _, events := range eventMap {
		for _, event := range events {
			if event.Duration > 0 {
				rootEvents = append(rootEvents, event)
			}
		}
	}
	return rootEvents, nil
}

func processLogEntry(entry *LogEntry, eventMap map[int][]*Event) error {
	timestamp, err := time.Parse(timeFormat, entry.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to parse timestamp: %v", err)
	}

	events, exists := eventMap[entry.ID]
	if !exists {
		events = []*Event{}
		eventMap[entry.ID] = events
	}

	if entry.Type == "begin" {
		newEvent := &Event{Name: entry.EventName, StartTime: timestamp}
		if len(events) > 0 {
			parent := findLastIncompleteEvent(events)
			if parent != nil {
				parent.Children = append(parent.Children, newEvent)
			}
		}
		events = append(events, newEvent)
		eventMap[entry.ID] = events
	} else if entry.Type == "end" {
		if len(events) > 0 {
			event := findLastIncompleteEvent(events)
			if event != nil && event.Name == entry.EventName {
				event.Duration = timestamp.Sub(event.StartTime)
			}
		}
	}
	return nil
}

func findLastIncompleteEvent(events []*Event) *Event {
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Duration == 0 {
			return events[i]
		}
	}
	return nil
}
