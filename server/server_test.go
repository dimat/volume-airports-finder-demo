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

func TestServerStartError(t *testing.T) {
	// Start the first server
	addr := "localhost:8081"
	server1 := New(addr)
	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	go func() {
		if err := server1.Start(ctx1); err != nil {
			t.Logf("Expected server1 to shut down gracefully, got: %v", err)
		}
	}()

	// Ensure server1 is running
	time.Sleep(500 * time.Millisecond)

	// Start the second server on the same port
	server2 := New(addr)
	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	errCh := make(chan error)
	go func() {
		errCh <- server2.Start(ctx2)
	}()

	// Wait for server2 to return an error
	select {
	case err := <-errCh:
		if err == nil {
			t.Error("Expected server2 to return an error due to the port being in use")
		} else {
			t.Logf("server2 returned an error as expected: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Error("server2 did not return an error within the expected time")
	}
}
