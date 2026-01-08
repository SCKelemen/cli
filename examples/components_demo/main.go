package main

import (
	"fmt"
	"time"

	"github.com/SCKelemen/cli/components"
	"github.com/SCKelemen/dataviz"
)

func main() {
	fmt.Println("=== DataViz Components Demo ===\n")

	// 1. Heatmap Component
	fmt.Println("1. Heatmap Component:")
	heatmapData := createSampleHeatmap()
	heatmap := components.NewHeatmap(heatmapData).
		WithSize(50, 1).
		WithColor("#39d353").
		WithTheme("default")

	heatmapNode := heatmap.ToStyledNode()
	fmt.Println(heatmapNode.Content)
	fmt.Println()

	// 2. StatCard Component
	fmt.Println("2. StatCard Component:")
	statData := dataviz.StatCardData{
		Title:    "Total Commits",
		Value:    "1,234",
		Subtitle: "past month",
		TrendData: []dataviz.TimeSeriesData{
			{Value: 5}, {Value: 8}, {Value: 12},
			{Value: 7}, {Value: 15}, {Value: 20},
			{Value: 18}, {Value: 22}, {Value: 25}, {Value: 30},
		},
	}
	statCard := components.NewStatCard(statData).
		WithSize(50, 10).
		WithColor("#3B82F6").
		WithTheme("default")

	statNode := statCard.ToStyledNode()
	fmt.Println(statNode.Content)
	fmt.Println()

	// 3. BarChart Component
	fmt.Println("3. BarChart Component:")
	barData := dataviz.BarChartData{
		Bars: []dataviz.BarData{
			{Value: 85, Label: "React"},
			{Value: 70, Label: "Vue"},
			{Value: 95, Label: "Svelte"},
			{Value: 60, Label: "Angular"},
			{Value: 80, Label: "Next.js"},
		},
	}
	barChart := components.NewBarChart(barData).
		WithSize(60, 7).
		WithColor("#3B82F6").
		WithTheme("default")

	barNode := barChart.ToStyledNode()
	fmt.Println(barNode.Content)
	fmt.Println()

	// 4. LineGraph Component
	fmt.Println("4. LineGraph Component:")
	lineData := dataviz.LineGraphData{
		Points: []dataviz.TimeSeriesData{
			{Value: 10}, {Value: 25}, {Value: 15},
			{Value: 35}, {Value: 30}, {Value: 45},
			{Value: 40}, {Value: 60}, {Value: 55},
			{Value: 70}, {Value: 65}, {Value: 80},
		},
	}
	lineGraph := components.NewLineGraph(lineData).
		WithSize(60, 12).
		WithColor("#3B82F6").
		WithTheme("default")

	lineNode := lineGraph.ToStyledNode()
	fmt.Println(lineNode.Content)

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nUsage Tips:")
	fmt.Println("- All components support fluent builder pattern")
	fmt.Println("- Use WithTheme() for: default, midnight, nord, paper, wrapped")
	fmt.Println("- Use WithColor() for custom primary colors")
	fmt.Println("- Use WithSize() to control dimensions")
}

// createSampleHeatmap generates sample heatmap data
func createSampleHeatmap() dataviz.HeatmapData {
	days := make([]dataviz.ContributionDay, 30)
	startDate := time.Now().AddDate(0, 0, -30)
	for i := 0; i < 30; i++ {
		days[i] = dataviz.ContributionDay{
			Date:  startDate.AddDate(0, 0, i),
			Count: (i * 3) % 20,
		}
	}
	return dataviz.HeatmapData{
		Days:      days,
		StartDate: startDate,
		EndDate:   time.Now(),
		Type:      "linear",
	}
}
