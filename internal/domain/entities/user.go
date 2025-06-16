package entities

import (
	"time"
)

type User struct {
	Id        string
	Name      string
	Note      []*Note
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
