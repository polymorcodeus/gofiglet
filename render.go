package gofiglet

import (
	"image/color"
	"strings"
)

// RenderOptions configures a single ASCIIRender.RenderOpts call: which
// font to render with, and optionally a per-character color cycle.
type RenderOptions struct {
	// FontName selects the font to render with. If the named font
	// cannot be found, RenderOpts returns an error.
	FontName string
	// FontColor, if non-empty, is applied cyclically across the
	// characters of the rendered string (character i gets
	// FontColor[i % len(FontColor)]). If empty, no color is applied.
	FontColor []color.Color
}

// NewRenderOptions creates a new RenderOptions with FontName set to
// defaultFont ("standard") and no FontColor.
func NewRenderOptions() *RenderOptions {
	return &RenderOptions{
		FontName: defaultFont,
	}
}

// ASCIIRender is the core rendering engine. It wraps a fontManager and
// exposes methods to render strings to ASCII art.
type ASCIIRender struct {
	fontMgr *fontManager
}

// NewASCIIRender creates a new ASCIIRender with a fresh fontManager,
// preloaded with the embedded builtin fonts.
func NewASCIIRender() *ASCIIRender {
	return &ASCIIRender{
		fontMgr: newFontManager(),
	}
}

// LoadFont registers all *.flf font files found recursively under
// fontPath, making them available for later rendering by name. Fonts
// are discovered but not parsed until they are actually requested.
func (ar *ASCIIRender) LoadFont(fontPath string) error {
	return ar.fontMgr.loadFontList(fontPath)
}

// Render renders str using default RenderOptions (the default font, no
// color). It is a convenience wrapper around RenderOpts.
func (ar *ASCIIRender) Render(str string) (string, error) {
	return ar.RenderOpts(str, NewRenderOptions())
}

// RenderOpts renders str as ASCII art according to opt, returning the
// fully composed multi-line output (including a trailing newline after
// each glyph row). It returns an error if opt.FontName cannot be found,
// or if str contains a rune outside the printable ASCII range (0-127).
func (ar *ASCIIRender) RenderOpts(str string, opt *RenderOptions) (string, error) {
	colored := len(opt.FontColor) > 0

	font, err := ar.fontMgr.getFont(opt.FontName)
	if err != nil {
		return "", err
	}

	chars := []*asciiChar{}

	curColorIndex := 0

	for _, char := range str {
		asciiChar, err := newASCIIChar(font, char)
		if err != nil {
			return "", err
		}

		if colored {
			if curColorIndex == len(opt.FontColor) {
				curColorIndex = 0
			}
			asciiChar.color = opt.FontColor[curColorIndex]
			curColorIndex++
		}

		chars = append(chars, asciiChar)
	}

	var result strings.Builder
	for curLine := 0; curLine < font.height; curLine++ {
		for i := range chars {
			result.WriteString(chars[i].GetLine(curLine))
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}
