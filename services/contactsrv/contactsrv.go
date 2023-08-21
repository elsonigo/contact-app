package contactsrv

import (
	"errors"

	"github.com/elsonigo/hypermediasystems/domain"
	"github.com/elsonigo/hypermediasystems/ports"
	"github.com/google/uuid"
)

type ContactService struct {
	repo ports.ContactRepo
}

func NewContactService(repo ports.ContactRepo) *ContactService {
	return &ContactService{
		repo: repo,
	}
}

func (cs *ContactService) validate(contact *domain.Contact) *domain.Contact {
	if contact.Email == "" {
		contact.Errors = map[string]string{
			"email": "Email required.",
		}

		return contact
	}

	all, err := cs.repo.All()
	if err != nil {
		return nil
	}

	for _, c := range all {
		if c.Email == contact.Email && c.ID != contact.ID {
			contact.Errors = map[string]string{
				"email": "Email already taken.",
			}
		}
	}

	return contact
}

func (cs *ContactService) Save(contact *domain.Contact) (*domain.Contact, error) {
	contact.ID = uuid.New()

	validated := cs.validate(contact)

	if validated.Errors != nil {
		return validated, errors.New("invalid contact")
	}

	return cs.repo.Save(contact)
}

func (cs *ContactService) Update(contact *domain.Contact) (*domain.Contact, error) {
	validated := cs.validate(contact)

	if validated.Errors != nil {
		return validated, errors.New("invalid contact")
	}

	return cs.repo.Update(contact)
}

func (cs *ContactService) All() ([]*domain.Contact, error) {
	return cs.repo.All()
}

func (cs *ContactService) Search(s string) ([]*domain.Contact, error) {
	return cs.repo.Search(s)
}

func (cs *ContactService) Delete(contact *domain.Contact) error {
	return cs.repo.Delete(contact)
}

func (cs *ContactService) Find(id string) (*domain.Contact, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id provided")
	}

	return cs.repo.Find(parsed)
}
