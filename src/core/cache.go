package core

import "context"

type Cache interface {
	Set(ctx context.Context, key, url string) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}
