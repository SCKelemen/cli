package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
	"golang.org/x/term"
)

func main() {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	if w == 0 {
		w = 120
	}
	if h == 0 {
		h = 40
	}

	// Seed random for simulated data
	rand.Seed(time.Now().UnixNano())

	screen := renderer.NewScreen(w, h)
	root := buildDashboard(w, h)

	constraints := layout.Tight(float64(w), float64(h))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(w),
		ViewportHeight: float64(h),
		RootFontSize:   16,
	}
	layout.Layout(root.Node, constraints, ctx)
	screen.Render(root)

	fmt.Print(screen.String())
}

func buildDashboard(width, height int) *renderer.StyledNode {
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(width)),
			Height:        layout.Px(float64(height)),
			Padding:       layout.Spacing{Top: layout.Px(1), Right: layout.Px(2), Bottom: layout.Px(1), Left: layout.Px(2)},
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header
	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgBlue, _ := color.ParseColor("oklch(0.5 0.25 250)")
	borderBlue, _ := color.ParseColor("oklch(0.7 0.25 250)")

	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(3),
			Margin:  layout.Spacing{Bottom: layout.Px(1)},
		},
	}
	headerStyle := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgBlue,
		Bold:        true,
		BorderColor: &borderBlue,
		TextAlign:   renderer.TextAlignCenter,
	}
	headerStyle.WithBorder(renderer.ThickBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = fmt.Sprintf("System Dashboard - %s - Uptime: 2h 34m 12s",
		time.Now().Format("15:04:05"))
	rootStyled.AddChild(headerStyled)

	// Resource gauges container
	gaugesContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width - 4)),
			Height:        layout.Px(8),
			Margin:        layout.Spacing{Bottom: layout.Px(1)},
		},
	}
	gaugesStyled := renderer.NewStyledNode(gaugesContainer, nil)

	// Resource gauges
	gauges := []struct {
		label string
		value float64
	}{
		{"CPU", 45.2},
		{"Memory", 68.5},
		{"Disk", 82.1},
		{"Network", 23.7},
	}

	gaugeWidth := (width - 4 - 6) / 4 // Divide by 4 gauges, account for spacing
	for i, g := range gauges {
		gauge := createGauge(g.label, g.value, gaugeWidth)
		gauge.Node.Style.Margin = layout.Spacing{Right: layout.Px(2)}
		if i == len(gauges)-1 {
			gauge.Node.Style.Margin.Right = layout.Px(0)
		}
		gaugesStyled.AddChild(gauge)
	}

	rootStyled.AddChild(gaugesStyled)

	// Bottom section: System Info + Activity
	// Use FlexGrow to fill remaining space instead of explicit height
	bottomContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(width - 4)),
			FlexGrow:      1,
		},
	}
	bottomStyled := renderer.NewStyledNode(bottomContainer, nil)

	// System Info (left)
	sysInfo := createSystemInfo()
	sysInfo.Node.Style.Width = layout.Px(float64((width - 4) / 2))
	sysInfo.Node.Style.FlexGrow = 1
	sysInfo.Node.Style.Margin.Right = layout.Px(1)
	bottomStyled.AddChild(sysInfo)

	// Activity Log (right)
	activityLog := createActivityLog()
	activityLog.Node.Style.Width = layout.Px(float64((width - 4) / 2))
	activityLog.Node.Style.FlexGrow = 1
	activityLog.Node.Style.Margin.Left = layout.Px(1)
	bottomStyled.AddChild(activityLog)

	rootStyled.AddChild(bottomStyled)

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width - 4)),
			Height:  layout.Px(1),
			Margin:  layout.Spacing{Top: layout.Px(1)},
		},
	}
	fgGray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
		TextAlign:  renderer.TextAlignCenter,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = "Dashboard Snapshot"
	rootStyled.AddChild(footerStyled)

	return rootStyled
}

func createGauge(label string, value float64, width int) *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(width)),
			Height:  layout.Px(8),
		},
	}

	borderGray, _ := color.ParseColor("#5A5A5A")
	style := &renderer.Style{
		BorderColor: &borderGray,
	}
	style.WithBorder(renderer.NormalBorder)

	styledNode := renderer.NewStyledNode(node, style)

	// Build content
	content := fmt.Sprintf("%s\n\n%.1f%%\n", label, value)

	// Progress bar
	barWidth := width - 4
	filled := int(float64(barWidth) * value / 100.0)
	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	content += bar + "\n\n"

	// Status
	status := "Normal"
	if value > 90 {
		status = "Critical"
	} else if value > 75 {
		status = "Warning"
	}
	content += status

	styledNode.Content = content

	return styledNode
}

func createSystemInfo() *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display:   layout.DisplayBlock,
			MinHeight: layout.Px(10), // Set minimum height to prevent collapse
		},
	}

	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgGreen, _ := color.ParseColor("oklch(0.4 0.15 150)")
	borderGreen, _ := color.ParseColor("oklch(0.6 0.15 150)")
	style := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgGreen,
		BorderColor: &borderGreen,
		TextWrap:    renderer.TextWrapNormal,
	}
	style.WithBorder(renderer.RoundedBorder)

	styledNode := renderer.NewStyledNode(node, style)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	content := fmt.Sprintf(`System Information

OS: %s
Arch: %s
CPUs: %d
Go Version: %s
Goroutines: %d
Heap Alloc: %.2f MB`,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		runtime.Version(),
		runtime.NumGoroutine(),
		float64(m.Alloc)/1024/1024,
	)

	styledNode.Content = content

	return styledNode
}

func createActivityLog() *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display:   layout.DisplayBlock,
			MinHeight: layout.Px(10), // Set minimum height to prevent collapse
		},
	}

	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgPurple, _ := color.ParseColor("oklch(0.4 0.15 270)")
	borderPurple, _ := color.ParseColor("oklch(0.6 0.15 270)")
	style := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgPurple,
		BorderColor: &borderPurple,
		TextWrap:    renderer.TextWrapNormal,
	}
	style.WithBorder(renderer.RoundedBorder)

	styledNode := renderer.NewStyledNode(node, style)

	activities := []string{
		"System started successfully",
		"Network connection established",
		"Services running normally",
		"Database query completed",
		"Cache updated successfully",
	}

	content := "Recent Activity\n\n"
	timestamp := time.Now().Format("15:04:05")
	for _, activity := range activities {
		content += fmt.Sprintf("[%s] %s\n", timestamp, activity)
	}

	styledNode.Content = content

	return styledNode
}
