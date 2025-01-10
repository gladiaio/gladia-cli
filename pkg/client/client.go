package client

import "net/http"

// Can be override by the developer, before initializing the client.
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
