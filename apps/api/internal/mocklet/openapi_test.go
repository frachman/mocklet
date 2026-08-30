package mocklet

import (
	"strings"
	"testing"
)

func TestPreviewOpenAPIUsesExamplesAndDeterministicFallback(t *testing.T) {
	document := []byte(`{"openapi":"3.0.3","info":{"title":"Checkout","version":"1.0.0"},"paths":{"/orders/{id}":{"get":{"parameters":[{"name":"id","in":"path","required":true,"schema":{"type":"string"}}],"responses":{"200":{"description":"ok","content":{"application/json":{"example":{"id":"order-1","status":"ready"}}}},"404":{"description":"missing","content":{"application/json":{"schema":{"type":"object","properties":{"error":{"type":"string"}}}}}}}}}}}`)
	preview, err := previewOpenAPI(document)
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.Endpoints) != 1 || len(preview.Endpoints[0].Scenarios) != 2 {
		t.Fatalf("unexpected preview: %+v", preview)
	}
	if preview.Endpoints[0].Path != "/orders/{id}" || preview.Endpoints[0].StatusCode != 200 {
		t.Fatalf("unexpected endpoint: %+v", preview.Endpoints[0])
	}
	if !strings.Contains(preview.Endpoints[0].Body, "order-1") {
		t.Fatalf("explicit example was not used: %s", preview.Endpoints[0].Body)
	}
	if preview.Endpoints[0].Scenarios[1].StatusCode != 404 || !strings.Contains(preview.Endpoints[0].Scenarios[1].Body, "error") {
		t.Fatalf("schema fallback missing: %+v", preview.Endpoints[0].Scenarios[1])
	}
}

func TestPreviewOpenAPIBoundsRoutes(t *testing.T) {
	paths := make([]string, 0, 6)
	for i := 0; i < 6; i++ {
		paths = append(paths, `"/route`+string(rune('a'+i))+`":{"get":{"responses":{"200":{"description":"ok"}}}}`)
	}
	document := []byte(`{"openapi":"3.0.3","info":{"title":"Too large","version":"1"},"paths":{` + strings.Join(paths, ",") + `}}`)
	if _, err := previewOpenAPI(document); err == nil {
		t.Fatal("expected endpoint limit rejection")
	}
}
