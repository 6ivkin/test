package reader

import "time"

type Reader struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
