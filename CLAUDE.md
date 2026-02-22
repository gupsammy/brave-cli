# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```sh
go build -o brave-cli .        # compile binary
go run . search "golang generics"  # run without building
go mod tidy                    # fix indirect dependency flags in go.mod
```

There are no tests yet. When adding them, use `go test ./...` to run all packages.

## Architecture

This is a Go CLI tool built with [cobra](https://github.com/spf13/cobra) that wraps the [Brave Search API](https://api.search.brave.com).

**Entry point:** `main.go` → `cmd.Execute()` → cobra dispatches to a subcommand.

**Subcommands** (`cmd/`):
- `search <query>` — calls the Brave Search API, optionally fetches full page content for each result
- `fetch <url>` — extracts readable text from an arbitrary URL (no API key required)
- `auth` — validates the API key with a minimal test search

**Internal packages** (`internal/`):
- `brave` — HTTP client for the Brave Search API + `FetchContent` (direct HTTP, not via Brave). The `Search` method paginates automatically: the Brave API's `offset` param is a zero-based page index (not an absolute item count), each page returns ≤20 results, max 10 pages (200 results).
- `config` — resolves `BRAVE_API_KEY` in priority order: env var → `~/.config/brave-cli/.env` → `~/.secrets` → `~/.zshenv` → `~/.zshrc` → `~/.bashrc` → `~/.bash_profile` → `~/.profile` → `~/.env` (all use `export KEY=value` or bare `KEY=value` format)
- `output` — terminal rendering (`PrintCompact` with ANSI colors) and JSON serialization for both search and fetch results. `NoColor` is a package-level bool toggled by `--no-color`.

**Key conventions:**
- All `RunE` handlers set `SilenceUsage: true` so cobra never dumps the usage banner on runtime errors — only on argument/flag parse errors.
- Progress messages (e.g., "fetching …") go to stderr; all data output goes to stdout. This makes the tool pipe-friendly.
- `FetchContent` loads the API key for API-key consistency across subcommands but does not actually use it — content fetching is a direct HTTP GET.
- The Free AI plan allows 1 req/s; pagination inserts a 1.2-second delay between pages (`pagingDelay` constant in `internal/brave/client.go`).
