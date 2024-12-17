package persistence

type Config struct {
	// The DB URL.
	URL string `env:"URL"`
	// The ID of the row in the `url_ids(id)` table which holds the monotonically increasing number
	// used to generate the shortened URLs' paths.
	IDTableID uint64 `env:"ID_TABLE_ID"`
}
