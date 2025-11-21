package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/aprksy/tinysvc/internal/domain"
	"github.com/aprksy/tinysvc/internal/repository"
	_ "github.com/mattn/go-sqlite3"
)

type pasteRepository struct {
	db *sql.DB
}

// NewPasteRepository creates a new SQLite paste repository
func NewPasteRepository(db *sql.DB) repository.PasteRepository {
	return &pasteRepository{db: db}
}

func (r *pasteRepository) Create(ctx context.Context, paste *domain.Paste) error {
	query := `
        INSERT INTO pastes (id, content, is_markdown, expires_at, created_at)
        VALUES (?, ?, ?, ?, ?)
    `

	var expiresAt interface{}
	if paste.ExpiresAt != nil {
		expiresAt = paste.ExpiresAt.Format(time.RFC3339)
	}

	_, err := r.db.ExecContext(ctx, query,
		paste.ID,
		paste.Content,
		paste.IsMarkdown,
		expiresAt,
		paste.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *pasteRepository) GetByID(ctx context.Context, id string) (*domain.Paste, error) {
	query := `
        SELECT id, content, is_markdown, expires_at, created_at
        FROM pastes
        WHERE id = ?
    `

	paste := &domain.Paste{}
	var expiresAt sql.NullString
	var createdAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&paste.ID,
		&paste.Content,
		&paste.IsMarkdown,
		&expiresAt,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrPasteNotFound
	}
	if err != nil {
		return nil, err
	}

	// Parse created_at
	if parsedCreatedAt, err := time.Parse(time.RFC3339, createdAt); err == nil {
		paste.CreatedAt = parsedCreatedAt
	}

	// Parse expires_at if exists
	if expiresAt.Valid {
		if parsedExpiresAt, err := time.Parse(time.RFC3339, expiresAt.String); err == nil {
			paste.ExpiresAt = &parsedExpiresAt
		}
	}

	return paste, nil
}

func (r *pasteRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM pastes WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrPasteNotFound
	}

	return nil
}

func (r *pasteRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `
        DELETE FROM pastes
        WHERE expires_at IS NOT NULL AND expires_at < ?
    `
	result, err := r.db.ExecContext(ctx, query, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
