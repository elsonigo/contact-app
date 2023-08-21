package services

import (
	"errors"

	"github.com/google/uuid"
)

type ContactService struct {
	repo ContactRepo
}

func NewContactService(repo ContactRepo) *ContactService {
	return &ContactService{
		repo: repo,
	}
}

func (cs *ContactService) validate(contact *Contact) *Contact {
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
		if c.Email == contact.Email {
			contact.Errors = map[string]string{
				"email": "Email already taken.",
			}
		}
	}

	return contact
}

func (cs *ContactService) Save(contact *Contact) (*Contact, error) {
	contact.ID = uuid.New()

	validated := cs.validate(contact)

	if validated.Errors != nil {
		return validated, errors.New("invalid contact")
	}

	return cs.repo.Save(contact)
}

func (cs *ContactService) Update(contact *Contact) (*Contact, error) {
	validated := cs.validate(contact)

	if validated.Errors != nil {
		return validated, errors.New("invalid contact")
	}

	return cs.repo.Update(contact)
}

func (cs *ContactService) All() ([]*Contact, error) {
	return cs.repo.All()
}

func (cs *ContactService) Search(s string) ([]*Contact, error) {
	return cs.repo.Search(s)
}

func (cs *ContactService) Delete(contact *Contact) error {
	return cs.repo.Delete(contact)
}

func (cs *ContactService) Find(id string) (*Contact, error) {
	return cs.repo.Find(id)
}
