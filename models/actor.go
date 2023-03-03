package models

import (
	"time"
)

type Actor struct {
	Id        uint64
	Email     string
	Password  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
