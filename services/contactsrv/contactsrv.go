package contactsrv

import (
	"errors"

	"github.com/elsonigo/contact-app/domain"
	"github.com/elsonigo/contact-app/ports"
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

func (cs *ContactService) ValidateEmail(email string, id uuid.UUID) error {
	if email == "" {
		return errors.New("email required")
	}

	for _, c := range cs.repo.All() {
		if c.Email == email && c.ID != id {
			return errors.New("email already taken")
		}
	}

	return nil
}

func (cs *ContactService) Save(contact *domain.Contact) (*domain.Contact, error) {
	contact.ID = uuid.New()

	validationError := cs.ValidateEmail(contact.Email, contact.ID)

	if validationError != nil {
		contact.Errors = map[string]string{
			"email": validationError.Error(),
		}

		return contact, errors.New("invalid contact")
	}

	return cs.repo.Save(contact)
}

func (cs *ContactService) Update(contact *domain.Contact) (*domain.Contact, error) {
	validationError := cs.ValidateEmail(contact.Email, contact.ID)

	if validationError != nil {
		contact.Errors = map[string]string{
			"email": validationError.Error(),
		}

		return contact, errors.New("invalid contact")
	}

	return cs.repo.Update(contact)
}

func (cs *ContactService) All() []*domain.Contact {
	return cs.repo.All()
}

func (cs *ContactService) Search(s string) ([]*domain.Contact, error) {
	return cs.repo.Search(s)
}

func (cs *ContactService) Delete(contact *domain.Contact) error {
	return cs.repo.Delete(contact)
}

func (cs *ContactService) Find(id string) *domain.Contact {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil
	}

	return cs.repo.Find(parsed)
}
