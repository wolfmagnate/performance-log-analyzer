package output

import (
	"fmt"
	"strings"
	"time"

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

// TreeNode represents a node in the event tree
type TreeNode struct {
	Event    *event.Event
	Children []*TreeNode
}

func PrintEventTreeWithStats(events []*event.Event, statsMap map[string]*stats.EventStats) {
	if len(events) == 0 {
		fmt.Println("No events to display")
		return
	}
	root := buildEventTree(events)
	printTreeNode(root, statsMap, 0)
}

func buildEventTree(events []*event.Event) *TreeNode {
	if len(events) == 0 {
		return nil
	}
	root := &TreeNode{Event: events[0], Children: []*TreeNode{}}
	for _, e := range events {
		if e.Name == root.Event.Name {
			addEventToTree(root, e)
		}
	}
	return root
}

func addEventToTree(node *TreeNode, currentEvent *event.Event) {
	// node.Event.Name == currentEvent.Name を仮定している
	for _, childEvent := range currentEvent.Children {
		var childNode *TreeNode
		for _, existingChild := range node.Children {
			if existingChild.Event.Name == childEvent.Name {
				childNode = existingChild
				break
			}
		}
		if childNode == nil {
			childNode = &TreeNode{Event: childEvent, Children: []*TreeNode{}}
			node.Children = append(node.Children, childNode)
		}
		addEventToTree(childNode, childEvent)
	}
}

func printTreeNode(node *TreeNode, statsMap map[string]*stats.EventStats, depth int) {
	indent := strings.Repeat("  ", depth)
	eventStats := statsMap[node.Event.Name]

	eventNameColor.Printf("%s%s ", indent, node.Event.Name)
	durationColor.Printf("(%s)\n", formatDuration(node.Event.Duration))

	if eventStats != nil {
		mean := stats.CalculateMean(eventStats)
		stdDev := stats.CalculateStdDev(eventStats)

		statsLabelColor.Printf("%s  Avg: ", indent)
		statsValueColor.Printf("%s, ", formatDuration(mean))
		statsLabelColor.Printf("StdDev: ")
		statsValueColor.Printf("%s, ", formatDuration(stdDev))
		statsLabelColor.Printf("Count: ")
		statsValueColor.Printf("%d\n", eventStats.Count)
	}

	for _, child := range node.Children {
		printTreeNode(child, statsMap, depth+1)
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2f μs", float64(d.Microseconds()))
	}
	return fmt.Sprintf("%.2f ms", float64(d.Microseconds())/1000)
}
