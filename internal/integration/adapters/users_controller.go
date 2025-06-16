package adapters

import (
	"context"
)

type UsersController interface {
	ListUsers(ctx context.Context) (response []byte, status int)
	GetUser(ctx context.Context, id string) (response []byte, status int)
	CreateUser(ctx context.Context, body []byte) (response []byte, status int)
	DeleteUser(ctx context.Context, id string) (response []byte, status int)
}
