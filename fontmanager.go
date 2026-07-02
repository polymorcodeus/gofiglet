package gofiglet

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// defaultFont is the font name NewRenderOptions uses by default. It is
// not a fallback for missing fonts — getFont returns an error if a
// requested font can't be found, even if that font is defaultFont.
const defaultFont string = "standard"

// extension is the file extension used for figlet font files.
const extension string = "flf"

// embeddedFonts holds the builtin fonts (fonts/*.flf) embedded into the
// binary at compile time.
//
//go:embed fonts
var embeddedFonts embed.FS

// fontManager holds the available fonts: fontLib caches parsed *font
// objects by name, and fontList maps font names to on-disk paths that
// have been discovered but not yet parsed (see loadFontList).
type fontManager struct {
	fontLib  map[string]*font
	fontList map[string]string
}

// newFontManager creates a new fontManager, eagerly loading and parsing
// every embedded built-in font (from the "fonts" directory bundled via go:embed)
// into fontLib. It panics if an embedded font is missing or
// fails to parse, since that indicates a broken build rather than a
// runtime condition callers can recover from.
func newFontManager() *fontManager {
	fm := &fontManager{
		fontLib:  make(map[string]*font),
		fontList: make(map[string]string),
	}

	entries, err := fs.ReadDir(embeddedFonts, "fonts")
	if err != nil {
		panic("failed to read embedded fonts: " + err.Error())
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "."+extension) {
			continue
		}
		fontName := strings.TrimSuffix(entry.Name(), "."+extension)
		data, err := embeddedFonts.ReadFile("fonts/" + entry.Name())
		if err != nil {
			panic("failed to read embedded font " + entry.Name() + ": " + err.Error())
		}
		font, err := parseFontContent(string(data))
		if err != nil {
			panic("failed to parse embedded font " + entry.Name() + ": " + err.Error())
		}
		fm.fontLib[fontName] = font
	}

	return fm
}

// getFont returns the font registered under fontName. If fontName is not
// already cached in fontLib, getFont attempts to load it from fontList
// via loadDiskFont. It returns an error if fontName is not a known
// builtin font and was never registered via loadFontList (or fails to
// parse) — there is no silent fallback to defaultFont.
func (fm *fontManager) getFont(fontName string) (*font, error) {
	f, ok := fm.fontLib[fontName]
	if !ok {
		if err := fm.loadDiskFont(fontName); err != nil {
			return nil, fmt.Errorf("font %q: %w", fontName, err)
		}
		f = fm.fontLib[fontName]
	}

	return f, nil
}

// loadFontList walks fontPath recursively and records every ".flf" file
// found in fontList, keyed by font name (the filename without its
// extension) with the on-disk path as the value. It does not parse or
// load the fonts at this point; parsing is deferred to getFont /
// loadDiskFont for performance.
func (fm *fontManager) loadFontList(fontPath string) error {
	return filepath.Walk(fontPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(info.Name(), "."+extension) {
			return nil
		}
		fontName := strings.TrimSuffix(info.Name(), "."+extension)
		fm.fontList[fontName] = path

		return nil
	})
}

// loadDiskFont reads, parses, and caches the font registered under
// fontName in fontList. It returns an error if fontName was never
// registered via loadFontList, if the file cannot be read, or if
// parsing fails.
func (fm *fontManager) loadDiskFont(fontName string) error {
	path, ok := fm.fontList[fontName]
	if !ok {
		return errors.New("not registered in font list")
	}

	fontStr, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	font, err := parseFontContent(string(fontStr))
	if err != nil {
		return err
	}

	fm.fontLib[fontName] = font

	return nil
}

// parseFontContent parses the raw text content of an FLF font file into
// a *font. It is used for both disk-loaded and embedded builtin fonts.
func parseFontContent(cont string) (*font, error) {
	lines := strings.Split(cont, "\n")

	if len(lines) < 1 || len(lines[0]) < 5 {
		return nil, errors.New("font header missing or too short")
	}

	// FLF signature must start with "flf2a"
	if !strings.HasPrefix(lines[0], "flf2a") {
		return nil, fmt.Errorf("invalid font signature: expected flf2a, got %q", lines[0][:min(len(lines[0]), 5)])
	}

	header := strings.Fields(lines[0])

	// Minimum header fields: flf2a$ Height Baseline MaxLength OldLayout CommentLines
	// That's 6 whitespace-separated tokens (signature and hardblank counts as one).
	if len(header) < 6 {
		return nil, fmt.Errorf("font header has %d fields, expected at least 6", len(header))
	}

	hardblank := header[0][len(header[0])-1:]

	height, err := strconv.Atoi(header[1])
	if err != nil || height <= 0 {
		return nil, fmt.Errorf("invalid font height %q: %w", header[1], err)
	}

	// header[5] is Comment_Lines
	commentLines, err := strconv.Atoi(header[5])
	if err != nil || commentLines < 0 {
		return nil, fmt.Errorf("invalid comment line count %q: %w", header[5], err)
	}

	// Verify we have enough lines for the comment block + at least one glyph
	minLines := 1 + commentLines + height*95 // 95 printable ASCII chars (32-126)
	if len(lines) < minLines {
		return nil, fmt.Errorf("font content truncated: got %d lines, need at least %d", len(lines), minLines)
	}

	font := &font{
		hardblank: hardblank,
		height:    height,
		fontSlice: lines[1+commentLines:],
	}

	return font, nil
}
