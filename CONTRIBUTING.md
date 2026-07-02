# Contributing to gofiglet

Thank you for considering a contribution! This is a pure-Go library with zero external dependencies. Keeping it lightweight and predictable is the top priority.

## Quick Start

```bash
git clone <your-fork>
cd gofiglet
go build ./...
go test ./...
go vet ./...
```

## Project Structure

All source lives in the root package (`gofiglet`). There is no `main` package.

| File | What it does |
| ------ | ------------- |
| `banner.go` | High-level `Banner` API with functional options |
| `render.go` | Core `AsciiRender` engine |
| `fontmanager.go` | Font discovery (embedded eager, disk lazy) |
| `font.go` | FLF font parsing |
| `char.go` | Per-character rendering (`asciiChar`) |
| `color.go` | Color system (ANSI, TrueColor, hex) |
| `fonts/` | Bundled `.flf` files (embedded via `//go:embed`) |

## Conventions

### Go Version

Target **Go 1.26.4**. Avoid language features you aren't certain exist in this version. When in doubt, check [go.dev/doc/go1.26](https://go.dev/doc/go1.26).

### No External Dependencies

The module uses **stdlib only**. Do not add third-party packages to `go.mod`. If a dependency seems necessary, open an issue to discuss alternatives first.

### Error Handling

Return errors rather than silently falling back. `fontManager.getFont()` and `AsciiRender.RenderOpts()` propagate errors to the caller. Missing fonts are failures, not warnings.

### Functional Options

`Banner` uses the functional options pattern:

```go
func NewCmdBanner(title []string, options ...BannerOptions) (*Banner, error)
```

If you add a new option, follow the existing `WithXxx` naming and apply it in `NewCmdBanner` before the final validation step.

### Color Spelling

Use American spelling: `Color`, `ResolveColor`, `Colors` — not `Colour`.

### `strings.Builder`

Prefer `strings.Builder` with `WriteString`/`WriteByte` over `fmt.Fprintf` for string assembly in hot paths.

## Code Quality

Run the full check before pushing:

```bash
make check          # fmt, vet, lint, test
```

Or manually:

```bash
go fmt ./...
go vet ./...
golangci-lint run   # config: .golangci.yml
go test ./...
```

### Linting

We use `golangci-lint` with a custom config (`.golangci.yml`). Key enabled linters: `errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`, `misspell`.

### Tests

There are currently **zero tests**, thats on me. New features or bug fixes **should** include tests. Prefer table-driven tests. Target files for initial coverage:

- `parseFontContent` in `font.go`
- `ResolveColor` in `color.go`
- `RenderOpts` in `render.go`

## Fonts

### Embedded Fonts

Built-in fonts live in `fonts/` and are embedded at compile time via `//go:embed fonts`. They are eagerly loaded in `newFontManager()`. If you add a new bundled font, place the `.flf` in `fonts/` — no path configuration is needed.

### Disk Fonts

Additional fonts can be loaded lazily from disk via `AsciiRender.LoadFont(path)` or `WithLocalFont(name, path)`. Disk fonts are not parsed until first use.

### FLF Format

The parser expects the classic figlet header (`flf2a...`). See the detailed comment block in `font.go` for the header spec. Invalid fonts return clear errors; they do not fall back to a default.

## Pull Request Process

1. **Open an issue first** for significant changes (new API surface, breaking changes, new dependencies).
2. **Fork and branch**: `git checkout -b fix/description` or `feature/description`.
3. **Write tests** for any new behavior or bug fix.
4. **Run `make check`** and ensure everything passes.
5. **Update docs** if you change the public API (`Banner`, `AsciiRender`, `RenderOptions`, color functions).
6. **Squash** logically related commits if the history is noisy.
7. **Reference the issue** in your PR description.

## Release Notes

This project uses [GoReleaser](https://goreleaser.com/) and conventional commit grouping for changelogs. While we don't enforce commit message format in PRs, clean history helps. The changelog groups commits by:

- `feat:` → "New!"
- `fix:` → "Fixed"
- `docs:` → "Docs"
- `(deps)` → "Deps"
- everything else → "Other stuff"

## Questions?

Open a [Discussion](https://github.com/polymorcodeus/gofiglet/discussions) or issue. For bug reports, include:

- Go version (`go version`)
- Input that triggers the issue
- Expected vs actual output
- Font name (if relevant)
