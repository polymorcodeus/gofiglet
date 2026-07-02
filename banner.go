// Package gofiglet renders ASCII art text using figlet fonts (.flf).
// It supports ANSI colors, true color (RGB), and per-character coloring.
// Built-in fonts are embedded and available without external file paths.
package gofiglet

import (
	"fmt"
	"os"
	"strings"
)

// Banner holds configuration for rendering a multi-segment ASCII banner.
// Each entry in Title is rendered with the corresponding color from Colors.
type Banner struct {
	// Title holds the banner's text segments. Segments are concatenated
	// with no separator before rendering; each segment is colored
	// independently via Colors.
	Title []string
	// Colors holds one color per Title segment, applied cyclically by
	// index (segment i gets Colors[i % len(Colors)]). NewCmdBanner
	// requires len(Colors) == len(Title).
	Colors []Color
	// FontName is the figlet font to render with, by name.
	FontName string
	// FontPath, if set, is an on-disk directory to load additional fonts
	// from (in addition to the embedded builtin fonts) before rendering.
	FontPath string
	// TopPadding, if true, adds a single leading newline before the
	// rendered output. It does not affect kerning or layout.
	TopPadding bool
}

// BannerOptions configures a Banner via the functional options pattern.
type BannerOptions func(b *Banner)

// WithColors sets the color palette for each Title segment.
func WithColors(colors ...string) BannerOptions {
	return func(b *Banner) {
		b.Colors = make([]Color, len(colors))
		for i, c := range colors {
			b.Colors[i] = ResolveColor(c)
		}
	}
}

// WithFont sets the figlet font name to use for rendering.
func WithFont(f string) BannerOptions {
	return func(b *Banner) {
		b.FontName = f
	}
}

// WithLocalFont sets the font name and loads additional fonts from a local directory.
func WithLocalFont(f string, p string) BannerOptions {
	return func(b *Banner) {
		b.FontName = f
		b.FontPath = p
	}
}

// WithZeroPadding disables the leading newline added by default when TopPadding is true.
func WithZeroPadding() BannerOptions {
	return func(b *Banner) {
		b.TopPadding = false
	}
}

// NewCmdBanner creates a Banner with sensible defaults for CLI tool banners.
// Title entries represent command and subcommand names (e.g., ["cmd", "sub"]).
// Colors must match the number of Title entries.
func NewCmdBanner(title []string, options ...BannerOptions) (*Banner, error) {
	b := &Banner{
		Title:      title,
		Colors:     []Color{ColorCyan, TrueColorPink206},
		FontName:   "smallsmursh",
		TopPadding: true,
	}
	for _, o := range options {
		o(b)
	}

	if len(b.Colors) != len(b.Title) {
		return nil, fmt.Errorf(
			"banner has %d title entries but %d colors provided; counts must match",
			len(b.Title), len(b.Colors),
		)
	}

	if b.FontPath != "" {
		if err := verifyExists(b.FontPath); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// CmdBanner renders b as a colored ASCII art string. It returns an error
// if b.FontPath is set but fails to load, or if rendering fails (e.g.
// b.FontName cannot be found, or a Title segment contains a non-ASCII
// character).
func CmdBanner(b *Banner) (string, error) {
	ascii := NewASCIIRender()
	if b.FontPath != "" {
		if err := ascii.LoadFont(b.FontPath); err != nil {
			return "", err
		}
	}

	figletOptions := NewRenderOptions()
	figletOptions.FontName = b.FontName

	bannerTitle := strings.Join(b.Title, "")

	figletColors := make([]Color, 0, len(bannerTitle))
	for i, entry := range b.Title {
		color := b.Colors[i%len(b.Colors)]
		for range entry {
			figletColors = append(figletColors, color)
		}
	}
	figletOptions.FontColor = figletColors

	var renderedString strings.Builder
	if b.TopPadding {
		renderedString.WriteByte('\n')
	}

	asciiString, err := ascii.RenderOpts(bannerTitle, figletOptions)
	if err != nil {
		return "", err
	}
	renderedString.WriteString(asciiString)

	return renderedString.String(), nil
}

// PrintCmdBanner renders and prints b to stdout. It returns an error if
// rendering fails; see CmdBanner.
func PrintCmdBanner(b *Banner) (int, error) {
	s, err := CmdBanner(b)
	if err != nil {
		return 0, err
	}
	return fmt.Print(s)
}

// verifyExists returns an error if filename does not exist on disk (e.g.
// os.ErrNotExist), or nil if it does.
func verifyExists(filename string) error {
	_, err := os.Stat(filename)
	return err
}
