package engine

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

// Client represents the custom HTTP client for scanning
type Client struct {
	HTTPClient *http.Client
	Timeout    time.Duration
	Retries    int
}

// NewClient initializes a new Client with default or custom settings
func NewClient(timeout time.Duration, proxyURL string, insecure bool) (*Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}

	if proxyURL != "" {
		u, err := url.Parse(proxyURL)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(u)
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   timeout,
		// Explicitly disable following redirects to handle them manually if needed
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Client{
		HTTPClient: httpClient,
		Timeout:    timeout,
		Retries:    3, // Default retries
	}, nil
}

// Do performs an HTTP request with retry logic
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < c.Retries; i++ {
		resp, err = c.HTTPClient.Do(req)
		if err == nil {
			return resp, nil
		}
		// Wait before retry (exponential backoff could be added here)
		time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
	}

	return nil, err
}
