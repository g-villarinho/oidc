package ports

import "context"

type Hasher interface {
	Hash(ctx context.Context, plaintext string) (string, error)
	Compare(ctx context.Context, plaintext string, hashed string) error
}
