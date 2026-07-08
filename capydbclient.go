// Package capydbclient is the shared HTTP transport for the CapyDB
// control-plane API clients - the CLI and the Terraform provider. It holds the
// retrying request engine, the APIError type, and the empty-slice normalization
// helper that both clients otherwise duplicated verbatim.
//
// Conventions mirrored from the control plane:
//   - All endpoints are JSON over HTTPS with bearer (org API key) auth.
//   - Idempotent GET requests are retried on network errors and 5xx responses;
//     non-GET requests are never retried.
//   - List endpoints may serialize empty Go slices as JSON null, so every list
//     decode should be passed through NormalizeList to get a non-nil slice.
package capydbclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIError is a non-2xx control-plane response.
type APIError struct {
	Message    string
	StatusCode int
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("capydb api request failed with status %d", e.StatusCode)
	}
	return fmt.Sprintf("capydb api request failed with status %d: %s", e.StatusCode, e.Message)
}

// IsNotFound reports whether err is an APIError with HTTP status 404.
func IsNotFound(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound
}

// Header is an extra request header applied to a single call - used for secrets
// that must not travel in the URL (e.g. the device-login poll token).
type Header struct {
	Key   string
	Value string
}

// Doer is the retrying HTTP transport shared by the control-plane clients. Each
// client wraps a *Doer with its own typed method surface.
type Doer struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
	// RetryBackoff holds the delay before each idempotent GET attempt. The
	// first entry is the delay before the initial attempt (zero), so the slice
	// length equals the total number of attempts. Non-GET requests are never
	// retried.
	RetryBackoff []time.Duration
}

// Do executes an API request, JSON-encoding payload when non-nil and decoding
// the response into dest when non-nil. Idempotent GET requests are retried on
// network errors and 5xx responses following RetryBackoff. extraHeaders are set
// on every attempt.
func (d *Doer) Do(ctx context.Context, method, path string, payload, dest any, extraHeaders ...Header) error {
	var encoded []byte
	if payload != nil {
		var err error
		encoded, err = json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode request body: %w", err)
		}
	}

	backoff := []time.Duration{0}
	if method == http.MethodGet && len(d.RetryBackoff) > 0 {
		backoff = d.RetryBackoff
	}

	var lastErr error
	for attempt, delay := range backoff {
		if delay > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		retryable, err := d.doOnce(ctx, method, path, encoded, payload != nil, dest, extraHeaders)
		if err == nil {
			return nil
		}
		lastErr = err
		if !retryable || ctx.Err() != nil || attempt == len(backoff)-1 {
			return err
		}
	}
	return lastErr
}

// doOnce performs a single HTTP round trip. The boolean return reports whether
// the failure is retryable (network error or 5xx response).
func (d *Doer) doOnce(ctx context.Context, method, path string, encoded []byte, hasBody bool, dest any, extraHeaders []Header) (bool, error) {
	var body io.Reader
	if hasBody {
		body = bytes.NewReader(encoded)
	}

	request, err := http.NewRequestWithContext(ctx, method, d.BaseURL+path, body)
	if err != nil {
		return false, fmt.Errorf("build request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", d.UserAgent)
	if hasBody {
		request.Header.Set("Content-Type", "application/json")
	}
	if d.APIKey != "" {
		request.Header.Set("Authorization", "Bearer "+d.APIKey)
	}
	for _, h := range extraHeaders {
		if h.Key != "" && h.Value != "" {
			request.Header.Set(h.Key, h.Value)
		}
	}

	response, err := d.HTTPClient.Do(request)
	if err != nil {
		return true, fmt.Errorf("perform request: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return true, fmt.Errorf("read response body: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var payload struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(raw, &payload)
		return response.StatusCode >= 500, &APIError{
			Message:    payload.Error,
			StatusCode: response.StatusCode,
		}
	}

	if dest == nil || len(raw) == 0 {
		return false, nil
	}
	if err := json.Unmarshal(raw, dest); err != nil {
		return false, fmt.Errorf("decode response body: %w", err)
	}
	return false, nil
}

// NormalizeList converts a nil slice (the control plane may serialize empty
// lists as JSON null) into an empty, non-nil slice.
func NormalizeList[T any](list []T) []T {
	if list == nil {
		return []T{}
	}
	return list
}
