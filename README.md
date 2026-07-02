# gofiglet

[![Go Version](https://img.shields.io/github/go-mod/go-version/polymorcodeus/gofiglet)](https://go.dev/) [![License](https://img.shields.io/github/license/polymorcodeus/gofiglet)](./LICENSE) [![Build Status](https://img.shields.io/github/actions/workflow/status/polymorcodeus/gofiglet/ci.yml?branch=main)](https://github.com/polymorcodeus/gofiglet/actions)[![Go Report Card](https://goreportcard.com/badge/github.com/polymorcodeus/gofiglet)](https://goreportcard.com/report/github.com/polymorcodeus/gofiglet)

`gofiglet` is a pure Go library for rendering ASCII art text from [figlet](http://www.figlet.org/) (`.flf`) fonts. It
supports ANSI colors, 24-bit true color, and per-character coloring, with a set of fonts bundled and embedded directly
into the package.

It's a reimagining of [mbndr/figlet4go](https://github.com/mbndr/figlet4go) with updated conventions, error handling, and unlike it's predecessor, there is no CLI. I've updated it primarily to display my CLI commands as a figlet banner, so those features have been prioritized.

## Features

- Render any string to ASCII art using classic FIGfont (`.flf`) files
- Builtin fonts embedded at compile time (`standard`, `small`, `ogre`, `smallsmursh`) — works out of the box with no
filesystem setup
- Load additional fonts from disk at runtime, download from [figlet](http://www.figlet.org/)
- ANSI 16-color, 24-bit true color, and no-color output
- Named colors, hex color strings (`#RRGGBB`), and per-character color cycling
- A high-level `Banner` API for quickly building colored CLI banners
- Stdlib only — no external dependencies

## Installation

```bash
go get github.com/polymorcodeus/gofiglet
```

## Quick Start

### Rendering with `AsciiRender`

`AsciiRender` is the core rendering engine. It comes preloaded with the embedded builtin fonts.

```go
package main

import (
  "fmt"

  "github.com/polymorcodeus/gofiglet"
)

func main() {
  ascii := gofiglet.NewAsciiRender()

  // Render with the default font ("standard"), no color.
  out, err := ascii.Render("Hello")
  if err != nil {
    panic(err)
  }
  fmt.Println(out)
}
```

### Rendering with options

Use `RenderOpts` to choose a font and apply color:

```go
ascii := gofiglet.NewAsciiRender()

opts := gofiglet.NewRenderOptions()
opts.FontName = "smallsmursh"
opts.FontColor = []gofiglet.Color{
  gofiglet.ColorCyan,
  gofiglet.ResolveColor("#ff5fAF") ,
}

out, err := ascii.RenderOpts("Hi!", opts)
if err != nil {
  panic(err)
}
fmt.Println(out)
```

`FontColor` is applied cyclically across the characters of the rendered string — with two colors and four characters,
colors alternate `0, 1, 0, 1`.

### Loading fonts from disk

Builtin fonts are embedded, but you can register additional `.flf` fonts from a directory at runtime:

```go
ascii := gofiglet.NewAsciiRender()
if err := ascii.LoadFont("/path/to/fonts"); err != nil {
  panic(err)
}

opts := gofiglet.NewRenderOptions()
opts.FontName = "my-custom-font"

out, _ := ascii.RenderOpts("Custom", opts)
fmt.Println(out)
```

`LoadFont` walks the directory recursively and registers every `.flf` file it finds, keyed by filename (without extension). Fonts aren't parsed until they're actually requested.

> **Note:** if a requested font can't be found (never loaded, or misspelled), `RenderOpts` (and anything built on it, like `Banner`) returns an error naming the missing font.

### The `Banner` convenience API

`Banner` is a higher-level wrapper aimed at CLI tool banners — multiple title segments, each with its own color:

```go
package main

import "github.com/polymorcodeus/gofiglet"

func main() {
  b, err := gofiglet.NewCmdBanner(
    []string{"my-cli", " sub"},
    gofiglet.WithColors("cyan", "pink"),
    gofiglet.WithFont("smallsmursh"),
  )
  if err != nil {
    panic(err)
  }

  if _, err := gofiglet.PrintCmdBanner(b); err != nil {
    panic(err)
  }
}
```

- `Colors` must have the same number of entries as `Title` — `NewCmdBanner` returns an error otherwise.
- `Title` segments are concatenated with no separator, so include spacing in the segments themselves if you want it.
- `TopPadding` defaults to `true`, which adds a single leading newline before the rendered output (use `WithZeroPadding()` to disable this).
- `CmdBanner` and `PrintCmdBanner` both return an error if `FontPath` fails to load or rendering fails (e.g. `FontName` can't be found).

#### Banner functional options

| Option | Effect |
| --- | --- |
| `WithColors(colors ...string)` | Sets the color palette, one per `Title` segment. Each string is resolved with `ResolveColor`. |
| `WithFont(name string)` | Selects a builtin or already-loaded font by name. |
| `WithLocalFont(name, path string)` | Sets the font name *and* a directory to load additional fonts from. |
| `WithZeroPadding()` | Disables the default leading newline (`TopPadding = false`). |

## Colors

Colors implement a common `Color` interface (`GetPrefix`, `GetSuffix`, `GetColorCode`), with three built-in implementations:

- `AnsiColor` — standard 16-color ANSI terminal colors (e.g. `ColorRed`, `ColorHiCyan`)
- `TrueColor` — 24-bit RGB colors (e.g. `TrueColorPink206`, or any hex string)
- `NoColor` — a no-op that emits no escape sequences

### Resolving colors by name or hex

```go
c1 := gofiglet.ResolveColor("cyan")      // named lookup
c2 := gofiglet.ResolveColor("#ff5fAF") // hex string
c3 := gofiglet.ResolveColor("bogus")     // falls back to TrueColorPink206
```

`ResolveColor` checks, in order: a named lookup in the `Colors` map, then a `#RRGGBB` / `RRGGBB` hex string, then falls back to `TrueColorPink206` if nothing matches.

See the `Colors` map in `color.go` for the full list of named colors (standard and high-intensity ANSI names, plus a few named true colors like `"pink"`, `"gold"`, `"neonyellow"`).

### Custom hex colors

```go
c, err := gofiglet.NewTrueColorFromHexString("#39ff14")
if err != nil {
  panic(err)
}
```

## Bundled fonts

The following fonts are embedded in the binary and available by name without any setup: `standard`, `small`, `ogre`, `smallsmursh`. Once kerning is supported the duplicate
fonts will be removed.

## Project status

This is a works-for-me library:

- There are currently no automated tests.
- Kerning is not yet supported.
- No autodection for ignoring `PrintCmdBanner` when running headless.

## License

[MIT](LICENSE)
