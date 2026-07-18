package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"gopbincli/internal/model"
	"gopbincli/internal/repository"
	"gopbincli/internal/utilities"
)

type AuthHandler struct {
	Repo *repository.PostgresRepo
}

type BinCreateRequest struct {
	Content    string `json:"content"`
	Expiration string `json:"expiration"`
}

func (h *AuthHandler) CreateBinHandler(w http.ResponseWriter, r *http.Request) {
	var req BinCreateRequest
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if req.Content == "" {
		http.Error(w, "Content Empty", http.StatusBadRequest)
		return
	}
	id, err := utilities.GenerateRandomKey(8)
	if err != nil {
		http.Error(w, "Error Generating Url ", http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()
	var expiresAt time.Time
	switch req.Expiration {
	case "10m":
		expiresAt = now.Add(10 * time.Minute)
	case "1h":
		expiresAt = now.Add(time.Hour)
	case "24h":
		expiresAt = now.Add(24 * time.Hour)
	case "never", "":
		expiresAt = now.AddDate(100, 0, 0)
	default:
		http.Error(w, "Invalid expiration duration", http.StatusBadRequest)
		return
	}
	pasteBinModel := model.PasteBinItem{ID: id, Content: req.Content, CreatedAt: now, ExpiresAt: expiresAt}
	err = h.Repo.CreateBin(ctx, &pasteBinModel)
	if err != nil {
		// TODO(teacher): It's good practice to log the actual 'err' here (e.g. log.Printf) so you can debug DB failures on your server!
		http.Error(w, "Error Creating bin", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":         id,
		"expires_at": expiresAt.Format(time.RFC3339),
	})
}

func (h *AuthHandler) GetBinByIdHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Missing bin ID", http.StatusBadRequest)
		return
	}

	content, err := h.Repo.GetContentById(ctx, id)
	if err != nil {
		http.Error(w, "Bin not found or has expired", http.StatusNotFound)
		return
	}

	res := struct {
		Content string `json:"content"`
	}{Content: content}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
