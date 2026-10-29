package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"shortener/internal/storage"
)

type Handler struct {
	store *storage.Memory
}

func NewHandler(store *storage.Memory) *Handler {
	return &Handler{
		store: store,
	}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "not correct JSON", http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(req.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		http.Error(w, "not correct URL", http.StatusBadRequest)
		return
	}

	code := storage.GenerateCode(6)
	h.store.Save(code, req.URL)

	w.Header().Set("Contenet-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shortenResponse{
		ShortURL: "http://" + r.Host + "/" + code,
	})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	link, ok := h.store.Get(code)

	if ok {
		http.Redirect(w, r, link, http.StatusFound)

	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}
