package mocklet

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidateEndpoint(t *testing.T) {
	valid := &createRequest{Method: "get", Path: "/users/{id}", StatusCode: 200}
	if err := validateEndpoint(valid); err != nil || valid.Method != "GET" || valid.ContentType != "application/json" {
		t.Fatalf("valid endpoint rejected: %v", err)
	}
	for _, in := range []*createRequest{{Method: "TRACE", Path: "/x"}, {Method: "GET", Path: "users"}, {Method: "GET", Path: "/x", StatusCode: 700}, {Method: "GET", Path: "/x", DelayMS: 10001}} {
		if validateEndpoint(in) == nil {
			t.Fatalf("expected rejection: %+v", in)
		}
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
