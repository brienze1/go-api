package models

import (
	"time"

	"github.com/brienze1/notes-api/internal/domain/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id        string    `gorm:"column:id;primaryKey" json:"id"`
	Name      string    `gorm:"column:name;not null;size:255" json:"name" validate:"required,min=3"`
	Note      []*Note   `gorm:"foreignKey:UserId;references:Id" json:"-"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(db *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	return
}

func (u User) ToEntity() *entities.User {
	return &entities.User{
		Id:        u.Id,
		Name:      u.Name,
		CreatedAt: &u.CreatedAt,
		UpdatedAt: &u.UpdatedAt,
	}
}

func (u *User) FromEntity(user entities.User) *User {
	u.Id = user.Id
	u.Name = user.Name

	if user.CreatedAt != nil {
		u.CreatedAt = *user.CreatedAt
	}
	if user.UpdatedAt != nil {
		u.UpdatedAt = *user.UpdatedAt
	}
	return u
}

type Users []*User

func (u Users) ToEntities() []*entities.User {
	users := make([]*entities.User, len(u))
	for i, user := range u {
		users[i] = user.ToEntity()
	}
	return users
}
