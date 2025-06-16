package adapters

import "context"

type ErrorHandler interface {
	HandlePanic(ctx context.Context, recovered any) (res []byte, statusCode int)
	HandleError(ctx context.Context, err error) (res []byte, statusCode int)
}
