package main

type EditResult struct {
	Modified bool
	// TODO: DidScroll, Error, etc.
}

var (
	BufferUnchanged = EditResult{Modified: false}
	BufferModified  = EditResult{Modified: true}
)

// #region SelectionHelpers

// TODO: get hovered area
func WithActiveArea(f func(*Area)) {
	for _, area := range state.Areas {
		if area.IsFocused {
			f(area)
			return
		}
	}
}

func WithFocusedArea(f func(*Area)) {
	for _, area := range state.Areas {
		if area.IsFocused {
			f(area)
			return
		}
	}
}

func WithFocusedBuffer(f func(*Buffer) EditResult) {
	// for _, area := range state.Areas {
	// 	if area.IsFocused {
	// 		change := BufferUnchanged
	// 		if area.Context == nil {
	// 			LogErr("Focused area has no context")
	// 			continue
	// 		}
	// 		if area.Buffer != nil {
	// 			change = f(area.Buffer)
	// 		} else {
	// 			LogErr("Focused area has no buffer")
	// 		}
	// 		if change == BufferModified {
	// 			area.Buffer.MarkModified()
	// 		}
	// 	}
	// }
}

func WithFocusedCursors(f func(*Buffer, *Cursor) EditResult) {
	// for _, area := range state.Areas {
	// 	if area.IsFocused {
	// 		change := BufferUnchanged
	// 		if area.Buffer != nil {
	// 			for _, cursor := range area.Cursors {
	// 				local_change := f(area.Buffer, &cursor)
	// 				if local_change == BufferModified {
	// 					change = BufferModified
	// 				}
	// 			}
	// 		} else {
	// 			LogErr("Focused area has no buffer")
	// 		}
	// 		if change == BufferModified {
	// 			area.Buffer.MarkModified()
	// 		}
	// 	}
	// }
}

// #endregion SelectionHelpers

func EditInsertRune(r rune) {
	if state.Prompt.Visible {
		PromptInput(r)
		return
	}
	// WithBufferAndCursor(func(buf *Buffer, cur *Cursor) EditResult {
	// 	// line := &buf.Lines[cur.Position.Y]
	// 	// runes := line.ToRunes()

	// 	// // newRunes := append(runes[:cur.Position.X], append([]rune{r}, runes[cur.Position.X:]...)...)
	// 	// // line.SetRunes(newRunes)
	// 	// // cur.Position.X += 1

	// 	return BufferUnchanged
	// })
}

func EditDelRight() {
	// WithBufferAndCursor(func(buf *Buffer, cur *Cursor) EditResult {
	// 	if cur.Position.Y >= len(buf.Lines) {
	// 		return BufferUnchanged
	// 	}
	// 	line := &buf.Lines[cur.Position.Y]
	// 	runes := line.ToRunes()

	// 	if cur.Position.X >= len(runes) {
	// 		return BufferUnchanged // nothing to delete
	// 	}

	// 	newRunes := append(runes[:cur.Position.X], runes[cur.Position.X+1:]...)
	// 	line.SetRunes(newRunes)

	// 	return BufferModified
	// })
}

func EditDelLeft() {
	// WithBufferAndCursor(func(buf *Buffer, cur *Cursor) EditResult {
	// 	if cur.Position.Y >= len(buf.Lines) {
	// 		return BufferUnchanged
	// 	}
	// 	line := &buf.Lines[cur.Position.Y]
	// 	runes := line.ToRunes()

	// 	if cur.Position.X == 0 {
	// 		// if on the first line, do nothing
	// 		if cur.Position.Y == 0 {
	// 			return BufferUnchanged
	// 		}

	// 		// otherwise, merge with previous line
	// 		prevLine := &buf.Lines[cur.Position.Y-1]
	// 		prevLineLen := len(prevLine.ToRunes())

	// 		// append current line to previous line
	// 		prevLine.SetRunes(append(prevLine.ToRunes(), runes...))

	// 		// remove current line
	// 		buf.Lines = append(buf.Lines[:cur.Position.Y], buf.Lines[cur.Position.Y+1:]...)

	// 		// move cursor
	// 		cur.Position.Y -= 1
	// 		cur.Position.X = prevLineLen

	// 		return BufferModified
	// 	}

	// 	newRunes := append(runes[:cur.Position.X-1], runes[cur.Position.X:]...)
	// 	line.SetRunes(newRunes)
	// 	cur.Position.X -= 1

	// 	return BufferModified
	// })
}

func EditLineBreak() {
	// WithBufferAndCursor(func(buf *Buffer, cur *Cursor) EditResult {
	// 	if cur.Position.Y >= len(buf.Lines) {
	// 		return BufferUnchanged
	// 	}

	// 	line := &buf.Lines[cur.Position.Y]
	// 	runes := line.ToRunes()

	// 	// Split runes into two parts
	// 	left := runes[:cur.Position.X]
	// 	right := runes[cur.Position.X:]

	// 	// Replace current line with left part
	// 	line.SetRunes(left)

	// 	// Insert right part as new line below
	// 	newLine := LineFromRunes(right)
	// 	buf.Lines = append(buf.Lines[:cur.Position.Y+1],
	// 		append([]Line{newLine}, buf.Lines[cur.Position.Y+1:]...)...)

	// 	// Move cursor to start of new line
	// 	cur.Position.Y += 1
	// 	cur.Position.X = 0

	// 	return BufferModified
	// })
}

func EditBackspace() {
	Log("Delete")
}

func EditDelete() {
	Log("Delete")
}
