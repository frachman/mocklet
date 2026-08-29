package mocklet

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrNotFound = errors.New("mock not found")

type Mock struct {
	ID, PublicKey, Name string
	ExpiresAt           time.Time
	Endpoint            Endpoint
}
type Endpoint struct {
	ID, Method, Path, Body, ContentType string
	StatusCode, DelayMS                 int
	Headers                             map[string]string
}

type Repository struct{ db *sql.DB }

func NewPostgresRepository(ctx context.Context, url string) (*Repository, error) {
	db, err := sql.Open("pgx", url)
	if err != nil {
		return nil, err
	}
	if err = db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return &Repository{db: db}, nil
}
func (r *Repository) Close() error { return r.db.Close() }

func (r *Repository) Create(ctx context.Context, name, tokenHash string, expiresAt time.Time, e Endpoint) (Mock, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Mock{}, err
	}
	defer tx.Rollback()
	var m Mock
	if err = tx.QueryRowContext(ctx, `INSERT INTO mock_apis (public_key, management_token_hash, name, expires_at) VALUES (encode(gen_random_bytes(12),'hex'), $1, $2, $3) RETURNING id, public_key, name, expires_at`, tokenHash, name, expiresAt).Scan(&m.ID, &m.PublicKey, &m.Name, &m.ExpiresAt); err != nil {
		return Mock{}, err
	}
	if err = insertEndpoint(ctx, tx, m.ID, e); err != nil {
		return Mock{}, err
	}
	if err = tx.QueryRowContext(ctx, `SELECT id FROM mock_endpoints WHERE mock_api_id=$1`, m.ID).Scan(&e.ID); err != nil {
		return Mock{}, err
	}
	m.Endpoint = e
	return m, tx.Commit()
}

func insertEndpoint(ctx context.Context, tx *sql.Tx, mockID string, e Endpoint) error {
	headers, _ := json.Marshal(e.Headers)
	_, err := tx.ExecContext(ctx, `INSERT INTO mock_endpoints (mock_api_id, method, path, status_code, headers, body, content_type, delay_ms) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, mockID, e.Method, e.Path, e.StatusCode, headers, e.Body, e.ContentType, e.DelayMS)
	return err
}

func (r *Repository) FindByPublicKey(ctx context.Context, key string) (Mock, error) {
	var m Mock
	var headers []byte
	row := r.db.QueryRowContext(ctx, `SELECT a.id,a.public_key,a.name,a.expires_at,e.id,e.method,e.path,e.status_code,e.headers,e.body,e.content_type,e.delay_ms FROM mock_apis a JOIN mock_endpoints e ON e.mock_api_id=a.id WHERE a.public_key=$1 ORDER BY e.created_at LIMIT 1`, key)
	if err := row.Scan(&m.ID, &m.PublicKey, &m.Name, &m.ExpiresAt, &m.Endpoint.ID, &m.Endpoint.Method, &m.Endpoint.Path, &m.Endpoint.StatusCode, &headers, &m.Endpoint.Body, &m.Endpoint.ContentType, &m.Endpoint.DelayMS); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Mock{}, ErrNotFound
		}
		return Mock{}, err
	}
	_ = json.Unmarshal(headers, &m.Endpoint.Headers)
	return m, nil
}

func (r *Repository) Authenticate(ctx context.Context, key, tokenHash string) (bool, error) {
	var found bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM mock_apis WHERE public_key=$1 AND management_token_hash=$2 AND expires_at > now())`, key, tokenHash).Scan(&found)
	return found, err
}
