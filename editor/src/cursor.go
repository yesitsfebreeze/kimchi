package main

type CursorTrailCell struct {
	Position Vec2
	Time     int64
}

type Cursor struct {
	Position Vec2
	LastX    int
	Area     *Area
	Trails   []CursorTrailCell
}

func (cur *Cursor) GetArea() (*Area, bool) {
	if cur.Area == nil {
		LogErr("Cursor has no Area")
		return nil, false
	}
	return cur.Area, true
}

func (cur *Cursor) Move(dx, dy int) {
	// area, ok := cur.GetArea()
	// if !ok {
	// 	LogErr("Cursor has no parent buffer")
	// 	return
	// }

	// lines := buf.Lines
	// if len(lines) == 0 {
	// 	cur.Position = Vec2{}
	// 	return
	// }

	// // Handle vertical movement first
	// if dy != 0 {
	// 	cur.Position.Y += dy
	// 	cur.Clamp()
	// 	// After moving vertically, try to restore LastX
	// 	lineLen := len(buf.Lines[cur.Position.Y].ToRunes())
	// 	cur.Position.X = min(cur.LastX, lineLen)
	// }

	// // Handle horizontal movement
	// if dx != 0 {
	// 	newX := cur.Position.X + dx
	// 	lineLen := len(buf.Lines[cur.Position.Y].ToRunes())
	// 	if newX < 0 {
	// 		if cur.Position.Y > 0 {
	// 			cur.Position.Y -= 1
	// 			cur.Position.X = len(buf.Lines[cur.Position.Y].ToRunes())
	// 		}
	// 	} else if newX > lineLen {
	// 		if cur.Position.Y < len(lines)-1 {
	// 			cur.Position.Y += 1
	// 			cur.Position.X = 0
	// 		}
	// 	} else {
	// 		cur.Position.X = newX
	// 	}
	// 	cur.LastX = cur.Position.X
	// }

	cur.Clamp()
	cur.AddTrail()
}

func (cur *Cursor) Clamp() {
	// buf, ok := cur.GetBuffer()
	// if !ok {
	// 	LogErr("Cursor has no parent buffer")
	// 	return
	// }
	// if cur.Position.Y >= len(buf.Lines) {
	// 	cur.Position.Y = len(buf.Lines) - 1
	// }
	// if cur.Position.Y < 0 {
	// 	cur.Position.Y = 0
	// }
	// line := buf.Lines[cur.Position.Y].ToRunes()
	// if cur.Position.X > len(line) {
	// 	cur.Position.X = len(line)
	// }
	// if cur.Position.X < 0 {
	// 	cur.Position.X = 0
	// }
}

func (cur *Cursor) AddTrail() {
	if !state.Config.CursorTrail.Enabled {
		return
	}

	if len(cur.Trails) == 0 {
		cur.Trails = append(cur.Trails, CursorTrailCell{
			Position: cur.Position,
			Time:     int64(state.Config.CursorTrail.Time),
		})
		return
	}

	last := cur.Trails[len(cur.Trails)-1].Position
	delta := cur.Position.Sub(last)

	// Only track orthogonal steps (one axis moves, not both)
	if (delta.X == 0 && delta.Y != 0) || (delta.Y == 0 && delta.X != 0) {
		cur.Trails = append(cur.Trails, CursorTrailCell{
			Position: cur.Position,
			Time:     int64(state.Config.CursorTrail.Time),
		})
	}
}

func (cur *Cursor) Draw() {
	area, ok := cur.GetArea()
	if !ok {
		LogErr("Cursor has no parent buffer")
		return
	}

	cur.DrawTrail()

	screen := cur.Position
	screen = screen.Add(area.Position)
	screen = screen.Sub(area.ScrollOffset)

	style := state.Theme.Cursor.Multi
	// if i == len(buf.Cursors)-1 {
	// 	style = kit.Theme.Cursor.Main
	// }

	if screen.Y >= 0 && screen.Y < len(state.Screen.Cells) &&
		screen.X >= 0 && screen.X < len(state.Screen.Cells[screen.Y]) {
		cell := state.Screen.Cells[screen.Y][screen.X]
		state.Screen.Cells[screen.Y][screen.X] = Cell{
			Position: screen,
			Rune:     cell.Rune,
			Style:    style,
		}
	}

}

func (cur *Cursor) DrawTrail() {
	// if !kit.Config.CursorTrail.Enabled {
	// 	return
	// }

	// buf, ok := cur.GetBuffer()
	// if !ok {
	// 	return
	// }

	// // length := kit.Config.CursorTrail.Length
	// for i := 1; i < len(cur.Trails); {
	// 	trail := &cur.Trails[i]
	// 	screen := trail.Position.Add(buf.View.Position).Sub(buf.View.Offset)

	// 	if screen.Y < 0 || screen.Y >= len(kit.Screen.Cells) ||
	// 		screen.X < 0 || screen.X >= len(kit.Screen.Cells[0]) {
	// 		// skip drawing if out of bounds
	// 		i++
	// 		continue
	// 	}

	// 	cell := kit.Screen.Cells[screen.Y][screen.X]
	// 	// fade := (float64(length-i) / float64(length))

	// 	style := kit.Theme.Cursor.Multi
	// 	// style.BG = BlendColor(style.BG, kit.Theme.Text.BG, fade)

	// 	kit.Screen.Cells[screen.Y][screen.X] = Cell{
	// 		Position: screen,
	// 		Rune:     cell.Rune,
	// 		Style:    style,
	// 	}

	// 	trail.Time -= kit.Delta
	// 	if trail.Time <= 0 {
	// 		// remove this trail
	// 		cur.Trails = append(cur.Trails[:i], cur.Trails[i+1:]...)
	// 		// don't increment i
	// 	} else {
	// 		// clamp and continue
	// 		trail.Time = max(0, trail.Time)
	// 		i++
	// 	}
	// }
}

func (buf *Buffer) PrimaryCursor() *Cursor {
	// if len(buf.Cursors) == 0 {
	// 	// LogErr("No cursors in buffer")
	// 	return nil
	// }
	// return &buf.Cursors[0]

	return nil
}

func CursorMove(dx, dy int) {
	// WithBuffer(func(buf *Buffer) {
	// 	WithCursor(buf, func(cur *Cursor) EditResult {
	// 		cur.Move(dx, dy)
	// 		return BufferUnchanged
	// 	})
	// })
}

func CursorMoveStart() {

}

func CursorMoveEnd() {

}
