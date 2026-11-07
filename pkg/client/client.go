package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GladiaApiEndpoint is the default Gladia API base URL.
var GladiaApiEndpoint = "https://api.gladia.io"

const (
	// DefaultPollInterval is the delay between two polls of a transcription result.
	DefaultPollInterval = 1 * time.Second
	// DefaultPollTimeout bounds how long the client waits for a transcription
	// to complete before giving up.
	DefaultPollTimeout = 1 * time.Hour
	// pollRequestTimeout bounds each individual polling request.
	pollRequestTimeout = 60 * time.Second
	// maxRedirects mirrors net/http's default redirect limit.
	maxRedirects = 10
)

type GladiaClient struct {
	ApiKey         string
	GladiaEndpoint string
	CLIVersion     string
	httpClient     *http.Client
	Verbose        bool

	// PollInterval is the delay between two polls of a transcription result.
	// Zero means DefaultPollInterval.
	PollInterval time.Duration
	// PollTimeout is the maximum time spent waiting for a transcription result.
	// Zero means DefaultPollTimeout.
	PollTimeout time.Duration
}

func NewGladiaClient(apiKey string, verbose bool, cliVersion string) *GladiaClient {
	if cliVersion == "" {
		cliVersion = "dev"
	}
	c := &GladiaClient{
		ApiKey:         apiKey,
		GladiaEndpoint: GladiaApiEndpoint,
		CLIVersion:     cliVersion,
		Verbose:        verbose,
	}
	c.httpClient = &http.Client{CheckRedirect: c.checkRedirect}
	return c
}

func (c *GladiaClient) apiURL(path string) string {
	return strings.TrimSuffix(c.GladiaEndpoint, "/") + path
}

// checkRedirect refuses redirects that leave the configured Gladia API origin.
// net/http strips Authorization and Cookie on cross-host redirects, but not
// custom headers such as x-gladia-key, so following such a redirect would hand
// the API key to a third party. Redirects that downgrade the scheme are refused
// for the same reason.
func (c *GladiaClient) checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return fmt.Errorf("stopped after %d redirects", maxRedirects)
	}
	if !c.isGladiaURL(req.URL) {
		return fmt.Errorf("refusing redirect to %q: not the Gladia API origin", req.URL.Redacted())
	}
	return nil
}

// isGladiaURL reports whether u has the same origin (scheme, host and port) as
// the configured Gladia endpoint.
func (c *GladiaClient) isGladiaURL(u *url.URL) bool {
	if u == nil {
		return false
	}
	endpoint, err := url.Parse(c.GladiaEndpoint)
	if err != nil || endpoint.Host == "" {
		return false
	}
	return sameOrigin(u, endpoint)
}

func sameOrigin(a, b *url.URL) bool {
	if !strings.EqualFold(a.Scheme, b.Scheme) {
		return false
	}
	if !strings.EqualFold(a.Hostname(), b.Hostname()) {
		return false
	}
	return effectivePort(a) == effectivePort(b)
}

func effectivePort(u *url.URL) string {
	if p := u.Port(); p != "" {
		return p
	}
	switch strings.ToLower(u.Scheme) {
	case "https":
		return "443"
	case "http":
		return "80"
	}
	return ""
}

// validateResultURL checks that a URL returned by the API (such as result_url)
// is an absolute http(s) URL on the configured Gladia API origin before the
// client fetches it with the API key attached.
func (c *GladiaClient) validateResultURL(rawURL string) error {
	if rawURL == "" {
		return errors.New("API returned an empty result_url")
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("API returned an invalid result_url: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("API returned a result_url with unsupported scheme %q", parsed.Scheme)
	}
	if parsed.User != nil {
		return errors.New("API returned a result_url containing credentials")
	}
	if !c.isGladiaURL(parsed) {
		return fmt.Errorf("API returned a result_url outside the Gladia API origin (%s)", parsed.Host)
	}
	return nil
}

// setRequestHeaders sets the CLI headers. The API key is only attached to
// requests sent to the configured Gladia API origin.
func (c *GladiaClient) setRequestHeaders(req *http.Request) {
	if c.isGladiaURL(req.URL) {
		req.Header.Set("x-gladia-key", c.ApiKey)
	}
	req.Header.Set("x-gladia-version", "cli/"+c.CLIVersion)
}
