package persistence

import (
	"context"

	"github.com/0x5d/hash/core"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	conn      *pgxpool.Pool
	idTableID uint64
}

func NewPGURLRepo(ctx context.Context, c Config) (*PostgresRepository, error) {
	conn, err := pgxpool.New(ctx, c.URL)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{conn: conn, idTableID: c.IDTableID}, nil
}

func (r *PostgresRepository) NextID(ctx context.Context) (uint64, error) {
	q := `update url_ids 
   	set url_id = url_id + 1
	where id = $1
	returning url_id;`
	row := r.conn.QueryRow(ctx, q, r.idTableID)
	var id uint64
	return id, row.Scan(&id)
}

func (r *PostgresRepository) Create(ctx context.Context, url core.ShortenedURL) error {
	q := `insert into urls(url, short_key, enabled) values ($1, $2, $3);`
	_, err := r.conn.Exec(ctx, q, url.Original, url.ShortKey, url.Enabled)
	return err
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, newURL *string, enabled *bool) (*core.ShortenedURL, error) {
	q := `update urls set
    url = coalesce($1::text, url),
    enabled = coalesce($2::boolean, enabled)
	where id = $3
	returning url, enabled;`
	var (
		url string
		en  bool
	)
	err := r.conn.QueryRow(ctx, q, newURL, enabled, id).Scan(&url, &en)
	if err != nil {
		return nil, err
	}
	return &core.ShortenedURL{Original: url, Enabled: en}, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id uint64) (*core.ShortenedURL, error) {
	q := `select url, enabled
	from urls
	where id = $1;`
	row := r.conn.QueryRow(ctx, q, id)
	var (
		url     string
		enabled bool
	)
	err := row.Scan(&url, &enabled)
	if err != nil {
		return nil, err
	}
	return &core.ShortenedURL{Original: url, Enabled: enabled}, nil
}
