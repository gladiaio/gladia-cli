package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestValidateResultURL(t *testing.T) {
	c := NewGladiaClient("k", false, "dev")
	c.GladiaEndpoint = "https://api.gladia.io"

	ok := []string{
		"https://api.gladia.io/v2/pre-recorded/abc",
		"https://API.Gladia.io/v2/pre-recorded/abc",
		"https://api.gladia.io:443/v2/pre-recorded/abc",
	}
	for _, u := range ok {
		if err := c.validateResultURL(u); err != nil {
			t.Errorf("validateResultURL(%q) = %v, want nil", u, err)
		}
	}

	bad := []string{
		"",
		"/v2/pre-recorded/abc",
		"file:///etc/passwd",
		"ftp://api.gladia.io/x",
		"http://api.gladia.io/v2/pre-recorded/abc",   // scheme downgrade
		"https://api.gladia.io:8443/v2/pre-recorded", // other port
		"https://evil.example/v2/pre-recorded/abc",
		"https://api.gladia.io.evil.example/x",
		"https://api.gladia.io@evil.example/x",
		"https://user:pass@api.gladia.io/x",
		"http://169.254.169.254/latest/meta-data/",
		"://bad",
	}
	for _, u := range bad {
		if err := c.validateResultURL(u); err == nil {
			t.Errorf("validateResultURL(%q) = nil, want error", u)
		}
	}
}

func TestSetRequestHeaders_keyOnlyForGladiaOrigin(t *testing.T) {
	c := NewGladiaClient("secret", false, "dev")
	c.GladiaEndpoint = "https://api.gladia.io"

	for target, want := range map[string]string{
		"https://api.gladia.io/v2/upload/": "secret",
		"https://other.example/x":          "",
		"http://api.gladia.io/x":           "",
	} {
		req := httptest.NewRequest(http.MethodGet, target, nil)
		c.setRequestHeaders(req)
		if got := req.Header.Get("x-gladia-key"); got != want {
			t.Errorf("%s: x-gladia-key = %q, want %q", target, got, want)
		}
		if got := req.Header.Get("x-gladia-version"); got != "cli/dev" {
			t.Errorf("%s: x-gladia-version = %q", target, got)
		}
	}
}

func TestTranscribeAudioURL_rejectsForeignResultURL(t *testing.T) {
	var foreignHits int32
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&foreignHits, 1)
		_, _ = w.Write([]byte(`{"status":"done"}`))
	}))
	defer foreign.Close()

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(TranscriptionResponse{ID: "t1", ResultURL: foreign.URL + "/result"})
	}))
	defer api.Close()

	c := NewGladiaClient("secret", false, "dev")
	c.GladiaEndpoint = api.URL
	_, err := c.TranscribeAudioURL("https://audio.example/x.wav", TranscriptionRequest{})
	if err == nil || !strings.Contains(err.Error(), "outside the Gladia API origin") {
		t.Fatalf("err = %v, want origin error", err)
	}
	if n := atomic.LoadInt32(&foreignHits); n != 0 {
		t.Fatalf("foreign host received %d requests, want 0", n)
	}
}

func TestPoll_refusesCrossOriginRedirect(t *testing.T) {
	var leaked atomic.Value
	leaked.Store("")
	foreign := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		leaked.Store(r.Header.Get("x-gladia-key") + "|hit")
		_, _ = w.Write([]byte(`{"status":"done"}`))
	}))
	defer foreign.Close()

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, foreign.URL+"/result", http.StatusFound)
	}))
	defer api.Close()

	c := NewGladiaClient("secret", false, "dev")
	c.GladiaEndpoint = api.URL
	_, err := c.pollForTranscriptionResult(api.URL + "/v2/pre-recorded/t1")
	if err == nil || !strings.Contains(err.Error(), "refusing redirect") {
		t.Fatalf("err = %v, want redirect refusal", err)
	}
	if got := leaked.Load().(string); got != "" {
		t.Fatalf("foreign host was contacted (%q)", got)
	}
}

func TestPoll_followsSameOriginRedirect(t *testing.T) {
	api := httptest.NewServer(nil)
	defer api.Close()
	api.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/old" {
			http.Redirect(w, r, "/new", http.StatusFound)
			return
		}
		if r.Header.Get("x-gladia-key") != "secret" {
			t.Errorf("missing api key after same-origin redirect")
		}
		_, _ = w.Write([]byte(`{"status":"done","result":{"transcription":{"full_transcript":"ok"}}}`))
	})

	c := NewGladiaClient("secret", false, "dev")
	c.GladiaEndpoint = api.URL
	res, err := c.pollForTranscriptionResult(api.URL + "/old")
	if err != nil {
		t.Fatal(err)
	}
	if res.Result.Transcription.FullTranscript != "ok" {
		t.Fatalf("got %+v", res)
	}
}

func TestCheckRedirect_limit(t *testing.T) {
	c := NewGladiaClient("k", false, "dev")
	c.GladiaEndpoint = "https://api.gladia.io"
	u, _ := url.Parse("https://api.gladia.io/x")
	req := &http.Request{URL: u}
	if err := c.checkRedirect(req, make([]*http.Request, maxRedirects)); err == nil {
		t.Fatal("expected redirect limit error")
	}
	if err := c.checkRedirect(req, nil); err != nil {
		t.Fatalf("same-origin redirect refused: %v", err)
	}
}

func TestPoll_stopsAfterTimeout(t *testing.T) {
	var hits int32
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		_, _ = w.Write([]byte(`{"status":"processing"}`))
	}))
	defer api.Close()

	c := NewGladiaClient("k", false, "dev")
	c.GladiaEndpoint = api.URL
	c.PollInterval = 10 * time.Millisecond
	c.PollTimeout = 200 * time.Millisecond

	start := time.Now()
	_, err := c.pollForTranscriptionResult(api.URL + "/v2/pre-recorded/t1")
	if err == nil || !strings.Contains(err.Error(), "did not complete within") {
		t.Fatalf("err = %v, want timeout error", err)
	}
	if elapsed := time.Since(start); elapsed > 5*time.Second {
		t.Fatalf("polling took %s, deadline not enforced", elapsed)
	}
	if atomic.LoadInt32(&hits) == 0 {
		t.Fatal("server never polled")
	}
}

func TestPoll_retriesTransientErrorsThenFailsOnClientError(t *testing.T) {
	var hits int32
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch atomic.AddInt32(&hits, 1) {
		case 1:
			w.WriteHeader(http.StatusServiceUnavailable)
		case 2:
			w.WriteHeader(http.StatusTooManyRequests)
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"message":"not found"}`))
		}
	}))
	defer api.Close()

	c := NewGladiaClient("k", false, "dev")
	c.GladiaEndpoint = api.URL
	c.PollInterval = time.Millisecond
	_, err := c.pollForTranscriptionResult(api.URL + "/v2/pre-recorded/t1")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("err = %v, want not found", err)
	}
	if n := atomic.LoadInt32(&hits); n != 3 {
		t.Fatalf("hits = %d, want 3", n)
	}
}
