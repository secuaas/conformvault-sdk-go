package conformvault

import (
	"fmt"
	"time"
)

// APIError represents an error returned by the ConformVault API.
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"error"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("conformvault: HTTP %d: %s", e.StatusCode, e.Message)
}

// RateLimitError is returned when the API returns 429 Too Many Requests.
type RateLimitError struct {
	APIError
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("conformvault: rate limited (retry after %s)", e.RetryAfter)
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == 404
	}
	return false
}

// IsRateLimited returns true if the error is a 429 Too Many Requests.
func IsRateLimited(err error) bool {
	_, ok := err.(*RateLimitError)
	return ok
}
