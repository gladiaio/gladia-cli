package client

import "testing"

func TestGladiaClient_apiURL(t *testing.T) {
	c := &GladiaClient{GladiaEndpoint: "https://api.gladia.io"}
	if got := c.apiURL("/v2/upload/"); got != "https://api.gladia.io/v2/upload/" {
		t.Fatalf("got %q", got)
	}

	c.GladiaEndpoint = "https://api.gladia.io/"
	if got := c.apiURL("/v2/pre-recorded"); got != "https://api.gladia.io/v2/pre-recorded" {
		t.Fatalf("got %q", got)
	}
}

func TestNewGladiaClient_defaults(t *testing.T) {
	c := NewGladiaClient("key", true)
	if c.ApiKey != "key" || !c.Verbose || c.GladiaEndpoint != GladiaApiEndpoint {
		t.Fatalf("unexpected client: %+v", c)
	}
}
