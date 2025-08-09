package health

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestServiceCheckLocalstackOK(t *testing.T) {
	// Mock HTTP client via http.DefaultTransport replacement using httptest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	svc := NewService(server.URL)
	resp := svc.Check(context.Background())

	if !resp.OK {
		t.Fatalf("expected OK true, got false")
	}
	if resp.LocalStack == "" {
		t.Fatalf("expected LocalStack endpoint to be set")
	}
}

func TestServiceCheckLocalstackDown(t *testing.T) {
	svc := NewService("http://127.0.0.1:0") // invalid port
	resp := svc.Check(context.Background())
	if resp.OK {
		t.Fatalf("expected OK false when localstack unreachable")
	}
}
