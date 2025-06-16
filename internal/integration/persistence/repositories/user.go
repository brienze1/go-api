package repositories

import (
	"context"
	"errors"
	"fmt"
	"github.com/brienze1/notes-api/internal/domain/entities"
	"github.com/brienze1/notes-api/internal/domain/exceptions"
	"github.com/brienze1/notes-api/internal/integration/adapters"
	"github.com/brienze1/notes-api/internal/integration/persistence/repositories/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	connection *gorm.DB
}

func (u *userRepository) Get(ctx context.Context, userId string) (*entities.User, error) {
	var user *models.User
	err := u.connection.WithContext(ctx).Where(
		&entities.User{
			Id: userId,
		}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exceptions.NewNotFoundError(fmt.Sprintf("user with id %s not found", userId))
	}

	return user.ToEntity(), nil
}

func (u *userRepository) Create(ctx context.Context, user entities.User) (*entities.User, error) {
	var userModel models.User

	err := u.connection.WithContext(ctx).Create(userModel.FromEntity(user)).Error
	if err != nil {
		return nil, exceptions.NewInternalServerError("failed to create user", err.Error())
	}

	return userModel.ToEntity(), nil
}

func (u *userRepository) Delete(ctx context.Context, userId string) error {
	tx := u.connection.WithContext(ctx).Select(clause.Associations).Delete(&models.User{Id: userId})
	err := tx.Error
	if err != nil {
		return exceptions.NewInternalServerError(fmt.Sprintf("failed to delete user with id %s", userId), err.Error())
	}

	if tx.RowsAffected == 0 {
		return exceptions.NewNotFoundError(fmt.Sprintf("user with id %s not found", userId))
	}

	return nil
}

func (u *userRepository) List(ctx context.Context) ([]*entities.User, error) {
	var users models.Users
	err := u.connection.WithContext(ctx).Find(&users).Error
	if err != nil {
		return nil, exceptions.NewInternalServerError("failed to list users", err.Error())
	}

	return users.ToEntities(), nil
}

func NewUser(connection *gorm.DB) adapters.UserRepository {
	return &userRepository{connection: connection}
}
