<div align="center">

<h1>brave-cli</h1>

<p>Search the web and extract page content from your terminal via the Brave Search API.</p>

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/gupsammy/brave-cli)](https://github.com/gupsammy/brave-cli/releases)

</div>

---

## Description

`brave-cli` is a Go command-line tool that wraps the [Brave Search API](https://api.search.brave.com) to bring web search directly to your terminal. It targets developers and power users who want scriptable, pipe-friendly access to search results — with optional full-page content extraction — without requiring a browser or Node.js runtime.

## Features

- Web search with up to 200 results via automatic pagination
- Optional full-page content extraction (`--content`) per search result
- Standalone URL content fetcher (`fetch`) — no API key required
- Multiple output formats: human-readable compact view or machine-readable JSON
- ANSI color output, suppressible with `--no-color` for clean pipe usage
- API key validation command (`auth`) to confirm setup in seconds
- Date-filtered search by freshness (`day`, `week`, `month`, `year`)
- Country and safe-search filtering for localized results

## Installation

### Quick Install (macOS / Linux)

```sh
curl -fsSL https://raw.githubusercontent.com/gupsammy/brave-cli/main/install.sh | sh
```

Detects your OS and architecture, downloads the pre-built binary from [GitHub Releases](https://github.com/gupsammy/brave-cli/releases), and installs it to `/usr/local/bin` (or `~/.local/bin` if `/usr/local/bin` is not writable).

### Via Go

Requires Go 1.21 or later:

```sh
go install github.com/gupsammy/brave-cli@latest
```

### Build from Source

```sh
git clone https://github.com/gupsammy/brave-cli.git
cd brave-cli
go build -o brave-cli .
```

### Prerequisites

- A [Brave Search API key](https://api-dashboard.search.brave.com/) (free tier available)

## Usage

```sh
# Basic search
brave-cli search "golang generics"

# Fetch up to 50 results as JSON
brave-cli search "site:github.com cobra cli" -n 50 -o json

# Search and extract full page text for each result
brave-cli search "rust async runtime" --content

# Extract readable text from a URL (no API key needed)
brave-cli fetch https://go.dev/blog/generics-proposal

# Verify your API key is configured correctly
brave-cli auth
# OK — key accepted. Test query returned 1 result(s).
```

## Configuration

`brave-cli` resolves `BRAVE_API_KEY` from the following sources in priority order:

| Priority | Source | Format |
|----------|--------|--------|
| 1 | Environment variable | already exported into the process |
| 2 | `~/.config/brave-cli/.env` | `BRAVE_API_KEY=your_key` |
| 3 | `~/.secrets` | `export BRAVE_API_KEY=your_key` |
| 4 | `~/.zshenv` | `export BRAVE_API_KEY=your_key` |
| 5 | `~/.zshrc` | `export BRAVE_API_KEY=your_key` |
| 6 | `~/.bashrc` | `export BRAVE_API_KEY=your_key` |
| 7 | `~/.bash_profile` | `export BRAVE_API_KEY=your_key` |
| 8 | `~/.profile` | `export BRAVE_API_KEY=your_key` |
| 9 | `~/.env` | `BRAVE_API_KEY=your_key` |

Both `KEY=value` and `export KEY=value` formats are accepted in any file. Quoted values (`"..."` or `'...'`) are also handled.

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--no-color` | `false` | Disable ANSI color in terminal output |

## API Reference

### `search <query>`

Search the web using the Brave Search API. Results go to stdout; progress messages (when `--content` is used) go to stderr, keeping the output pipe-friendly.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--count` | `-n` | `10` | Number of results to return (1–200; auto-paginated at 20 per page) |
| `--freshness` | | `""` | Date filter: `day`, `week`, `month`, `year` |
| `--country` | | `""` | 2-character country code (e.g. `US`, `GB`) |
| `--extra-snippets` | | `false` | Fetch additional text excerpts per result |
| `--content` | | `false` | Fetch and extract full page text for each result |
| `--safesearch` | | `moderate` | Content filter: `off`, `moderate`, `strict` |
| `--output` | `-o` | `compact` | Output format: `compact` or `json` |
| `--quiet` | `-q` | `false` | Suppress progress lines on stderr |

```sh
brave-cli search "openai gpt-4" -n 20 --freshness week -o json
```

### `fetch <url>`

Fetch a web page and extract clean, readable text from the HTML. Script/style blocks, navigation, and boilerplate are stripped. Does not require an API key.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `text` | Output format: `text` or `json` |

```sh
brave-cli fetch https://pkg.go.dev/net/http -o json
```

### `auth`

Runs a minimal test search to verify that `BRAVE_API_KEY` is set and accepted by the API.

```sh
brave-cli auth
# OK — key accepted. Test query returned 1 result(s).
```

## Acknowledgments

- [Brave Search API](https://api.search.brave.com) — the underlying search engine powering all search functionality
- [cobra](https://github.com/spf13/cobra) — CLI framework for Go
- [golang.org/x/net](https://pkg.go.dev/golang.org/x/net) — HTML parsing used for content extraction

## License

MIT — see [LICENSE](LICENSE) for full text.
