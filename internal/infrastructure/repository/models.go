package repository

import (
	"time"
)

type UserModel struct {
	id        int64
	name      string
	email     *string
	phone     string
	password  string
	createdAt time.Time
}
