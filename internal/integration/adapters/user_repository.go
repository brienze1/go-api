package adapters

import "github.com/brienze1/notes-api/internal/domain/usecases"

type UserRepository interface {
	usecases.GetUserUseCase
	usecases.ListUsersUseCase
	usecases.CreateUserUseCase
	usecases.DeleteUserUseCase
}
