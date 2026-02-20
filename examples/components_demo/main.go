package main

import (
	"fmt"
	"time"

	"github.com/SCKelemen/cli/components"
	"github.com/SCKelemen/dataviz"
)

func main() {
	fmt.Print("=== DataViz Components Demo ===\n")

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
	fmt.Println()

	// 5. AreaChart Component
	fmt.Println("5. AreaChart Component:")
	areaData := dataviz.AreaChartData{
		Points: []dataviz.TimeSeriesData{
			{Value: 10}, {Value: 25}, {Value: 15},
			{Value: 35}, {Value: 30}, {Value: 45},
			{Value: 40}, {Value: 60}, {Value: 55},
			{Value: 70}, {Value: 65}, {Value: 80},
		},
		Color:       "#3B82F6",
		FillColor:   "#3B82F6",
		UseGradient: true,
		Smooth:      true,
		Tension:     0.3,
	}
	areaChart := components.NewAreaChart(areaData).
		WithSize(60, 12).
		WithColor("#3B82F6").
		WithTheme("default")

	areaNode := areaChart.ToStyledNode()
	fmt.Println(areaNode.Content)
	fmt.Println()

	// 6. ScatterPlot Component
	fmt.Println("6. ScatterPlot Component:")
	scatterPoints := make([]dataviz.ScatterPoint, 12)
	for i := 0; i < 12; i++ {
		scatterPoints[i] = dataviz.ScatterPoint{
			Date:  time.Now().AddDate(0, 0, i-12),
			Value: 10 + (i * 5) + ((i % 3) * 10),
			Size:  3.0,
		}
	}
	scatterData := dataviz.ScatterPlotData{
		Points:     scatterPoints,
		Color:      "#3B82F6",
		MarkerType: "circle",
		MarkerSize: 3.0,
	}
	scatterPlot := components.NewScatterPlot(scatterData).
		WithSize(60, 12).
		WithColor("#3B82F6").
		WithTheme("default")

	scatterNode := scatterPlot.ToStyledNode()
	fmt.Println(scatterNode.Content)

	fmt.Print("\n=== Demo Complete ===")
	fmt.Print("\nUsage Tips:")
	fmt.Println("- All components support fluent builder pattern")
	fmt.Println("- Use WithTheme() for: default, midnight, nord, paper, wrapped")
	fmt.Println("- Use WithColor() for custom primary colors")
	fmt.Println("- Use WithSize() to control dimensions")
	fmt.Println("- AreaChart supports smooth curves with tension control")
	fmt.Println("- ScatterPlot supports various marker types and sizes")
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
