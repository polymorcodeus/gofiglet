package gofiglet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"image/color"
	"regexp"
	"strings"
)

// escape is the ANSI escape character used to build all terminal color
// sequences in this file.
const escape = "\x1b"

// Standard 16-color ANSI terminal colors (normal intensity) represented
// as standard image/color RGBA values.
var (
	ColorBlack     = color.RGBA{0, 0, 0, 255}
	ColorRed       = color.RGBA{170, 0, 0, 255}
	ColorGreen     = color.RGBA{0, 170, 0, 255}
	ColorYellow    = color.RGBA{170, 85, 0, 255}
	ColorBlue      = color.RGBA{0, 0, 170, 255}
	ColorMagenta   = color.RGBA{170, 0, 170, 255}
	ColorCyan      = color.RGBA{0, 170, 170, 255}
	ColorWhite     = color.RGBA{170, 170, 170, 255}
	ColorHiBlack   = color.RGBA{85, 85, 85, 255}
	ColorHiRed     = color.RGBA{255, 85, 85, 255}
	ColorHiGreen   = color.RGBA{85, 255, 85, 255}
	ColorHiYellow  = color.RGBA{255, 255, 85, 255}
	ColorHiBlue    = color.RGBA{85, 85, 255, 255}
	ColorHiMagenta = color.RGBA{255, 85, 255, 255}
	ColorHiCyan    = color.RGBA{85, 255, 255, 255}
	ColorHiWhite   = color.RGBA{255, 255, 255, 255}
)

// Preset 24-bit TrueColor values used as defaults elsewhere in the
// package (e.g. Colors["default"], NewCmdBanner's default palette).
var (
	TrueColorPink206    = color.RGBA{255, 95, 175, 255}
	TrueColorYellowNeon = color.RGBA{207, 255, 4, 255}
	TrueColorGold       = color.RGBA{255, 215, 64, 255}
)

// ColorNone is the no-op color; using it renders text without any ANSI
// color escape sequences.
var ColorNone color.Color = nil

// Colors maps human-friendly color names to color.Color values, used by
// ResolveColor for named lookups. Named entries mirror the variables
// declared above.
var Colors = map[string]color.Color{
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

// ansiColorCodes maps the predefined RGBA colors to their legacy ANSI
// escape codes for terminals that do not support 24-bit true color.
var ansiColorCodes = map[color.RGBA]int{
	{0, 0, 0, 255}:       30,
	{170, 0, 0, 255}:     31,
	{0, 170, 0, 255}:     32,
	{170, 85, 0, 255}:    33,
	{0, 0, 170, 255}:     34,
	{170, 0, 170, 255}:   35,
	{0, 170, 170, 255}:   36,
	{170, 170, 170, 255}: 37,
	{85, 85, 85, 255}:    90,
	{255, 85, 85, 255}:   91,
	{85, 255, 85, 255}:   92,
	{255, 255, 85, 255}:  93,
	{85, 85, 255, 255}:   94,
	{255, 85, 255, 255}:  95,
	{85, 255, 255, 255}:  96,
	{255, 255, 255, 255}: 97,
}

// GetPrefix returns the ANSI escape sequence that switches the terminal
// to the given color. If c is nil, it returns an empty string.
func GetPrefix(c color.Color) string {
	if c == nil {
		return ""
	}
	if code, ok := ansiCodeFor(c); ok {
		return fmt.Sprintf("%v[0;%dm", escape, code)
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("%v[38;2;%d;%d;%dm", escape, uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

// GetColorCode returns the raw ANSI color code as a string, without any
// escape sequence wrapping. If c is nil, it returns an empty string.
func GetColorCode(c color.Color) string {
	if c == nil {
		return ""
	}
	if code, ok := ansiCodeFor(c); ok {
		return fmt.Sprintf("%d", code)
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("38;2;%d;%d;%d", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

// GetSuffix returns the ANSI escape sequence that resets terminal
// formatting back to default. If c is nil, it returns an empty string.
func GetSuffix(c color.Color) string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("%v[0m", escape)
}

// ansiCodeFor looks up the legacy ANSI color code for a predefined
// palette color. If the color is not in the predefined palette, it
// returns ok == false so the caller can fall back to 24-bit true color.
func ansiCodeFor(c color.Color) (int, bool) {
	r, g, b, a := c.RGBA()
	rgba := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
	code, ok := ansiColorCodes[rgba]
	return code, ok
}

// ResolveColor returns a color.Color by named lookup or hex string
// (#RRGGBB). Falls back to TrueColorPink206 if the input is
// unrecognized.
func ResolveColor(c string) color.Color {
	var hex6 = regexp.MustCompile(`^#?[0-9a-fA-F]{6}$`)

	if lookup, ok := Colors[c]; ok {
		return lookup
	}
	if hex6.MatchString(c) {
		if newHexColor, err := NewColorFromHexString(c); err == nil {
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

// NewColorFromHexString returns a color.Color parsed from a hex string.
func NewColorFromHexString(c string) (color.Color, error) {
	rgb, err := hexToRGB(c)
	if err != nil {
		return nil, err
	}

	return color.RGBA{rgb[0], rgb[1], rgb[2], 255}, nil
}
