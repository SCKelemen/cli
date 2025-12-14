package components

import (
	"strings"
	"time"

	"github.com/SCKelemen/color"
	"github.com/SCKelemen/layout"
	"github.com/SCKelemen/cli/renderer"
)

// LoadingDots represents an animated loading indicator
type LoadingDots struct {
	Phase       int
	MaxPhases   int
	Interval    time.Duration
	LastUpdate  time.Time
	Foreground  *color.Color
}

// NewLoadingDots creates a new loading dots component
func NewLoadingDots() *LoadingDots {
	fg, _ := color.ParseColor("#7D56F4")
	return &LoadingDots{
		Phase:      0,
		MaxPhases:  4,
		Interval:   500 * time.Millisecond,
		LastUpdate: time.Now(),
		Foreground: &fg,
	}
}

// Update updates the animation state
func (l *LoadingDots) Update(now time.Time) bool {
	if now.Sub(l.LastUpdate) >= l.Interval {
		l.Phase = (l.Phase + 1) % l.MaxPhases
		l.LastUpdate = now
		return true // Needs redraw
	}
	return false
}

// ToStyledNode converts the loading dots to a styled node
func (l *LoadingDots) ToStyledNode() *renderer.StyledNode {
	// Generate the loading text based on phase
	var text string
	switch l.Phase {
	case 0:
		text = "Loading   "
	case 1:
		text = "Loading.  "
	case 2:
		text = "Loading.. "
	case 3:
		text = "Loading..."
	}

	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(len(text)),
			Height:  1,
		},
	}

	style := &renderer.Style{
		Foreground: l.Foreground,
	}

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = text

	return styledNode
}

// SpinnerDots represents a spinning dot animation
type SpinnerDots struct {
	Phase      int
	Interval   time.Duration
	LastUpdate time.Time
	Foreground *color.Color
	frames     []string
}

// NewSpinnerDots creates a new spinner dots component
func NewSpinnerDots() *SpinnerDots {
	fg, _ := color.ParseColor("#7D56F4")
	return &SpinnerDots{
		Phase:      0,
		Interval:   100 * time.Millisecond,
		LastUpdate: time.Now(),
		Foreground: &fg,
		frames:     []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	}
}

// Update updates the animation state
func (s *SpinnerDots) Update(now time.Time) bool {
	if now.Sub(s.LastUpdate) >= s.Interval {
		s.Phase = (s.Phase + 1) % len(s.frames)
		s.LastUpdate = now
		return true
	}
	return false
}

// ToStyledNode converts the spinner to a styled node
func (s *SpinnerDots) ToStyledNode() *renderer.StyledNode {
	text := s.frames[s.Phase]

	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   2, // Unicode spinners can be wider
			Height:  1,
		},
	}

	style := &renderer.Style{
		Foreground: s.Foreground,
	}

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = text

	return styledNode
}

// ProgressBar represents a progress bar component
type ProgressBar struct {
	Progress   float64 // 0.0 to 1.0
	Width      int
	Foreground *color.Color
	Background *color.Color
}

// NewProgressBar creates a new progress bar
func NewProgressBar(width int) *ProgressBar {
	fg, _ := color.ParseColor("#7D56F4")
	bg, _ := color.ParseColor("#3C3C3C")
	return &ProgressBar{
		Progress:   0.0,
		Width:      width,
		Foreground: &fg,
		Background: &bg,
	}
}

// SetProgress sets the progress value (0.0 to 1.0)
func (p *ProgressBar) SetProgress(progress float64) {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	p.Progress = progress
}

// ToStyledNode converts the progress bar to a styled node
func (p *ProgressBar) ToStyledNode() *renderer.StyledNode {
	filled := int(float64(p.Width) * p.Progress)
	empty := p.Width - filled

	text := strings.Repeat("█", filled) + strings.Repeat("░", empty)

	node := &layout.Node{
		Style: layout.Style{
			Display: layout.DisplayBlock,
			Width:   float64(p.Width),
			Height:  1,
		},
	}

	style := &renderer.Style{
		Foreground: p.Foreground,
	}

	styledNode := renderer.NewStyledNode(node, style)
	styledNode.Content = text

	return styledNode
}
