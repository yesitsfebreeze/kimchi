package main

type Encoding int

const (
	EncodingUTF8 Encoding = iota
	EncodingUTF16
)

type Char struct {
	Rune  rune
	Style Style
}

type Line struct {
	Index int
	Chars []Char
}

type BufferType int

const (
	BufferTypeStandard BufferType = iota
	BufferTypeMini
)

type BufferView struct {
	Size     Vec2
	Position Vec2
	Offset   Vec2
}

type Buffer struct {
	Type     BufferType
	Name     string
	Path     string
	AbsPath  string
	Encoding Encoding
	Lines    []Line
	Modified bool
	MaxLines int
}

func (buf *Buffer) MarkModified() {
	buf.Modified = true
}

func (buf *Buffer) IsEmpty() bool {
	return len(buf.Lines) == 0
}

func (l *Line) ToRunes() []rune {
	runes := make([]rune, len(l.Chars))
	for i, char := range l.Chars {
		runes[i] = char.Rune
	}
	return runes
}

func (l *Line) SetRuneAt(x int, runes []rune, style Style) {
	if x < 0 || x >= len(l.Chars) {
		LogErrF("Index out of bounds: %d for line with %d characters", x, len(l.Chars))
		return
	}
	if x >= len(runes) {
		LogErrF("Index out of bounds: %d for runes with %d characters", x, len(runes))
		return
	}

	l.Chars[x] = Char{Rune: runes[x], Style: style}
}

func (l *Line) SetRunes(runes []rune, style Style) {
	l.Chars = make([]Char, len(runes))
	for i, r := range runes {
		l.Chars[i] = Char{Rune: r, Style: style}
	}
}

func LineFromRunesDefaultStyle(runes []rune) Line {
	line := Line{}
	line.SetRunes(runes, state.Theme.Text)
	return line
}

func LineFromRunes(runes []rune, style Style) Line {
	line := Line{}
	line.SetRunes(runes, style)
	return line
}

// #region RUNES

func (buf *Buffer) IsValidRune(line, col int) bool {
	return buf.IsValidLine(line) && col >= 0 && col < len(buf.Lines[line].ToRunes())
}

func (buf *Buffer) InsertRune(line, col int, r rune) {
	if !buf.IsValidRune(line, col) {
		return // or panic
	}

	lineRunes := buf.Lines[line].ToRunes()
	if col < 0 || col > len(lineRunes) {
		col = len(lineRunes) // append at the end
	}

	lineRunes = append(lineRunes[:col], append([]rune{r}, lineRunes[col:]...)...)
	buf.Lines[line] = LineFromRunesDefaultStyle(lineRunes)
	buf.MarkModified()
}

func (buf *Buffer) ReplaceRune(line, col int, r rune) {
	if !buf.IsValidRune(line, col) {
		return // or panic
	}

	lineRunes := buf.Lines[line].ToRunes()
	if col < 0 || col >= len(lineRunes) {
		return // or panic
	}

	lineRunes[col] = r
	buf.Lines[line] = LineFromRunesDefaultStyle(lineRunes)
	buf.MarkModified()
}

func (buf *Buffer) DeleteRune(line, col int) {
	if !buf.IsValidRune(line, col) {
		return // or panic
	}

	lineRunes := buf.Lines[line].ToRunes()
	if col < 0 || col >= len(lineRunes) {
		return // or panic
	}

	lineRunes = append(lineRunes[:col], lineRunes[col+1:]...)
	buf.Lines[line] = LineFromRunesDefaultStyle(lineRunes)
	buf.MarkModified()
}

func (buf *Buffer) GetRune(line, col int) rune {
	if !buf.IsValidRune(line, col) {
		return ' '
	}
	return buf.Lines[line].ToRunes()[col]
}

// #endregion RUNES

// #region LINES

func (buf *Buffer) EnsureLine(index int) {
	if index < 0 {
		ThrowErrF("Invalid line index: %d", index)
	}
	if index >= len(buf.Lines) {
		for i := len(buf.Lines); i <= index; i++ {
			buf.Lines = append(buf.Lines, Line{Index: i, Chars: []Char{}})
		}
	}
}

func (buf *Buffer) IsValidLine(index int) bool {
	return index >= 0 && index < len(buf.Lines)
}

func (buf *Buffer) LineCount() int {
	return len(buf.Lines)
}

func (buf *Buffer) LineLength(line int) int {
	if line < 0 || line >= len(buf.Lines) {
		return 0
	}
	return len(buf.Lines[line].ToRunes())
}

// inclusive start, exclusive end
func (buf *Buffer) GetLines(start, end int) [][]rune {
	if start < 0 || end > len(buf.Lines) || start >= end {
		return nil // or panic
	}
	lines := make([][]rune, end-start)
	for i := start; i < end; i++ {
		lines[i-start] = buf.Lines[i].ToRunes()
	}
	return lines
}

func (buf *Buffer) GetLine(line int) []rune {
	if !buf.IsValidLine(line) {
		return nil
	}
	return buf.Lines[line].ToRunes()
}

func (buf *Buffer) SetLine(index int, r []rune) {
	if !buf.IsValidLine(index) {
		buf.EnsureLine(index)
	}
	buf.Lines[index] = LineFromRunesDefaultStyle(r)
	buf.MarkModified()
}

func (buf *Buffer) InsertLine(index int, r []rune) {
	if index < 0 || index > len(buf.Lines) {
		buf.EnsureLine(index)
	}
	line := LineFromRunesDefaultStyle(r)
	if index == len(buf.Lines) {
		buf.Lines = append(buf.Lines, line)
	} else {
		buf.Lines = append(buf.Lines[:index+1], buf.Lines[index:]...)
		buf.Lines[index] = line
	}

	if buf.MaxLines > 0 && len(buf.Lines) > buf.MaxLines {
		// Trim excess lines from the top
		excess := len(buf.Lines) - buf.MaxLines
		buf.Lines = buf.Lines[excess:]
	}

	buf.Modified = true
}

func (buf *Buffer) ReplaceLines(start, end int, lines [][]rune) {
	if start < 0 || end > len(buf.Lines) || start >= end {
		return // or panic
	}

	if len(lines) != end-start {
		return // or panic, or handle resizing
	}

	for i := start; i < end; i++ {
		buf.Lines[i] = LineFromRunesDefaultStyle(lines[i-start])
	}

	buf.MarkModified()
}

func (buf *Buffer) DeleteLine(index int) {
	if index < 0 || index >= len(buf.Lines) {
		return // or panic
	}
	buf.Lines = append(buf.Lines[:index], buf.Lines[index+1:]...)
	buf.MarkModified()
}

func (buf *Buffer) AppendLine(r []rune) {
	buf.SetLine(len(buf.Lines), r)
}

func (buf *Buffer) PrependLine(r []rune) {
	buf.SetLine(0, r)
}

// #endregion LINES
