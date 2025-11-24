package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/aprksy/tinysvc/internal/domain"
	"github.com/aprksy/tinysvc/internal/repository"
)

type urlRepository struct {
	db *sql.DB
}

// NewURLRepository creates a new SQLite URL repository
func NewURLRepository(db *sql.DB) repository.URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(ctx context.Context, url *domain.ShortURL) error {
	query := `
        INSERT INTO urls (id, long_url, short_code, views, expires_at, created_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `

	var expiresAt interface{}
	if url.ExpiresAt != nil {
		expiresAt = url.ExpiresAt.Format(time.RFC3339)
	}

	_, err := r.db.ExecContext(ctx, query,
		url.ID,
		url.LongURL,
		url.ShortCode,
		url.Views,
		expiresAt,
		url.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *urlRepository) GetByCode(ctx context.Context, code string) (*domain.ShortURL, error) {
	query := `
        SELECT id, long_url, short_code, views, expires_at, created_at
        FROM urls
        WHERE short_code = ?
    `

	url := &domain.ShortURL{}
	var expiresAt sql.NullString
	var createdAt string

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&url.ID,
		&url.LongURL,
		&url.ShortCode,
		&url.Views,
		&expiresAt,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrURLNotFound
	}
	if err != nil {
		return nil, err
	}

	// Parse created_at
	if parsedCreatedAt, err := time.Parse(time.RFC3339, createdAt); err == nil {
		url.CreatedAt = parsedCreatedAt
	}

	// Parse expires_at if exists
	if expiresAt.Valid {
		if parsedExpiresAt, err := time.Parse(time.RFC3339, expiresAt.String); err == nil {
			url.ExpiresAt = &parsedExpiresAt
		}
	}

	return url, nil
}

func (r *urlRepository) GetByID(ctx context.Context, id string) (*domain.ShortURL, error) {
	query := `
        SELECT id, long_url, short_code, views, expires_at, created_at
        FROM urls
        WHERE id = ?
    `

	url := &domain.ShortURL{}
	var expiresAt sql.NullString
	var createdAt string

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&url.ID,
		&url.LongURL,
		&url.ShortCode,
		&url.Views,
		&expiresAt,
		&createdAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrURLNotFound
	}
	if err != nil {
		return nil, err
	}

	// Parse dates
	if parsedCreatedAt, err := time.Parse(time.RFC3339, createdAt); err == nil {
		url.CreatedAt = parsedCreatedAt
	}

	if expiresAt.Valid {
		if parsedExpiresAt, err := time.Parse(time.RFC3339, expiresAt.String); err == nil {
			url.ExpiresAt = &parsedExpiresAt
		}
	}

	return url, nil
}

func (r *urlRepository) IncrementViews(ctx context.Context, code string) error {
	query := `UPDATE urls SET views = views + 1 WHERE short_code = ?`
	_, err := r.db.ExecContext(ctx, query, code)
	return err
}

func (r *urlRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM urls WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return domain.ErrURLNotFound
	}

	return nil
}

func (r *urlRepository) DeleteExpired(ctx context.Context) (int64, error) {
	query := `
        DELETE FROM urls
        WHERE expires_at IS NOT NULL AND expires_at < ?
    `
	result, err := r.db.ExecContext(ctx, query, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *urlRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	query := `SELECT COUNT(*) FROM urls WHERE short_code = ?`
	var count int
	err := r.db.QueryRowContext(ctx, query, code).Scan(&count)
	return count > 0, err
}
