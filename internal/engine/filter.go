package engine

import "net/http"

// Filter defines criteria for excluding or highlighting results
type Filter struct {
	ExcludeStatusCodes []int
	ExcludeLengths     []int64
}

// ShouldShow returns true if the result should be displayed based on filters
func (f *Filter) ShouldShow(result Result) bool {
	for _, code := range f.ExcludeStatusCodes {
		if result.StatusCode == code {
			return false
		}
	}
	for _, length := range f.ExcludeLengths {
		if result.Length == length {
			return false
		}
	}
	return true
}

// DetectBypass check if the result indicates a possible bypass (e.g., 200 instead of 403)
func DetectBypass(originalStatus, bypassStatus int) bool {
	return originalStatus == http.StatusForbidden && bypassStatus == http.StatusOK
}
