package entities

import (
	"time"
)

type Note struct {
	Id        string
	Title     string
	Content   string
	UserId    string
	User      User
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
