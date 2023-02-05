package data

import "time"

type UserSession struct {
	ID        string
	Token     string
	CreatedAt time.Time
	UserAgent string
	IP        string
}
