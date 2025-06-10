package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_healthHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	api := &API{
		log: slog.Default(),
	}
	api.healthHandler(w, r)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	body, err := io.ReadAll(w.Body)
	require.NoError(t, err, "Failed to read response body")

	if string(body) != "OK\n" {
		t.Errorf("Expected body 'OK', got '%s'", w.Body.String())
	}
}
