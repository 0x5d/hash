package core

import (
	"context"
)

type URLRepository interface {
	NextID(ctx context.Context) (uint64, error)
	Create(ctx context.Context, url ShortenedURL) error
	Update(ctx context.Context, id uint64, newURL *string, enabled *bool) error
	Get(ctx context.Context, id uint64) (*ShortenedURL, error)
}
