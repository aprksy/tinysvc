package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/aprksy/tinysvc/internal/domain"
	"github.com/aprksy/tinysvc/internal/repository"
)

type jsonRepository struct {
	db *sql.DB
}

// NewJSONRepository creates a new SQLite JSON repository
func NewJSONRepository(db *sql.DB) repository.JSONRepository {
	return &jsonRepository{db: db}
}

func (r *jsonRepository) Create(ctx context.Context, jsonBin *domain.JSONBin) error {
	query := `
        INSERT INTO json_bins (id, content, name, is_public, views, expires_at, created_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `

	var expiresAt interface{}
	if jsonBin.ExpiresAt != nil {
		expiresAt = jsonBin.ExpiresAt.Format(time.RFC3339)
	}

	_, err := r.db.ExecContext(ctx, query,
		jsonBin.ID,
		string(jsonBin.Content),
		jsonBin.Name,
		jsonBin.IsPublic,
		jsonBin.Views,
		expiresAt,
		jsonBin.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *jsonRepository) GetByID(ctx context.Context, id string) (*domain.JSONBin, error) {
	query := `
        SELECT id, content, name, is_public, views, expires_at, created_at
        FROM json_bins
        WHERE id = ?
    `

	jsonBin := &domain.JSONBin{}
	var expiresAt sql.NullString
	var createdAt string
	var content string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&jsonBin.ID,
		&content,
		&jsonBin.Name,
		&jsonBin.IsPublic,
		&jsonBin.Views,
		&expiresAt,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrJSONNotFound
	}
	if err != nil {
		return nil, err
	}

	// Convert content to RawMessage
	jsonBin.Content = json.RawMessage(content)

	// Parse created_at
	if parsedCreatedAt, err := time.Parse(time.RFC3339, createdAt); err == nil {
		jsonBin.CreatedAt = parsedCreatedAt
	}

	// Parse expires_at if exists
	if expiresAt.Valid {
		if parsedExpiresAt, err := time.Parse(time.RFC3339, expiresAt.String); err == nil {
			jsonBin.ExpiresAt = &parsedExpiresAt
		}
	}

	return jsonBin, nil
}

func (r *jsonRepository) Update(ctx context.Context, id string, content []byte) error {
	query := `UPDATE json_bins SET content = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, string(content), id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrJSONNotFound
	}

	return nil
}

func (r *jsonRepository) IncrementViews(ctx context.Context, id string) error {
	query := `UPDATE json_bins SET views = views + 1 WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *jsonRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM json_bins WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrJSONNotFound
	}

	return nil
}

func (r *jsonRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `
        DELETE FROM json_bins
        WHERE expires_at IS NOT NULL AND expires_at < ?
    `
	result, err := r.db.ExecContext(ctx, query, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
