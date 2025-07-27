package main

const MIN_ZINDEX int = 1
const MAX_ZINDEX int = 5000

const BASE_ZINDEX int = 1000
const MID_ZINDEX int = 2000
const FRONT_ZINDEX int = 3000

type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderRounded
	BorderSquare
	BorderBold
	BorderDouble
	BorderCustom
)

var BorderStyles = map[BorderStyle]string{
	BorderNone:    "       ",
	BorderRounded: "│╭─╮╰─╯",
	BorderSquare:  "│┌─┐└─┘",
	BorderBold:    "┃┏━┓┗━┛",
	BorderDouble:  "║╔═╗╚═╝",
	BorderCustom:  "custom",
}

type BorderUsage struct {
	Top    bool
	Bottom bool
	Left   bool
	Right  bool
}

type Padding struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

func PaddingXY(x, y int) Padding {
	return Padding{
		Top:    y,
		Right:  x,
		Bottom: y,
		Left:   x,
	}
}

func (p Padding) Start() Vec2 {
	return Vec2{
		X: p.Left,
		Y: p.Top,
	}
}

type Area struct {
	Name         string
	Hidden       bool
	ClampIgnored bool
	Size         Vec2
	Position     Vec2
	ScrollOffset Vec2
	IsFocused    bool
	ZIndex       int
	BorderStyle  BorderStyle
	BorderUsage  BorderUsage
	Padding      Padding
	Context      *Context
	Style        Style
}

func CreateArea(name string) *Area {
	a := &Area{
		Name: name,
	}
	state.Areas = append(state.Areas, a)
	state.NamedAreas[name] = a

	a.BorderStyle = BorderNone

	return a
}

func GetAreaByName(name string) *Area {
	if area, ok := state.NamedAreas[name]; ok {
		return area
	}
	LogErr("Area not found:", name)
	return nil
}

func DeleteAreaByName(name string) {
	for i := range state.Areas {
		if state.Areas[i].Name == name {
			state.Areas = append(state.Areas[:i], state.Areas[i+1:]...)
		}
	}
}

func (a *Area) IgnoreClamp(flag bool) *Area {
	a.ClampIgnored = flag
	return a
}

func (a *Area) SetHidden(hidden bool) *Area {
	a.Hidden = hidden
	if hidden {
		a.IsFocused = false
	}
	return a
}

func (a *Area) ClampSize() *Area {
	minY := 0
	maxY := state.Screen.Size.Y
	maxHeight := state.Screen.Size.Y

	statusbarHeight := 0
	if state.Config.Statusbar.Enabled {
		statusbarHeight = 1
		if state.Config.Statusbar.UseSeparator {
			statusbarHeight += 1
		}
	}

	maxHeight -= statusbarHeight

	if a.Size.Y > maxHeight {
		a.Size.Y = maxHeight
	}

	if state.Config.Statusbar.Position == StatusbarTop {
		minY += statusbarHeight
		if a.Position.Y < minY {
			a.Position.Y = minY
		}
	}

	if state.Config.Statusbar.Position == StatusbarBottom {
		maxY -= statusbarHeight
		if a.Position.Y+a.Size.Y > maxY {
			a.Position.Y -= statusbarHeight
		}
	}

	return a
}

func (a *Area) SetSize(size Vec2) *Area {
	a.Size = size
	if a.ClampIgnored {
		return a
	}
	a.ClampSize()
	return a
}

func (a *Area) SetPosition(pos Vec2) *Area {
	a.Position = pos
	if a.ClampIgnored {
		return a
	}
	a.ClampSize()
	return a
}

func (a *Area) SetContext(ctx *Context) *Area {
	a.Context = ctx
	return a
}

func (a *Area) SetBorderStyle(style BorderStyle) *Area {
	a.BorderStyle = style
	return a
}

func (a *Area) SetBorderUsage(usage BorderUsage) *Area {
	a.BorderUsage = usage
	return a
}

