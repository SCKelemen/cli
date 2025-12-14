# CLI Layout Engine

A proof of concept for integrating a CSS-based layout engine with terminal UIs in Go. This project demonstrates how to build beautifully laid out terminal interfaces using CSS Grid, Flexbox, and the box model.

## Features

- **CSS-based Layouts**: Full support for Flexbox and Grid layouts in the terminal
- **Advanced Color Support**: OKLCH color gradients and wide-gamut color spaces via the `color` library
- **Responsive Design**: Automatically relayouts on terminal resize
- **Component System**: Reusable UI components (collapsible sections, loading indicators, progress bars)
- **Smart Rendering**: Efficient screen buffer with ANSI escape code optimization
- **Animation Support**: 30fps animation timeline for loading indicators and transitions

## Architecture

### Core Components

- **renderer/**: Core rendering system
  - `style.go` - Visual styling (colors, borders, text attributes)
  - `ansi.go` - ANSI escape code generation
  - `screen.go` - Screen buffer and rendering logic

- **components/**: Reusable UI components
  - `message.go` - Styled message blocks
  - `loading.go` - Loading indicators, spinners, progress bars
  - `collapsible.go` - Expandable/collapsible sections

- **examples/**: Demo applications
  - `demo.go` - Claude Code-like UI demonstration

### Design Principles

1. **Separation of Concerns**
   - Layout engine handles sizing, positioning, and spacing
   - Style system handles colors, borders, and visual effects
   - No property duplication between systems

2. **Professional Color Handling**
   - Perceptually uniform color operations using OKLCH
   - Support for CSS color parsing
   - Wide-gamut color space support

3. **Responsive by Default**
   - Layouts automatically adapt to terminal size
   - Bubbletea integration for resize events

## Dependencies

- [layout](https://github.com/SCKelemen/layout) - CSS Grid/Flexbox layout engine
- [color](https://github.com/SCKelemen/color) - Professional color manipulation
- [bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework

## Running the Demo

```bash
cd examples
go build -o demo demo.go
./demo
```

### Demo Controls

- **'1'** - Toggle first collapsible section
- **'2'** - Toggle second collapsible section
- **Space** - Toggle all sections
- **'q'** or **ESC** or **Ctrl+C** - Quit

## Example Usage

### Creating a Simple Layout

```go
package main

import (
    "github.com/SCKelemen/cli/renderer"
    "github.com/SCKelemen/cli/components"
    "github.com/SCKelemen/layout"
    "github.com/SCKelemen/color"
)

func main() {
    // Create a message block
    msg := components.NewMessageBlock("Hello, World!")

    // Create layout tree
    root := &layout.Node{
        Style: layout.Style{
            Display: layout.DisplayFlex,
            FlexDirection: layout.FlexDirectionColumn,
            Width: 80,
            Height: 24,
            Padding: layout.Uniform(2),
        },
    }

    rootStyled := renderer.NewStyledNode(root, nil)
    rootStyled.AddChild(msg.ToStyledNode())

    // Compute layout
    layout.Layout(root, layout.Tight(80, 24))

    // Render to screen
    screen := renderer.NewScreen(80, 24)
    screen.Render(rootStyled)
    print(screen.String())
}
```

### Creating Custom Components

```go
type CustomComponent struct {
    Title string
    Color *color.Color
}

func (c *CustomComponent) ToStyledNode() *renderer.StyledNode {
    node := &layout.Node{
        Style: layout.Style{
            Display: layout.DisplayBlock,
            Width: 40,
            Height: 3,
        },
    }

    style := &renderer.Style{
        Foreground: c.Color,
        Bold: true,
    }
    style.WithBorder(renderer.RoundedBorder)

    styledNode := renderer.NewStyledNode(node, style)
    styledNode.Content = c.Title

    return styledNode
}
```

## Future Enhancements

- Line-based CLI mode (print-only, non-redrawable content)
- Multiline, multicursor text editor component
- More advanced animation easing functions
- Text-based units (em, rem, ch) in layout engine
- Viewport units (vh, vw) in layout engine
- Bottom-up content flow and reflow
- More component types (tables, lists, menus)

## Contributing

This is a proof of concept project. Feel free to experiment and extend it for your own use cases.

## License

MIT
