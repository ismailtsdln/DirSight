package engine

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestScanner_Scan(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client, _ := NewClient(5*time.Second, "", false)
	scanner := NewScanner(client, 2)
	wordlist := []string{"test", "notfound"}

	ctx := context.Background()
	go scanner.Scan(ctx, server.URL, wordlist)

	foundTest := false
	for result := range scanner.Results {
		if result.URL == server.URL+"/test" && result.StatusCode == http.StatusOK {
			foundTest = true
		}
	}

	if !foundTest {
		t.Error("Scanner failed to find /test path")
	}
}