func (a *Area) SetPadding(padding Padding) *Area {
	a.Padding = padding
	return a
}

func (a *Area) SetZIndex(z int) *Area {
	if state.ZIndexedAreas == nil {
		state.ZIndexedAreas = make(map[int][]*Area)
	}

	if z < MIN_ZINDEX {
		z = MIN_ZINDEX
	}
	if z > MAX_ZINDEX {
		z = MAX_ZINDEX
	}

	if a.ZIndex != z && a.ZIndex != 0 {
		slice := state.ZIndexedAreas[a.ZIndex]
		for i, area := range slice {
			if area == a {
				state.ZIndexedAreas[a.ZIndex] = append(slice[:i], slice[i+1:]...)
				break
			}
		}
	}

	state.ZIndexedAreas[z] = append(state.ZIndexedAreas[z], a)
	a.ZIndex = z

	return a
}

func (a *Area) SetStyle(style Style) *Area {
	a.Style = style
	return a
}

func (a *Area) Focus() {
	for _, area := range state.Areas {
		if area != a {
			area.IsFocused = false
		}
	}
	a.IsFocused = true
	state.FocusedArea = a
}

func (a *Area) GetCharAt(coord Vec2) *Char {
	if a.Context == nil {
		return &Char{Rune: ' ', Style: a.Style}
	}
	buf := a.Context.Buffer

	if buf == nil || len(buf.Lines) == 0 {
		return &Char{Rune: ' ', Style: a.Style}
	}

	if coord.Y < 0 || coord.Y >= len(buf.Lines) {
		return &Char{Rune: ' ', Style: a.Style}
	}

	line := buf.Lines[coord.Y]
	if coord.X < 0 || coord.X >= len(line.Chars) {
		return &Char{Rune: ' ', Style: a.Style}
	}

	char := &line.Chars[coord.X]

	// if transparent use area bg
	if char.Style.BG.TrueColor().Hex() == -1 {
		char.Style.BG = a.Style.BG
	}

	return char
}

func (a *Area) GetCharAtScrollOffset(coord Vec2) *Char {
	return a.GetCharAt(Vec2{X: a.ScrollOffset.X + coord.X, Y: a.ScrollOffset.Y + coord.Y})
}

func (a *Area) Draw() {
	if a.Hidden {
		return
	}

	var darken float32 = 0
	if a == state.Panes.One.Area || a == state.Panes.Two.Area {
		if !a.IsFocused {
			darken = state.Config.UnfocusedDarken
		}
	}

	scr := state.Screen
	s := Vec2{X: a.Size.X, Y: a.Size.Y}
	p := Vec2{X: a.Position.X, Y: a.Position.Y}
	pad := Padding{
		Top:    a.Padding.Top,
		Bottom: a.Padding.Bottom,
		Left:   a.Padding.Left,
		Right:  a.Padding.Right,
	}

	// If border is enabled, content is inset by 1 character
	if a.BorderStyle != BorderNone {
		pad.Top += 1
		pad.Bottom += 1
		pad.Left += 1
		pad.Right += 1
	}

	s.X = max(s.X, 1+pad.Left+pad.Right)
	s.Y = max(s.Y, 1+pad.Top+pad.Bottom)

	// Draw content area (inset by border if present)
	contentWidth := s.X - pad.Left - pad.Right
	contentHeight := s.Y - pad.Top - pad.Bottom

	for y := 0; y < contentHeight; y++ {
		for x := 0; x < contentWidth; x++ {
			charCoord := Vec2{X: x, Y: y}
			char := a.GetCharAtScrollOffset(charCoord)

			// Calculate screen coordinates with border offset
			coord := Vec2.Add(p, Vec2{X: x, Y: y})
			coord = coord.Add(Vec2{X: pad.Left, Y: pad.Top})
			scr.Cell(coord, ' ', a.Style.Darken(darken))
			if char != nil {
				scr.Cell(coord, char.Rune, char.Style.Darken(darken))
			}
		}
	}

	// Draw border if enabled
	if a.BorderStyle != BorderNone {
		border := BorderStyles[a.BorderStyle]
		if len(border) < 7 {
			LogErr("Border style must have 7 characters")
			return
		}
		b := []rune(border)

		w, h := s.X, s.Y
		x, y := p.X, p.Y

		// Top and bottom borders
		for i := 1; i < w-1; i++ {
			scr.CellXY(x+i, y, b[2], a.Style)     // Top border
			scr.CellXY(x+i, y+h-1, b[2], a.Style) // Bottom border
		}

		// Left and right borders
		for j := 1; j < h-1; j++ {
			scr.CellXY(x, y+j, b[0], a.Style)     // Left border
			scr.CellXY(x+w-1, y+j, b[0], a.Style) // Right border
		}

		// Corner characters
		scr.CellXY(x, y, b[1], a.Style)         // Top-left corner
		scr.CellXY(x+w-1, y, b[3], a.Style)     // Top-right corner
		scr.CellXY(x, y+h-1, b[4], a.Style)     // Bottom-left corner
		scr.CellXY(x+w-1, y+h-1, b[6], a.Style) // Bottom-right corner
	}
}

