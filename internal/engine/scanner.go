package engine

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Result represents the outcome of a single scan request
type Result struct {
	URL        string
	StatusCode int
	Length     int64
	Method     string
}

// Scanner handles the concurrent scanning process
type Scanner struct {
	Client  *Client
	Threads int
	Results chan Result
}

// NewScanner initializes a new Scanner
func NewScanner(client *Client, threads int) *Scanner {
	return &Scanner{
		Client:  client,
		Threads: threads,
		Results: make(chan Result),
	}
}

// Scan performs a concurrent scan of the given target and wordlist
func (s *Scanner) Scan(ctx context.Context, target string, wordlist []string) {
	var wg sync.WaitGroup
	jobs := make(chan string, s.Threads)

	// Start worker pool
	for i := 0; i < s.Threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					s.processPath(target, path)
				}
			}
		}()
	}

	// Send paths to workers
	go func() {
		for _, path := range wordlist {
			jobs <- path
		}
		close(jobs)
	}()

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(s.Results)
	}()
}

func (s *Scanner) processPath(target, path string) {
	fullURL := fmt.Sprintf("%s/%s", strings.TrimRight(target, "/"), strings.TrimLeft(path, "/"))
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	s.Results <- Result{
		URL:        fullURL,
		StatusCode: resp.StatusCode,
		Length:     resp.ContentLength,
		Method:     req.Method,
	}
}
