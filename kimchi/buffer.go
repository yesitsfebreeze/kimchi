package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Buffer struct {
	name     string
	path     string
	modified bool
	cursor_x int
	cursor_y int
	content  [][]rune
	cursors  Cursors
}

// load_buffer reads a file into a buffer.
func load_buffer(path string) (*Buffer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Normalize line endings for cross-platform sanity
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(text, "\n")

	content := make([][]rune, len(lines))
	for i, line := range lines {
		content[i] = []rune(line)
	}

	// Ensure at least one line exists
	if len(content) == 0 {
		content = append(content, []rune{})
	}

	name := filepath.Base(path)

	return &Buffer{
		name:    name,
		path:    path,
		content: content,
	}, nil
}

// new_empty_buffer creates a new in-memory buffer (not yet tied to a file).
func new_empty_buffer(name string) *Buffer {
	return &Buffer{
		name:    name,
		content: [][]rune{{}}, // start with one empty line
	}
}

// save_buffer writes the buffer contents back to its file.
func (b *Buffer) save_buffer() error {
	var sb strings.Builder
	for i, line := range b.content {
		sb.WriteString(string(line))
		if i < len(b.content)-1 {
			sb.WriteByte('\n')
		}
	}

	err := os.WriteFile(b.path, []byte(sb.String()), 0644)
	if err != nil {
		return err
	}

	b.modified = false
	return nil
}
