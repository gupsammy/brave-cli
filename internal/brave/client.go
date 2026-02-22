package brave

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const baseURL = "https://api.search.brave.com/res/v1"

var freshnessMap = map[string]string{
	"day":   "pd",
	"week":  "pw",
	"month": "pm",
	"year":  "py",
}

// Client wraps the Brave Search API.
type Client struct {
	apiKey string
	http   *http.Client
}

// NewClient returns a Client using the given API key.
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 30 * time.Second},
	}
}

// SearchResult is a normalised result item.
type SearchResult struct {
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	Description   string   `json:"description,omitempty"`
	ExtraSnippets []string `json:"extra_snippets,omitempty"`
	Content       string   `json:"content,omitempty"`
}

// SearchResponse is the top-level response returned by Search.
type SearchResponse struct {
	Items []SearchResult `json:"items"`
	Query string         `json:"query"`
	Total int            `json:"total"`
}

// SearchOpts controls the search request.
type SearchOpts struct {
	Count         int
	Freshness     string // day | week | month | year
	Country       string
	ExtraSnippets bool
	Safesearch    string // off | moderate | strict
}

// apiItem mirrors a single result from the Brave JSON payload.
type apiItem struct {
	Title         string   `json:"title"`
	URL           string   `json:"url"`
	Description   string   `json:"description"`
	ExtraSnippets []string `json:"extra_snippets"`
}

type apiPayload struct {
	Web struct {
		Results []apiItem `json:"results"`
	} `json:"web"`
}

// Search fetches up to opts.Count results, paginating automatically when
// the request exceeds the API's 20-result-per-page limit.
func (c *Client) Search(query string, opts SearchOpts) (*SearchResponse, error) {
	if opts.Count <= 0 {
		opts.Count = 10
	}
	if opts.Count > 100 {
		opts.Count = 100
	}
	if opts.Safesearch == "" {
		opts.Safesearch = "moderate"
	}

	var all []SearchResult
	remaining := opts.Count
	offset := 0

	for remaining > 0 {
		batch := remaining
		if batch > 20 {
			batch = 20
		}
		items, err := c.fetchPage(query, batch, offset, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, items...)
		if len(items) < batch {
			break // API returned fewer results than requested — no more pages
		}
		remaining -= batch
		offset += batch
	}

	return &SearchResponse{Items: all, Query: query, Total: len(all)}, nil
}

func (c *Client) fetchPage(query string, count, offset int, opts SearchOpts) ([]SearchResult, error) {
	params := url.Values{}
	params.Set("q", query)
	params.Set("count", strconv.Itoa(count))
	params.Set("offset", strconv.Itoa(offset))
	params.Set("text_decorations", "false")
	params.Set("safesearch", opts.Safesearch)

	if f, ok := freshnessMap[opts.Freshness]; ok {
		params.Set("freshness", f)
	}
	if opts.Country != "" {
		params.Set("country", strings.ToLower(opts.Country))
	}
	if opts.ExtraSnippets {
		params.Set("extra_snippets", "true")
	}

	req, err := http.NewRequest("GET", baseURL+"/web/search?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("X-Subscription-Token", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 || resp.StatusCode == 403 {
		return nil, fmt.Errorf("invalid API key (HTTP %d)", resp.StatusCode)
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var payload apiPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	out := make([]SearchResult, 0, len(payload.Web.Results))
	for _, r := range payload.Web.Results {
		out = append(out, SearchResult{
			Title:         r.Title,
			URL:           r.URL,
			Description:   r.Description,
			ExtraSnippets: r.ExtraSnippets,
		})
	}
	return out, nil
}

// FetchContent GETs a URL and returns readable plain text extracted from the HTML.
func (c *Client) FetchContent(pageURL string) (string, error) {
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024)) // 2 MB cap
	if err != nil {
		return "", err
	}

	return extractText(string(body)), nil
}

// extractText parses HTML and returns clean readable text.
func extractText(rawHTML string) string {
	doc, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return rawHTML // fall back to raw on parse error
	}

	var sb strings.Builder
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "noscript", "iframe", "head", "nav", "footer":
				return // skip entire subtree
			case "p", "div", "section", "article", "h1", "h2", "h3", "h4", "h5", "h6", "li", "tr", "br":
				sb.WriteString("\n")
			}
		}
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	walk(doc)

	// Collapse runs of whitespace / blank lines
	lines := strings.Split(sb.String(), "\n")
	out := make([]string, 0, len(lines))
	blanks := 0
	for _, l := range lines {
		l = strings.TrimRight(l, " \t")
		if l == "" {
			blanks++
			if blanks <= 1 {
				out = append(out, "")
			}
		} else {
			blanks = 0
			out = append(out, l)
		}
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}
