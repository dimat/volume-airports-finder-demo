package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestServer_ContextInterruption(t *testing.T) {
	addr := "localhost:12346"

	srv := New(addr)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := srv.Start(ctx)
	assert.NoError(t, err)
	elapsed := time.Since(start)

	assert.GreaterOrEqual(t, int64(elapsed), int64(500*time.Millisecond), "Server should wait for context cancellation")
}

func TestRegister(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Hello, world!"))
	})

	server := New("localhost:8080")
	server.Register("/test", handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	server.srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if rec.Body.String() != "Hello, world!" {
		t.Errorf("Expected response body 'Hello, world!', got '%s'", rec.Body.String())
	}
}
