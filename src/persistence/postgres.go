package persistence

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"strings"

	"github.com/0x5d/hash/core"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/tern/v2/migrate"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	c         Config
	log       *zap.Logger
	conn      *pgxpool.Pool
	idTableID uint64
}

//go:embed migrations
var migrations embed.FS

func NewPGURLRepo(log *zap.Logger, ctx context.Context, c Config) (*PostgresRepository, error) {
	conn, err := pgxpool.New(ctx, c.URL)
	if err != nil {
		return nil, err
	}
	return &PostgresRepository{log: log, c: c, conn: conn, idTableID: c.IDTableID}, nil
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
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &core.ErrNotFound{}
	}
	if err != nil {
		return nil, err
	}
	return &core.ShortenedURL{Original: url, Enabled: en}, nil
}

func (r *PostgresRepository) Get(ctx context.Context, id uint64) (*core.ShortenedURL, error) {
	q := `select url, enabled
	from urls
	where id = $1;`
	var (
		url     string
		enabled bool
	)
	err := r.conn.QueryRow(ctx, q, id).Scan(&url, &enabled)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &core.ErrNotFound{}
	}
	if err != nil {
		return nil, err
	}
	return &core.ShortenedURL{Original: url, Enabled: enabled}, nil
}

func (r *PostgresRepository) Migrate(ctx context.Context) error {
	conn, err := r.conn.Acquire(ctx)
	if err != nil {
		return err
	}
	m, err := migrate.NewMigrator(ctx, conn.Conn(), r.c.MigrationsTable)
	if err != nil {
		return err
	}
	err = fs.WalkDir(migrations, ".", func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.HasSuffix(path, ".sql") {
			str, rerr := migrations.ReadFile(path)
			if rerr != nil {
				return rerr
			}
			m.AppendMigration(path, string(str), "")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return m.Migrate(ctx)
}
