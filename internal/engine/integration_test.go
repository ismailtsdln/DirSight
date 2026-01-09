package engine

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ismailtsdln/DirSight/internal/bypass"
)

func TestIntegration_FullScanWithBypass(t *testing.T) {
	// Mock server that returns 403 normally but 200 with a bypass header
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-For") == "127.0.0.1" && r.URL.Path == "/forbidden" {
			w.WriteHeader(http.StatusOK)
		} else if r.URL.Path == "/forbidden" {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, _ := NewClient(5*time.Second, "", false)
	scanner := NewScanner(client, 2)

	// We'll test with the bypass logic
	bypassMethods := bypass.GetBypassMethods()

	ctx := context.Background()

	// Simulate one scan without bypass
	req1, _ := http.NewRequest("GET", server.URL+"/forbidden", nil)
	resp1, _ := client.Do(req1)
	if resp1.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403, got %d", resp1.StatusCode)
	}

	// Simulate scan with bypass application
	foundBypass := false
	for _, m := range bypassMethods {
		req, _ := http.NewRequest("GET", server.URL+"/forbidden", nil)
		m.Apply(req)
		resp, _ := client.Do(req)
		if resp.StatusCode == http.StatusOK {
			foundBypass = true
			break
		}
	}

	if !foundBypass {
		t.Error("Bypass integration test failed: could not bypass 403")
	}
}
