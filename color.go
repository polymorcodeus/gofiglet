package gofiglet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// escape is the ANSI escape character used to build all terminal color
// sequences in this file.
const escape = "\x1b"

// Standard 16-color ANSI terminal colors (normal intensity).
var (
	ColorBlack     AnsiColor = AnsiColor{30}
	ColorRed       AnsiColor = AnsiColor{31}
	ColorGreen     AnsiColor = AnsiColor{32}
	ColorYellow    AnsiColor = AnsiColor{33}
	ColorBlue      AnsiColor = AnsiColor{34}
	ColorMagenta   AnsiColor = AnsiColor{35}
	ColorCyan      AnsiColor = AnsiColor{36}
	ColorWhite     AnsiColor = AnsiColor{37}
	ColorHiBlack   AnsiColor = AnsiColor{90}
	ColorHiRed     AnsiColor = AnsiColor{91}
	ColorHiGreen   AnsiColor = AnsiColor{92}
	ColorHiYellow  AnsiColor = AnsiColor{93}
	ColorHiBlue    AnsiColor = AnsiColor{94}
	ColorHiMagenta AnsiColor = AnsiColor{95}
	ColorHiCyan    AnsiColor = AnsiColor{96}
	ColorHiWhite   AnsiColor = AnsiColor{97}
)

// Preset 24-bit TrueColor values used as defaults elsewhere in the
// package (e.g. Colors["default"], NewCmdBanner's default palette).
var (
	TrueColorPink206    TrueColor = TrueColor{r: 255, g: 95, b: 175}
	TrueColorYellowNeon TrueColor = TrueColor{r: 207, g: 255, b: 4}
	TrueColorGold       TrueColor = TrueColor{r: 255, g: 215, b: 64}
)

// ColorNone is the no-op Color; using it renders text without any ANSI
// color escape sequences.
var (
	ColorNone NoColor = NoColor{}
)

// Colors maps human-friendly color names to Color values, used by
// ResolveColor for named lookups. Named entries mirror the ANSI/
// TrueColor variables declared above.
var Colors = map[string]Color{
	"default": TrueColorPink206,
	"none":    ColorNone,
	// ANSI
	"black":   ColorBlack,
	"red":     ColorRed,
	"green":   ColorGreen,
	"yellow":  ColorYellow,
	"blue":    ColorBlue,
	"magenta": ColorMagenta,
	"cyan":    ColorCyan,
	"white":   ColorWhite,
	// High Intensity ANSI
	"darkGray":     ColorHiBlack,
	"lightRed":     ColorHiRed,
	"lightGreen":   ColorHiGreen,
	"lightYellow":  ColorHiYellow,
	"lightBlue":    ColorHiBlue,
	"lightMagenta": ColorHiMagenta,
	"lightCyan":    ColorHiCyan,
	"lightWhite":   ColorHiWhite,
	// True Colors
	"pink":       TrueColorPink206,
	"neonyellow": TrueColorYellowNeon,
	"gold":       TrueColorGold,
}

// Color wraps ANSI escape sequences for terminal coloring.
type Color interface {
	GetPrefix() string
	GetSuffix() string
	GetColorCode() string
}

// AnsiColor is a standard 16-color ANSI terminal color.
type AnsiColor struct {
	code int
}

// GetPrefix returns the ANSI escape sequence that switches the terminal
// to this color.
func (ac AnsiColor) GetPrefix() string {
	return fmt.Sprintf("%v[0;%dm", escape, ac.code)
}

// GetColorCode returns the raw ANSI color code as a string, without any
// escape sequence wrapping.
func (ac AnsiColor) GetColorCode() string {
	return fmt.Sprintf("%d", ac.code)
}

// GetSuffix returns the ANSI escape sequence that resets terminal
// formatting back to default.
func (ac AnsiColor) GetSuffix() string {
	return fmt.Sprintf("%v[0m", escape)
}

// TrueColor is a 24-bit RGB terminal color.
type TrueColor struct {
	r int
	g int
	b int
}

// GetPrefix returns the ANSI 24-bit escape sequence that switches the
// terminal to this RGB color.
func (tc TrueColor) GetPrefix() string {
	return fmt.Sprintf("%v[38;2;%d;%d;%dm", escape, tc.r, tc.g, tc.b)
}

// GetColorCode returns the raw ANSI 24-bit color code as a string,
// without any escape sequence wrapping.
func (tc TrueColor) GetColorCode() string {
	return fmt.Sprintf("38;2;%d;%d;%d", tc.r, tc.g, tc.b)
}

// GetSuffix returns the ANSI escape sequence that resets terminal
// formatting back to default.
func (tc TrueColor) GetSuffix() string {
	return fmt.Sprintf("%v[0m", escape)
}

// NoColor is a no-op Color that produces no escape sequences.
type NoColor struct{}

// GetPrefix returns an empty string; NoColor applies no formatting.
func (n NoColor) GetPrefix() string { return "" }

// GetColorCode returns an empty string; NoColor has no underlying code.
func (n NoColor) GetColorCode() string { return "" }

// GetSuffix returns an empty string; NoColor applies no formatting.
func (n NoColor) GetSuffix() string { return "" }

// ResolveColor returns a Color by named lookup or hex string (#RRGGBB).
// Falls back to TrueColorPink206 if the input is unrecognized.
func ResolveColor(c string) Color {
	var hex6 = regexp.MustCompile(`^#?[0-9a-fA-F]{6}$`)

	if lookup, ok := Colors[c]; ok {
		return lookup
	}
	if hex6.MatchString(c) {
		if newHexColor, err := NewTrueColorFromHexString(c); err == nil {
			return newHexColor
		}
	}

	return Colors["default"]
}

// hexToRGB decodes a hex color string (with or without a leading "#")
// into its raw RGB bytes.
func hexToRGB(c string) ([]byte, error) {
	trimHex := strings.TrimPrefix(c, "#")
	rgb, err := hex.DecodeString(trimHex)
	if err != nil {
		return nil, errors.New("Invalid color given (" + c + ")")
	}
	return rgb, nil
}

// NewTrueColorFromHexString returns a TrueColor parsed from a hex string.
func NewTrueColorFromHexString(c string) (*TrueColor, error) {
	rgb, err := hexToRGB(c)
	if err != nil {
		return nil, err
	}

	return &TrueColor{
		r: int(rgb[0]),
		g: int(rgb[1]),
		b: int(rgb[2]),
	}, nil
}
