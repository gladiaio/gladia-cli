package client

// Can be override by the developer, before initializing the client.
var GladiaApiEndpoint = "https://api.gladia.io/v2"

type GladiaClient struct {
	ApiKey         string
	GladiaEndpoint string
}

func NewGladiaClient(apiKey string) *GladiaClient {
	return &GladiaClient{
		ApiKey:         apiKey,
		GladiaEndpoint: GladiaApiEndpoint,
	}
}
