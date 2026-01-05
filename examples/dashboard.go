package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
	"golang.org/x/term"
)

type tickMsg time.Time

type model struct {
	width      int
	height     int
	cpuUsage   float64
	memUsage   float64
	diskUsage  float64
	netUsage   float64
	activities []string
	uptime     time.Duration
	startTime  time.Time
}

func initialModel() model {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	if w == 0 {
		w = 120
	}
	if h == 0 {
		h = 40
	}

	return model{
		width:      w,
		height:     h,
		cpuUsage:   45.2,
		memUsage:   68.5,
		diskUsage:  82.1,
		netUsage:   23.7,
		activities: []string{
			"System started successfully",
			"Network connection established",
			"Services running normally",
		},
		startTime: time.Now(),
	}
}

func (m model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		// Simulate resource usage changes
		m.cpuUsage = math.Max(0, math.Min(100, m.cpuUsage+rand.Float64()*10-5))
		m.memUsage = math.Max(0, math.Min(100, m.memUsage+rand.Float64()*5-2.5))
		m.diskUsage = math.Max(0, math.Min(100, m.diskUsage+rand.Float64()*2-1))
		m.netUsage = math.Max(0, math.Min(100, m.netUsage+rand.Float64()*15-7.5))

		// Update uptime
		m.uptime = time.Since(m.startTime)

		// Occasionally add activity
		if rand.Float64() < 0.2 {
			activities := []string{
				"Service health check passed",
				"Database query completed",
				"API request processed",
				"Cache updated successfully",
				"Backup job completed",
				"Log rotation performed",
			}
			newActivity := activities[rand.Intn(len(activities))]
			m.activities = append([]string{newActivity}, m.activities...)
			if len(m.activities) > 5 {
				m.activities = m.activities[:5]
			}
		}

		return m, tickCmd()
	}

	return m, nil
}

func (m model) View() string {
	screen := renderer.NewScreen(m.width, m.height)
	root := m.buildDashboard()

	constraints := layout.Tight(float64(m.width), float64(m.height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(m.width),
		ViewportHeight: float64(m.height),
		RootFontSize:   16,
	}
	layout.Layout(root.Node, constraints, ctx)
	screen.Render(root)

	return screen.String()
}

func (m model) buildDashboard() *renderer.StyledNode {
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(m.width)),
			Height:        layout.Px(float64(m.height)),
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
			Width:   layout.Px(float64(m.width - 4)),
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
	headerStyled.Content = fmt.Sprintf("System Dashboard - %s - Uptime: %s",
		time.Now().Format("15:04:05"), formatUptime(m.uptime))
	rootStyled.AddChild(headerStyled)

	// Resource gauges container
	gaugesContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(m.width - 4)),
			Height:        layout.Px(8),
			Margin:        layout.Spacing{Bottom: layout.Px(1)},
		},
	}
	gaugesStyled := renderer.NewStyledNode(gaugesContainer, nil)

	// Resource gauges
	gauges := []struct {
		label string
		value float64
		color string
	}{
		{"CPU", m.cpuUsage, "#FF6B6B"},
		{"Memory", m.memUsage, "#4ECDC4"},
		{"Disk", m.diskUsage, "#45B7D1"},
		{"Network", m.netUsage, "#96CEB4"},
	}

	gaugeWidth := (m.width - 4 - 6) / 4 // Divide by 4 gauges, account for spacing
	for i, g := range gauges {
		gauge := m.createGauge(g.label, g.value, gaugeWidth, g.color)
		gauge.Node.Style.Margin = layout.Spacing{Right: layout.Px(2)}
		if i == len(gauges)-1 {
			gauge.Node.Style.Margin.Right = layout.Px(0)
		}
		gaugesStyled.AddChild(gauge)
	}

	rootStyled.AddChild(gaugesStyled)

	// Bottom section: System Info + Activity
	bottomContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Width:         layout.Px(float64(m.width - 4)),
			FlexGrow:      1,
		},
	}
	bottomStyled := renderer.NewStyledNode(bottomContainer, nil)

	// System Info (left)
	sysInfo := m.createSystemInfo()
	sysInfo.Node.Style.Width = layout.Px(float64((m.width - 4) / 2))
	sysInfo.Node.Style.Margin.Right = layout.Px(1)
	bottomStyled.AddChild(sysInfo)

	// Activity Log (right)
	activityLog := m.createActivityLog()
	activityLog.Node.Style.Width = layout.Px(float64((m.width - 4) / 2))
	activityLog.Node.Style.Margin.Left = layout.Px(1)
	bottomStyled.AddChild(activityLog)

	rootStyled.AddChild(bottomStyled)

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(m.width - 4)),
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
	footerStyled.Content = "Press 'q' or Ctrl+C to quit"
	rootStyled.AddChild(footerStyled)

	return rootStyled
}

func (m model) createGauge(label string, value float64, width int, colorHex string) *renderer.StyledNode {
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

func (m model) createSystemInfo() *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display:   layout.DisplayBlock,
			FlexGrow:  1,
			MinHeight: layout.Px(10), // Prevent collapse to zero height
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

	var m2 runtime.MemStats
	runtime.ReadMemStats(&m2)

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
		float64(m2.Alloc)/1024/1024,
	)

	styledNode.Content = content

	return styledNode
}

func (m model) createActivityLog() *renderer.StyledNode {
	node := &layout.Node{
		Style: layout.Style{
			Display:   layout.DisplayBlock,
			FlexGrow:  1,
			MinHeight: layout.Px(10), // Prevent collapse to zero height
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

	content := "Recent Activity\n\n"
	for _, activity := range m.activities {
		timestamp := time.Now().Format("15:04:05")
		content += fmt.Sprintf("[%s] %s\n", timestamp, activity)
	}

	styledNode.Content = content

	return styledNode
}

func formatUptime(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running dashboard: %v\n", err)
		os.Exit(1)
	}
}
