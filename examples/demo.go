package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/components"
	"github.com/SCKelemen/cli/renderer"
)

// tickMsg is sent on every animation frame
type tickMsg time.Time

// Model represents the application state
type Model struct {
	width      int
	height     int
	screen     *renderer.Screen
	loading    *components.LoadingDots
	spinner    *components.SpinnerDots
	progress   *components.ProgressBar
	section1   *components.Collapsible
	section2   *components.Collapsible
	quitting   bool
}

// NewModel creates a new model
func NewModel() Model {
	return Model{
		width:    80,
		height:   24,
		screen:   renderer.NewScreen(80, 24),
		loading:  components.NewLoadingDots(),
		spinner:  components.NewSpinnerDots(),
		progress: components.NewProgressBar(40),
		section1: components.NewCollapsible(
			"Project Information",
			"This is a proof of concept for integrating\nthe layout engine with terminal UIs.\n\nWe're building beautiful terminal interfaces!",
		),
		section2: components.NewCollapsible(
			"Features",
			"- Responsive layouts using CSS Grid/Flexbox\n- OKLCH color gradients\n- Collapsible sections\n- Animated loading indicators\n- Dynamic resizing",
		),
		quitting: false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
	)
}

// tickCmd returns a command that sends tick messages
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Handle terminal resize
		m.width = msg.Width
		m.height = msg.Height
		m.screen.Resize(m.width, m.height)
		return m, nil

	case tickMsg:
		// Update animations
		now := time.Time(msg)
		needsRedraw := false

		if m.loading.Update(now) {
			needsRedraw = true
		}
		if m.spinner.Update(now) {
			needsRedraw = true
		}

		// Update progress bar (demo: oscillate between 0 and 1)
		progress := float64(now.UnixNano()%2000000000) / 2000000000.0
		m.progress.SetProgress(progress)
		needsRedraw = true

		if needsRedraw || m.quitting {
			return m, tickCmd()
		}
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit

		case "1":
			m.section1.Toggle()
			return m, nil

		case "2":
			m.section2.Toggle()
			return m, nil

		case " ":
			// Toggle both sections
			m.section1.Toggle()
			m.section2.Toggle()
			return m, nil
		}
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	// Build the layout tree
	root := m.buildLayout()

	// Perform layout calculation
	constraints := layout.Tight(float64(m.width), float64(m.height))
	ctx := &layout.LayoutContext{
		ViewportWidth:  float64(m.width),
		ViewportHeight: float64(m.height),
		RootFontSize:   16,
	}
	layout.Layout(root.Node, constraints, ctx)

	// Render to screen
	m.screen.Render(root)

	return m.screen.String()
}

// buildLayout constructs the UI layout tree
func (m Model) buildLayout() *renderer.StyledNode {
	// Create root container with flexbox column layout
	root := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionColumn,
			Width:         layout.Px(float64(m.width)),
			Height:        layout.Px(float64(m.height)),
			Padding:       layout.Spacing{Top: layout.Px(2), Right: layout.Px(2), Bottom: layout.Px(2), Left: layout.Px(2)},
		},
	}
	rootStyled := renderer.NewStyledNode(root, nil)

	// Header
	headerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(m.width - 4)),
			Height:  layout.Px(5),
			Margin:  layout.Spacing{Top: layout.Px(0), Right: layout.Px(0), Bottom: layout.Px(1), Left: layout.Px(0)},
		},
	}
	fgWhite, _ := color.ParseColor("#FFFFFF")
	bgPurple, _ := color.ParseColor("oklch(0.5 0.2 270)")
	borderPurple, _ := color.ParseColor("oklch(0.7 0.2 270)")
	headerStyle := &renderer.Style{
		Foreground:  &fgWhite,
		Background:  &bgPurple,
		Bold:        true,
		BorderColor: &borderPurple,
	}
	headerStyle.WithBorder(renderer.ThickBorder)
	headerStyled := renderer.NewStyledNode(headerNode, headerStyle)
	headerStyled.Content = "  Terminal Layout Engine Demo  "
	rootStyled.AddChild(headerStyled)

	// Info message
	borderGray, _ := color.ParseColor("#5A5A5A")
	info := components.NewMessageBlock(
		"Press '1' or '2' to toggle sections\nPress 'q' to quit",
	).WithBorderColor(&borderGray)
	rootStyled.AddChild(info.ToStyledNode())

	// Loading indicators container
	loadingContainer := &layout.Node{
		Style: layout.Style{
			Display:       layout.DisplayFlex,
			FlexDirection: layout.FlexDirectionRow,
			Margin:        layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(1), Left: layout.Px(0)},
		},
	}
	loadingStyled := renderer.NewStyledNode(loadingContainer, nil)

	// Add spinner
	spinnerNode := m.spinner.ToStyledNode()
	spinnerNode.Node.Style.Margin = layout.Spacing{Top: layout.Px(0), Right: layout.Px(2), Bottom: layout.Px(0), Left: layout.Px(0)}
	loadingStyled.AddChild(spinnerNode)

	// Add loading dots
	loadingNode := m.loading.ToStyledNode()
	loadingNode.Node.Style.Margin = layout.Spacing{Top: layout.Px(0), Right: layout.Px(2), Bottom: layout.Px(0), Left: layout.Px(0)}
	loadingStyled.AddChild(loadingNode)

	rootStyled.AddChild(loadingStyled)

	// Progress bar
	progressNode := m.progress.ToStyledNode()
	progressNode.Node.Style.Margin = layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(1), Left: layout.Px(0)}
	rootStyled.AddChild(progressNode)

	// Collapsible sections
	section1Node := m.section1.ToStyledNode()
	section1Node.Node.Style.Margin = layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)}
	rootStyled.AddChild(section1Node)

	section2Node := m.section2.ToStyledNode()
	section2Node.Node.Style.Margin = layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)}
	rootStyled.AddChild(section2Node)

	// Footer
	footerNode := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   layout.Px(float64(m.width - 4)),
			Height:  layout.Px(3),
			Margin:  layout.Spacing{Top: layout.Px(1), Right: layout.Px(0), Bottom: layout.Px(0), Left: layout.Px(0)},
		},
	}
	fgGray, _ := color.ParseColor("#888888")
	footerStyle := &renderer.Style{
		Foreground: &fgGray,
		Dim:        true,
	}
	footerStyled := renderer.NewStyledNode(footerNode, footerStyle)
	footerStyled.Content = fmt.Sprintf("Terminal: %dx%d", m.width, m.height)
	rootStyled.AddChild(footerStyled)

	return rootStyled
}

func main() {
	p := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
