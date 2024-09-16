package event

import (
	"bufio"
	"encoding/json"
	"io"
)

type LogEntry struct {
	EventName string `json:"eventName"`
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	ID        int    `json:"id"`
}

type Event struct {
	Name     string
	Duration int64
	Children []*Event
}

func ProcessLogFile(reader io.Reader) ([]*Event, error) {
	scanner := bufio.NewScanner(reader)
	eventMap := make(map[int][]*Event)

	for scanner.Scan() {
		var entry LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, err
		}
		processLogEntry(&entry, eventMap)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
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

func processLogEntry(entry *LogEntry, eventMap map[int][]*Event) {
	events, exists := eventMap[entry.ID]
	if !exists {
		events = []*Event{}
		eventMap[entry.ID] = events
	}

	if entry.Type == "begin" {
		newEvent := &Event{Name: entry.EventName}
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
				event.Duration = entry.Timestamp - events[len(events)-1].Duration
			}
		}
	}
}

func findLastIncompleteEvent(events []*Event) *Event {
	for i := len(events) - 1; i >= 0; i-- {
		if events[i].Duration == 0 {
			return events[i]
		}
	}
	return nil
}
