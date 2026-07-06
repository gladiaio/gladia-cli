package client

import (
	"net/http"
	"strings"
)

// GladiaApiEndpoint is the default Gladia API base URL.
var GladiaApiEndpoint = "https://api.gladia.io"

type GladiaClient struct {
	ApiKey         string
	GladiaEndpoint string
	CLIVersion     string
	httpClient     *http.Client
	Verbose        bool
}

func NewGladiaClient(apiKey string, verbose bool, cliVersion string) *GladiaClient {
	if cliVersion == "" {
		cliVersion = "dev"
	}
	return &GladiaClient{
		ApiKey:         apiKey,
		GladiaEndpoint: GladiaApiEndpoint,
		CLIVersion:     cliVersion,
		httpClient:     &http.Client{},
		Verbose:        verbose,
	}
}

func (c *GladiaClient) apiURL(path string) string {
	return strings.TrimSuffix(c.GladiaEndpoint, "/") + path
}

func (c *GladiaClient) setRequestHeaders(req *http.Request) {
	req.Header.Set("x-gladia-key", c.ApiKey)
	req.Header.Set("x-gladia-version", "cli/"+c.CLIVersion)
}
