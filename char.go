package gofiglet

import (
	"errors"
)

// asciiChar represents a single rendered ASCII-art character: the set of
// text lines that make up its glyph, plus an optional color applied when
// the line is emitted.
type asciiChar struct {
	lines []string
	color Color
}

// newASCIIChar builds an asciiChar for char using font's glyph data.
// It returns an error if char falls outside the printable ASCII range
// (0-127), since font only defines glyphs for ASCII characters.
func newASCIIChar(font *font, char rune) (*asciiChar, error) {
	if char < 0 || char > 127 {
		return nil, errors.New("not Ascii character")
	}
	lines, err := font.getCharSlice(char)
	if err != nil {
		return nil, err
	}

	return &asciiChar{lines: lines}, nil
}

// GetLine returns the line at index, wrapped in the character's color
// escape sequences if color is set, or unwrapped otherwise.
func (char *asciiChar) GetLine(index int) string {
	prefix := ""
	suffix := ""

	line := char.lines[index]

	if char.color != nil {
		prefix = char.color.GetPrefix()
		suffix = char.color.GetSuffix()
	}

	return prefix + line + suffix
}
