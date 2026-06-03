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
	httpClient     *http.Client
	Verbose        bool
}

func NewGladiaClient(apiKey string, verbose bool) *GladiaClient {
	return &GladiaClient{
		ApiKey:         apiKey,
		GladiaEndpoint: GladiaApiEndpoint,
		httpClient:     &http.Client{},
		Verbose:        verbose,
	}
}

func (c *GladiaClient) apiURL(path string) string {
	return strings.TrimSuffix(c.GladiaEndpoint, "/") + path
}
