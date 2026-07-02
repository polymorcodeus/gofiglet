package gofiglet

// Explanation of the .flf file header
// THE HEADER LINE
//
// The header line gives information about the FIGfont.  Here is an example
// showing the names of all parameters:
//
//           flf2a$ 6 5 20 15 3 0 143 229    NOTE: The first five characters in
//             |  | | | |  |  | |  |   |     the entire file must be "flf2a".
//            /  /  | | |  |  | |  |   \
//   Signature  /  /  | |  |  | |   \   Codetag_Count
//     Hardblank  /  /  |  |  |  \   Full_Layout*
//          Height  /   |  |   \  Print_Direction
//          Baseline   /    \   Comment_Lines
//           Max_Length      Old_Layout*
//
//   * The two layout parameters are closely related and fairly complex.
//       (See "INTERPRETATION OF LAYOUT PARAMETERS".)
//

import (
	"fmt"
	"strings"
)

// font represents a single parsed FLF (figlet) font: its hardblank
// replacement character, glyph line height, and the raw slice of glyph
// data lines read from the font file (after the header and comment
// block).
type font struct {
	hardblank string
	height    int
	fontSlice []string
}

// getCharSlice returns the rendered lines that make up char's glyph in
// this font, with the font's hardblank character replaced by a literal
// space and "@" end-of-line markers stripped. char must be a printable
// ASCII rune (32-127); the caller is responsible for validating this,
// as getCharSlice does no bounds checking.
func (f *font) getCharSlice(char rune) ([]string, error) {
	if char < 32 || char > 126 {
		return nil, fmt.Errorf("unsupported character %q (code %d); printable ASCII 32-126 only", char, char)
	}

	height := f.height
	beginRow := (int(char) - 32) * height

	// Defensive: should never trigger if parseFontContent validated correctly
	if beginRow+height > len(f.fontSlice) {
		return nil, fmt.Errorf("font data truncated for character %q", char)
	}

	lines := make([]string, height)
	for i := range height {
		row := f.fontSlice[beginRow+i]
		row = strings.ReplaceAll(row, "@", "")
		row = strings.ReplaceAll(row, f.hardblank, " ")
		lines[i] = row
	}

	return lines, nil
}
