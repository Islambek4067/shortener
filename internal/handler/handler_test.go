package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shortener/internal/storage"
	"strings"
	"testing"
)

func newMux() (*http.ServeMux, *storage.Memory) {
	store := storage.NewMemory()
	h := NewHandler(store)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /shorten", h.Shorten)
	mux.HandleFunc("GET /{code}", h.Redirect)
	mux.HandleFunc("GET /stats/{code}", h.Stats)

	return mux, store
}

func TestShorten_ok(t *testing.T) {
	mux, _ := newMux()

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(`{"url":"https://go.dev"}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("waited 201, got %d", rec.Code)
	}

	var resp shortenResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)

	if err != nil {
		t.Fatalf("answer isn't json %v", err)
	}

	if resp.ShortURL == "" {
		t.Error("short_url empty")
	}
}

func TestShorten_BadInput(t *testing.T) {
	bodies := []string{
		`{"url":"не ссылка"}`,
		`{"url":"ftp://example.com"}`,
		`{"url":""}`,
		`это вообще не JSON`,
	}

	for _, body := range bodies {
		mux, _ := newMux()

		req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(body))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("тело %q: ожидали 400, получили %d", body, rec.Code)
		}
	}
}

func TestRedirect_ok(t *testing.T) {
	mux, store := newMux()
	store.Save("abc123", "https://go.dev")

	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("waited 302, getting  %d", rec.Code)
	}

	location := rec.Header().Get("Location")
	if location != "https://go.dev" {
		t.Errorf("waited Location https://go.dev gets %q", location)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	mux, _ := newMux()
	req := httptest.NewRequest(http.MethodGet, "/nonexistent", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("waited 404 gets %d", rec.Code)
	}
}

func TestStats_CountClicks(t *testing.T) {
	mux, store := newMux()
	store.Save("abc123", "https://go.dev")

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/abc123", nil)
		mux.ServeHTTP(httptest.NewRecorder(), req)
	}

	req := httptest.NewRequest(http.MethodGet, "/stats/abc123", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("waited 200 gets %d", rec.Code)
	}

	var resp statsResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("response isn't JSON: %v", err)
	}

	if resp.Clicks != 2 {
		t.Errorf("waited 2 click gets %d", resp.Clicks)
	}

	if resp.Original != "https://go.dev" {
		t.Errorf("incorrect Original: %q", resp.Original)
	}
}
