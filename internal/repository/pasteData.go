package repository

import (
	"context"
	"errors"
	"time"

	"gopbincli/internal/model"
)

var (
	errorBinExpired   = errors.New("pastebin expired bad luck")
	errorInvalidBinID = errors.New("No with that ID found")
)

func (r *PostgresRepo) CreateBin(ctx context.Context, paste *model.PasteBinItem) error {
	query := `INSERT INTO pastes (public_id, content, created_at, expires_at) values ($1,$2,$3,$4)`
	_, err := r.db.Exec(ctx, query, paste.ID, paste.Content, paste.CreatedAt, paste.ExpiresAt)
	return err
}

var ErrPasteNotFound = errors.New("paste not found or expired")

func (r *PostgresRepo) GetContentById(ctx context.Context, id string) (string, error) {
	var res string
	now := time.Now().UTC()
	query := `SELECT content from pastes where public_id = $1 and expires_at > $2`

	err := r.db.QueryRow(ctx, query, id, now).Scan(&res)
	if err != nil {
		return "", err
	}
	return res, err
}
