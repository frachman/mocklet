package mocklet

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestValidateEndpoint(t *testing.T) {
	valid := &createRequest{Method: "get", Path: "/users/{id}", StatusCode: 200}
	if err := validateEndpoint(valid); err != nil || valid.Method != "GET" || valid.ContentType != "application/json" {
		t.Fatalf("valid endpoint rejected: %v", err)
	}
	for _, in := range []*createRequest{{Method: "TRACE", Path: "/x"}, {Method: "GET", Path: "users"}, {Method: "GET", Path: "/x", StatusCode: 700}, {Method: "GET", Path: "/x", DelayMS: 10001}, {Method: "GET", Path: "/x", ContentType: "text/html"}, {Method: "GET", Path: "/x", ContentType: "text/event-stream"}} {
		if validateEndpoint(in) == nil {
			t.Fatalf("expected rejection: %+v", in)
		}
	}
}

func TestMatchPathPattern(t *testing.T) {
	for _, path := range []string{"/products/123", "/products/abc"} {
		if !matchPathPattern("/products/{id}", path) {
			t.Errorf("expected template to match %s", path)
		}
	}
	for _, path := range []string{"/products", "/products/123/extra"} {
		if matchPathPattern("/products/{id}", path) {
			t.Errorf("expected template not to match %s", path)
		}
	}
	if !matchPathPattern("/health", "/health") || matchPathPattern("/health", "/healthz") {
		t.Fatal("literal path matching changed")
	}
}

func TestUsageBuffer(t *testing.T) {
	buffer := newUsageBuffer()
	for range 3 {
		if err := buffer.Add("runtime_requests"); err != nil {
			t.Fatal(err)
		}
	}
	if err := buffer.Add("unknown"); err == nil {
		t.Fatal("expected unknown usage event to be rejected")
	}
	counts := buffer.Take()
	if counts["runtime_requests"] != 3 || len(buffer.Take()) != 0 {
		t.Fatalf("unexpected buffered counts: %+v", counts)
	}
}

func TestPageViewRequiresSentinel(t *testing.T) {
	h := &Handler{repo: &Repository{usage: newUsageBuffer()}}
	for _, body := range []string{"", `{"source":"other"}`, `{"source":"landing"} {}`} {
		recorder := httptest.NewRecorder()
		h.pageView(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/telemetry/page-view", strings.NewReader(body)))
		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("body %q returned %d", body, recorder.Code)
		}
	}
	recorder := httptest.NewRecorder()
	h.pageView(recorder, httptest.NewRequest(http.MethodPost, "/api/v1/telemetry/page-view", strings.NewReader(`{"source":"landing"}`)))
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("valid sentinel returned %d", recorder.Code)
	}
}

func TestTokenHashIsStable(t *testing.T) {
	if hashToken("abc") != hashToken("abc") || hashToken("abc") == hashToken("def") {
		t.Fatal("token hashing is not stable/unique")
	}
}

func TestCORSPreflight(t *testing.T) {
	recorder := httptest.NewRecorder()
	cors(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Fatal("next handler should not run") })).ServeHTTP(recorder, httptest.NewRequest(http.MethodOptions, "/api/v1/mocks", nil))
	if recorder.Code != http.StatusNoContent || recorder.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("unexpected preflight response: %d %v", recorder.Code, recorder.Header())
	}
}

func TestRateLimiter(t *testing.T) {
	limiter := newRateLimiter(2, time.Hour)
	if !limiter.Allow("client") || !limiter.Allow("client") || limiter.Allow("client") {
		t.Fatal("expected third request to be limited")
	}
}
