package domain

import "github.com/google/uuid"

type Contact struct {
	ID     uuid.UUID         `json:"id"`
	Email  string            `json:"email"`
	First  string            `json:"first,omitempty"`
	Last   string            `json:"last,omitempty"`
	Phone  string            `json:"phone,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

// page size for pagination
const PAGE_SIZE = 3
