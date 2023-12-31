package ports

import (
	"github.com/elsonigo/contact-app/domain"
	"github.com/google/uuid"
)

type ContactRepo interface {
	Save(*domain.Contact) (*domain.Contact, error)
	Update(*domain.Contact) (*domain.Contact, error)
	All() []*domain.Contact
	Page(int) []*domain.Contact
	Search(string) ([]*domain.Contact, error)
	Delete(*domain.Contact) error
	Find(uuid.UUID) *domain.Contact
}
