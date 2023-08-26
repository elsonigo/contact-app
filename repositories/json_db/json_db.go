package json_db

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/elsonigo/contact-app/domain"
	"github.com/elsonigo/contact-app/ports"
	"github.com/google/uuid"
)

type JsonDatabase struct {
	contacts []*domain.Contact
}

var _ ports.ContactRepo = &JsonDatabase{}

func OpenJsonDatabase() (*JsonDatabase, error) {
	data, err := os.ReadFile("db.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return &JsonDatabase{}, nil
		}

		return nil, err
	}

	contacts := []*domain.Contact{}
	err = json.Unmarshal(data, &contacts)
	if err != nil {
		return nil, err
	}

	return &JsonDatabase{
		contacts: contacts,
	}, nil
}

func (db *JsonDatabase) Save(contact *domain.Contact) (*domain.Contact, error) {
	db.contacts = append(db.contacts, contact)
	err := saveToFile(db.contacts)
	if err != nil {
		return contact, err
	}

	return contact, nil
}

func saveToFile(contacts []*domain.Contact) error {
	marshalledJson, _ := json.MarshalIndent(contacts, "", "  ")
	err := os.WriteFile("db.json", marshalledJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *JsonDatabase) All() []*domain.Contact {
	if db.contacts == nil {
		return nil
	}

	return db.contacts
}

func (db *JsonDatabase) Page(p int) []*domain.Contact {
	if db.contacts == nil {
		return nil
	}

	start := (p - 1) * domain.PAGE_SIZE
	end := start + domain.PAGE_SIZE
	max := len(db.contacts)

	if end > max {
		end = max
	}

	return db.contacts[start:end]
}

func (db *JsonDatabase) Search(q string) ([]*domain.Contact, error) {
	if q == "" {
		return nil, errors.New("no query string given")
	}

	results := []*domain.Contact{}
	query := strings.ToLower(q)

	for _, contact := range db.contacts {
		if strings.Contains(strings.ToLower(contact.First), query) {
			results = append(results, contact)
			continue
		}

		if strings.Contains(strings.ToLower(contact.Last), query) {
			results = append(results, contact)
			continue
		}

		if strings.Contains(strings.ToLower(contact.Phone), query) {
			results = append(results, contact)
			continue
		}

		if strings.Contains(contact.Email, query) {
			results = append(results, contact)
			continue
		}
	}

	return results, nil
}

func (db *JsonDatabase) Delete(contact *domain.Contact) error {
	for i, con := range db.contacts {
		if con.ID == contact.ID {
			if len(db.contacts) == i+1 {
				db.contacts = db.contacts[:i]
				saveToFile(db.contacts)
				return nil
			}

			db.contacts = append(db.contacts[:i], db.contacts[i+1:]...)
			saveToFile(db.contacts)
			return nil
		}
	}

	return errors.New("could not delete contact, no such contact found")
}

func (db *JsonDatabase) Update(contact *domain.Contact) (*domain.Contact, error) {
	for i, con := range db.contacts {
		if con.ID == contact.ID {
			db.contacts[i] = contact
			saveToFile(db.contacts)
			return contact, nil
		}
	}

	return nil, errors.New("could not update contact, no such contact found")
}

func (db *JsonDatabase) Find(id uuid.UUID) *domain.Contact {
	for _, c := range db.contacts {
		if c.ID == id {
			return c
		}
	}

	return nil
}
