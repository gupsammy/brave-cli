<div align="center">

<h1>brave-cli</h1>

<p>Search the web and extract page content from your terminal via the Brave Search API.</p>

[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/gupsammy/brave-cli?color=brightgreen)](https://github.com/gupsammy/brave-cli/releases)
[![Go](https://img.shields.io/badge/go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev)

</div>

---

`brave-cli` is a Go command-line tool that brings web search directly to your terminal. It wraps the [Brave Search API](https://api.search.brave.com) to deliver scriptable, pipe-friendly search results — with optional full-page content extraction — without a browser or Node.js runtime. Built for developers and power users who want to search, filter, and pipe results programmatically.

<a id="features"></a>

## ✨ Features

- **Web search** with up to 200 results via automatic pagination
- **Full-page content extraction** (`--content`) — fetch and strip each result to readable text
- **Standalone URL fetcher** (`fetch`) — extract readable text from any URL, no API key required
- **Multiple output formats** — human-readable compact view or machine-readable JSON
- **ANSI color output** with `--no-color` for clean pipe usage
- **API key validation** (`auth`) — confirm your key is active in one command
- **Date-filtered search** by freshness: `day`, `week`, `month`, `year`
- **Country and safe-search filtering** for localized, audience-appropriate results

<a id="quickstart"></a>

## 🚀 Quick Start

### Installation

**macOS / Linux — one command:**

```sh
curl -fsSL https://raw.githubusercontent.com/gupsammy/brave-cli/main/install.sh | sh
```

**Via Go** (requires Go 1.21+):

```sh
go install github.com/gupsammy/brave-cli@latest
```

**Build from source:**

```sh
git clone https://github.com/gupsammy/brave-cli.git
cd brave-cli
go build -o brave-cli .
```

### Get an API Key

Sign up for a free key at [api-dashboard.search.brave.com](https://api-dashboard.search.brave.com/). The free tier allows 1 request/second.

Set it in your shell config (or any of the [supported sources](#configuration)):

```sh
echo 'export BRAVE_API_KEY=your_key_here' >> ~/.zshrc
```

Then verify it works:

```sh
brave-cli auth
# OK — key accepted. Test query returned 1 result(s).
```

### Basic Usage

```sh
# Search the web
brave-cli search "golang generics"

# Get 50 results as JSON
brave-cli search "site:github.com cobra cli" -n 50 -o json

# Search with full page content for each result
brave-cli search "rust async runtime" --content

# Extract readable text from any URL
brave-cli fetch https://go.dev/blog/generics-proposal
```

<a id="configuration"></a>

## ⚙️ Configuration

`brave-cli` resolves `BRAVE_API_KEY` from the following sources in priority order:

| Priority | Source | Format |
|:--------:|--------|--------|
| 1 | Environment variable | already exported into the process |
| 2 | `~/.config/brave-cli/.env` | `BRAVE_API_KEY=your_key` |
| 3 | `~/.secrets` | `export BRAVE_API_KEY=your_key` |
| 4 | `~/.zshenv` | `export BRAVE_API_KEY=your_key` |
| 5 | `~/.zshrc` | `export BRAVE_API_KEY=your_key` |
| 6 | `~/.bashrc` | `export BRAVE_API_KEY=your_key` |
| 7 | `~/.bash_profile` | `export BRAVE_API_KEY=your_key` |
| 8 | `~/.profile` | `export BRAVE_API_KEY=your_key` |
| 9 | `~/.env` | `BRAVE_API_KEY=your_key` |

Both `KEY=value` and `export KEY=value` formats are accepted. Quoted values (`"..."` or `'...'`) are handled.

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--no-color` | `false` | Disable ANSI color in terminal output |

<a id="api"></a>

## 📖 API Reference

### `brave-cli search <query>`

Search the web using the Brave Search API. Results go to stdout; progress messages (when `--content` is used) go to stderr, keeping output pipe-friendly.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--count` | `-n` | `10` | Number of results (1–200; auto-paginated at 20/page) |
| `--freshness` | | `""` | Date filter: `day`, `week`, `month`, `year` |
| `--country` | | `""` | 2-character country code (e.g. `US`, `GB`) |
| `--extra-snippets` | | `false` | Fetch additional text excerpts per result |
| `--content` | | `false` | Fetch and extract full page text for each result |
| `--safesearch` | | `moderate` | Content filter: `off`, `moderate`, `strict` |
| `--output` | `-o` | `compact` | Output format: `compact` or `json` |
| `--quiet` | `-q` | `false` | Suppress progress lines on stderr |

```sh
# 20 results from the past week, as JSON
brave-cli search "openai gpt-4" -n 20 --freshness week -o json

# Pipe search results into another tool
brave-cli search "rust crates 2025" -o json | jq '.[].url'

# Search and read every result page
brave-cli search "go embed tutorial" --content -q
```

---

### `brave-cli fetch <url>`

Fetch a web page and extract clean, readable text from the HTML. Script/style blocks, navigation, and boilerplate are stripped. Does not require an API key.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--output` | `-o` | `text` | Output format: `text` or `json` |

```sh
brave-cli fetch https://pkg.go.dev/net/http
brave-cli fetch https://pkg.go.dev/net/http -o json
```

---

### `brave-cli auth`

Runs a minimal test search to verify that `BRAVE_API_KEY` is set and accepted by the API. Exits with code `0` on success, `3` on failure.

```sh
brave-cli auth
# OK — key accepted. Test query returned 1 result(s).
# First result: https://example.com/...
```

<a id="roadmap"></a>

## 🗺 Roadmap

- [ ] `news` subcommand — dedicated News Search endpoint (up to 50 results/page)
- [ ] `images` subcommand — Image Search with URL, thumbnail, and dimensions output
- [ ] `--search-lang` flag — filter results by page language (ISO 639-1 code)
- [ ] Custom freshness date ranges — `--since 2025-01-01 --until 2025-06-30`
- [ ] Rich results — structured weather, stocks, and knowledge panel cards in terminal
- [ ] Local / POI search — business addresses, hours, and phone via `/local/pois`

<a id="faq"></a>

## ❓ FAQ

### Do I need Go installed to use brave-cli?

No — the [Quick Install](#quickstart) curl command downloads a pre-built binary for your OS and architecture. Go is only required if you want to install via `go install` or build from source.

### How many requests can I make?

The Brave Search API free tier allows 1 request per second and 2,000 queries per month. `brave-cli` automatically inserts a 1.2-second delay between pagination requests to respect this limit.

### Why does `brave-cli search` print progress to stderr instead of stdout?

Progress messages (e.g. `fetching https://...`) go to stderr so they don't contaminate data output on stdout. This means `brave-cli search "..." --content -o json | jq ...` works cleanly in pipes without the progress lines polluting the JSON stream.

<a id="acknowledgments"></a>

## 💙 Acknowledgments

- [Brave Search API](https://api.search.brave.com) — the underlying search engine powering all search functionality
- [cobra](https://github.com/spf13/cobra) — CLI framework for Go
- [golang.org/x/net](https://pkg.go.dev/golang.org/x/net) — HTML parsing for content extraction

<a id="license"></a>

## 📄 License

MIT — see [LICENSE](LICENSE) for full text.
