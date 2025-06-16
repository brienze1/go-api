package adapters

import "context"

type Secret interface {
	Get(ctx context.Context, key string) ([]byte, error)
}
