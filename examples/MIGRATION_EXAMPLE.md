# Migration Example: Using New CSS Units

## Before (Old API with float64):

```go
root := &layout.Node{
    Style: layout.Style{
        Display:             layout.DisplayGrid,
        Width:               float64(m.width),  // Hardcoded pixels
        Height:              float64(m.height),
        GridTemplateColumns: []layout.GridTrack{
            layout.FractionTrack(1),
            layout.FractionTrack(1),
        },
        GridGap: 1,
    },
}

constraints := layout.Tight(float64(m.width), float64(m.height))
layout.Layout(root, constraints)
```

## After (New API with CSS units):

```go
// Create layout context with viewport dimensions
ctx := layout.NewLayoutContext(m.width, m.height, 16) // width, height, root font size

root := &layout.Node{
    Style: layout.Style{
        Display:             layout.DisplayGrid,
        Width:               layout.Vw(100),  // 100% of viewport width
        Height:              layout.Vh(100),  // 100% of viewport height
        GridTemplateColumns: []layout.GridTrack{
            layout.FractionTrack(1),
            layout.FractionTrack(1),
        },
        GridGap: layout.Rem(1),  // 1x root font size spacing
    },
}

// Header spanning full width with viewport-relative height
headerNode := &layout.Node{
    Style: layout.Style{
        Width:   layout.Vw(100),
        Height:  layout.Ch(3),  // 3 character heights
        Padding: layout.Uniform(layout.Rem(1)),
    },
}

// Panel with mixed units
panelNode := &layout.Node{
    Style: layout.Style{
        Width:    layout.Vw(50),      // 50% of viewport
        MinWidth: layout.Ch(30),      // Minimum 30 characters wide
        Padding:  layout.Uniform(layout.Em(1)),  // 1x element font size
        Margin: layout.Spacing{
            Top:  layout.Vh(2),       // 2% of viewport height
            Left: layout.Rem(0.5),    // Half root font size
        },
    },
}

constraints := layout.Tight(layout.Vw(100), layout.Vh(100))
layout.Layout(root, constraints, ctx)  // Pass context as third parameter
```

## Benefits:

1. **Semantic units**: `Vh(50)` is clearer than `float64(height/2)`
2. **Responsive by default**: Viewport units adapt automatically
3. **Text-aware**: `Ch()` and `Em()` units respect font metrics
4. **CSS familiarity**: Same units web developers know
5. **Mixed units**: Can combine `Vw()`, `Rem()`, `Ch()` in same layout

## Common Patterns:

### Full-screen layout
```go
Width: layout.Vw(100),
Height: layout.Vh(100),
```

### Text-based sizing
```go
Width: layout.Ch(80),  // 80 characters wide (terminal standard)
Height: layout.Ch(24), // 24 lines tall
```

### Responsive panels
```go
Width: layout.Vw(33),  // 1/3 of viewport
MinWidth: layout.Ch(40), // But never smaller than 40 chars
```

### Consistent spacing
```go
Padding: layout.Uniform(layout.Rem(1)),  // Scale with root font size
GridGap: layout.Rem(0.5),
```
