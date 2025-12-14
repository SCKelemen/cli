package renderer

import "github.com/SCKelemen/color"

// Style defines visual attributes without sizing properties.
// Sizing and layout are handled by the layout engine.
type Style struct {
	// Colors
	Foreground *color.Color
	Background *color.Color

	// Text attributes
	Bold          bool
	Italic        bool
	Underline     bool
	Strikethrough bool
	Dim           bool
	Blink         bool
	Reverse       bool

	// Borders
	Border      *BorderStyle
	BorderColor *color.Color
}

// BorderStyle defines which borders to render
type BorderStyle struct {
	Top    bool
	Right  bool
	Bottom bool
	Left   bool

	// Border characters
	Chars BorderChars
}

// BorderChars defines the characters used for borders
type BorderChars struct {
	TopLeft     rune
	TopRight    rune
	BottomLeft  rune
	BottomRight rune
	Horizontal  rune
	Vertical    rune
}

// Predefined border styles
var (
	RoundedBorder = BorderChars{
		TopLeft:     '╭',
		TopRight:    '╮',
		BottomLeft:  '╰',
		BottomRight: '╯',
		Horizontal:  '─',
		Vertical:    '│',
	}

	DoubleBorder = BorderChars{
		TopLeft:     '╔',
		TopRight:    '╗',
		BottomLeft:  '╚',
		BottomRight: '╝',
		Horizontal:  '═',
		Vertical:    '║',
	}

	ThickBorder = BorderChars{
		TopLeft:     '┏',
		TopRight:    '┓',
		BottomLeft:  '┗',
		BottomRight: '┛',
		Horizontal:  '━',
		Vertical:    '┃',
	}

	NormalBorder = BorderChars{
		TopLeft:     '┌',
		TopRight:    '┐',
		BottomLeft:  '└',
		BottomRight: '┘',
		Horizontal:  '─',
		Vertical:    '│',
	}
)

// NewStyle creates a new empty style
func NewStyle() *Style {
	return &Style{}
}

// WithForeground sets the foreground color
func (s *Style) WithForeground(c *color.Color) *Style {
	s.Foreground = c
	return s
}

// WithBackground sets the background color
func (s *Style) WithBackground(c *color.Color) *Style {
	s.Background = c
	return s
}

// WithBold sets bold text
func (s *Style) WithBold(bold bool) *Style {
	s.Bold = bold
	return s
}

// WithItalic sets italic text
func (s *Style) WithItalic(italic bool) *Style {
	s.Italic = italic
	return s
}

// WithUnderline sets underlined text
func (s *Style) WithUnderline(underline bool) *Style {
	s.Underline = underline
	return s
}

// WithBorder sets a border with all sides enabled
func (s *Style) WithBorder(chars BorderChars) *Style {
	s.Border = &BorderStyle{
		Top:    true,
		Right:  true,
		Bottom: true,
		Left:   true,
		Chars:  chars,
	}
	return s
}

// WithBorderColor sets the border color
func (s *Style) WithBorderColor(c *color.Color) *Style {
	s.BorderColor = c
	return s
}
