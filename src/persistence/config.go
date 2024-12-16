package persistence

type Config struct {
	URL       string `env:"URL"`
	IDTableID uint64 `env:"ID_TABLE_ID"`
}
