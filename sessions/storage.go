package sessions

import "context"

type Storage interface {
	Get(ctx context.Context, id string) (*Data, error)
	Save(ctx context.Context, data *Data) error
	Delete(ctx context.Context, id string) error
}
