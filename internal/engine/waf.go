package engine

import (
	"net/http"
	"strings"
)

// WAFType represents the detected WAF
type WAFType string

const (
	WAFNone        WAFType = "None"
	WAFCloudflare  WAFType = "Cloudflare"
	WAFAkamai      WAFType = "Akamai"
	WAFModSecurity WAFType = "ModSecurity"
	WAFAWS         WAFType = "AWS WAF"
)

// DetectWAF attempts to identify a WAF based on response headers
func DetectWAF(resp *http.Response) WAFType {
	server := strings.ToLower(resp.Header.Get("Server"))
	if strings.Contains(server, "cloudflare") {
		return WAFCloudflare
	}
	if resp.Header.Get("CF-RAY") != "" {
		return WAFCloudflare
	}
	if strings.Contains(server, "akamai") || resp.Header.Get("X-Akamai-Transformed") != "" {
		return WAFAkamai
	}
	if resp.Header.Get("X-AWS-WAF-Attributes") != "" {
		return WAFAWS
	}
	// Simple check for ModSecurity or others often found in headers
	for k := range resp.Header {
		if strings.Contains(strings.ToLower(k), "mod_security") || strings.Contains(strings.ToLower(k), "waf") {
			return WAFModSecurity
		}
	}
	return WAFNone
}