// func (area *Area) AddCursor(x, y int) {
// 	// cur := Cursor{
// 	// 	Position: Vec2{X: x, Y: y},
// 	// 	Parent:   buf,
// 	// }
// 	// buf.Cursors = append(buf.Cursors, cur)
// }

// func (a *Area) UpdateViewSize() {
// 	// TODO: clampto screen size
// 	// a.Size = state.Screen.Size

// 	// // remove 1 for statusbar
// 	// a.Size.Y -= 1

// 	// a.Position.Y = 0
// 	// if state.Config.Statusbar.Position == StatusbarTop {
// 	// 	a.Position.Y += 1
// 	// }
// }

// func (buf *Buffer) ScrollToCursor(cur *Cursor) {
// 	// surrounding := kit.Config.SurroundingLines
// 	// topPad := min(surrounding, a.Size.Y/2)
// 	// bottomPad := min(surrounding, a.Size.Y/2)

// 	// // Calculate the visible window
// 	// visibleTop := a.ScrollOffset.Y + topPad
// 	// visibleBottom := a.ScrollOffset.Y + a.Size.Y - 1 - bottomPad

// 	// // newX := a.ScrollOffset.X
// 	// newY := a.ScrollOffset.Y

// 	// maxY := max(0, len(buf.Lines)-a.Size.Y)

// 	// if cur.Position.Y < visibleTop {
// 	// 	// Cursor is above the visible window, scroll up
// 	// 	newY = Clamp(cur.Position.Y-topPad, 0, maxY)
// 	// } else if cur.Position.Y > visibleBottom {
// 	// 	// Cursor is below the visible window, scroll down
// 	// 	newY = Clamp(cur.Position.Y-a.Size.Y+1+bottomPad, 0, maxY)
// 	// }
// 	// if newY != a.ScrollOffset.Y {
// 	// 	buf.ScrollTo(a.ScrollOffset.X, newY)
// 	// }
// }

func (a *Area) ScrollTo(x int, y int) {
	a.ScrollOffset.X = x
	a.ScrollOffset.Y = y
	a.ScrollOffset.X = Clamp(a.ScrollOffset.X, 0, a.Size.X)
	a.ScrollOffset.Y = Clamp(a.ScrollOffset.Y, 0, a.Size.Y)

	// clamp IF we have a buffer
	// for y := 0; y < a.Size.Y && (y+a.ScrollOffset.Y) < len(buf.Lines); y++ {
	// 	screenY := y + buf.View.Position.Y
	// 	kit.Screen.ClearLine(screenY)
	// }
}

func (a *Area) Scroll(dx int, dy int) {
	a.ScrollTo(a.ScrollOffset.X+dx, a.ScrollOffset.Y+dy)
}
