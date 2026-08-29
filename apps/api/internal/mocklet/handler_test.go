package mocklet

import "testing"

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
