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
	ID        string     `json:"id,omitempty"`
	PublicKey string     `json:"public_key"`
	Name      string     `json:"name"`
	ExpiresAt time.Time  `json:"expires_at"`
	Endpoints []Endpoint `json:"endpoints"`
}
type Endpoint struct {
	ID          string            `json:"id,omitempty"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Body        string            `json:"body"`
	ContentType string            `json:"content_type"`
	StatusCode  int               `json:"status_code"`
	DelayMS     int               `json:"delay_ms"`
	Headers     map[string]string `json:"headers,omitempty"`
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
func (r *Repository) Ready(ctx context.Context) error { return r.db.PingContext(ctx) }

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
	m.Endpoints = []Endpoint{e}
	return m, tx.Commit()
}

func insertEndpoint(ctx context.Context, tx *sql.Tx, mockID string, e Endpoint) error {
	headers, _ := json.Marshal(e.Headers)
	_, err := tx.ExecContext(ctx, `INSERT INTO mock_endpoints (mock_api_id, method, path, status_code, headers, body, content_type, delay_ms) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, mockID, e.Method, e.Path, e.StatusCode, headers, e.Body, e.ContentType, e.DelayMS)
	return err
}

func (r *Repository) FindByPublicKey(ctx context.Context, key string) (Mock, error) {
	var m Mock
	if err := r.db.QueryRowContext(ctx, `SELECT id,public_key,name,expires_at FROM mock_apis WHERE public_key=$1`, key).Scan(&m.ID, &m.PublicKey, &m.Name, &m.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Mock{}, ErrNotFound
		}
		return Mock{}, err
	}
	rows, err := r.db.QueryContext(ctx, `SELECT id,method,path,status_code,headers,body,content_type,delay_ms FROM mock_endpoints WHERE mock_api_id=$1 ORDER BY created_at`, m.ID)
	if err != nil {
		return Mock{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var e Endpoint
		var headers []byte
		if err := rows.Scan(&e.ID, &e.Method, &e.Path, &e.StatusCode, &headers, &e.Body, &e.ContentType, &e.DelayMS); err != nil {
			return Mock{}, err
		}
		_ = json.Unmarshal(headers, &e.Headers)
		m.Endpoints = append(m.Endpoints, e)
	}
	if err := rows.Err(); err != nil {
		return Mock{}, err
	}
	return m, nil
}

func (r *Repository) AddEndpoint(ctx context.Context, mockID string, e Endpoint) (Endpoint, error) {
	headers, _ := json.Marshal(e.Headers)
	err := r.db.QueryRowContext(ctx, `INSERT INTO mock_endpoints (mock_api_id,method,path,status_code,headers,body,content_type,delay_ms) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`, mockID, e.Method, e.Path, e.StatusCode, headers, e.Body, e.ContentType, e.DelayMS).Scan(&e.ID)
	return e, err
}

func (r *Repository) UpdateEndpoint(ctx context.Context, mockID, endpointID string, e Endpoint) (Endpoint, error) {
	headers, _ := json.Marshal(e.Headers)
	err := r.db.QueryRowContext(ctx, `UPDATE mock_endpoints SET method=$1,path=$2,status_code=$3,headers=$4,body=$5,content_type=$6,delay_ms=$7,updated_at=now() WHERE id=$8 AND mock_api_id=$9 RETURNING id`, e.Method, e.Path, e.StatusCode, headers, e.Body, e.ContentType, e.DelayMS, endpointID, mockID).Scan(&e.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return Endpoint{}, ErrNotFound
	}
	return e, err
}

func (r *Repository) DeleteEndpoint(ctx context.Context, mockID, endpointID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM mock_endpoints WHERE id=$1 AND mock_api_id=$2`, endpointID, mockID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) CountEndpoints(ctx context.Context, mockID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT count(*) FROM mock_endpoints WHERE mock_api_id=$1`, mockID).Scan(&count)
	return count, err
}

func (r *Repository) DeleteExpired(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM mock_apis WHERE expires_at <= now()`)
	return err
}

func (r *Repository) Authenticate(ctx context.Context, key, tokenHash string) (bool, error) {
	var found bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM mock_apis WHERE public_key=$1 AND management_token_hash=$2 AND expires_at > now())`, key, tokenHash).Scan(&found)
	return found, err
}
