package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Contact struct {
	ID     uuid.UUID         `json:"id"`
	Email  string            `json:"email"`
	First  string            `json:"first,omitempty"`
	Last   string            `json:"last,omitempty"`
	Phone  string            `json:"phone,omitempty"`
	Errors map[string]string `json:"errors,omitempty"`
}

// https://github.com/bigskysoftware/contact-app/blob/master/contacts_model.py#L92
type ContactRepo interface {
	Save(*Contact) (*Contact, error)
	Update(*Contact) (*Contact, error)
	All() ([]*Contact, error)
	Search(string) ([]*Contact, error)
	Delete(*Contact) error
	Find(string) (*Contact, error)
}

type JsonDatabase struct {
	contacts []*Contact
}

func OpenJsonDatabase() (*JsonDatabase, error) {
	data, err := os.ReadFile("db.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return &JsonDatabase{}, nil
		}

		return nil, err
	}

	contacts := []*Contact{}
	err = json.Unmarshal(data, &contacts)
	if err != nil {
		return nil, err
	}

	return &JsonDatabase{
		contacts: contacts,
	}, nil
}

var _ ContactRepo = &JsonDatabase{}

func (db *JsonDatabase) Save(contact *Contact) (*Contact, error) {
	db.contacts = append(db.contacts, contact)
	err := saveToFile(db.contacts)
	if err != nil {
		return contact, err
	}

	return contact, nil
}

func saveToFile(contacts []*Contact) error {
	marshalledJson, _ := json.MarshalIndent(contacts, "", "  ")
	err := os.WriteFile("db.json", marshalledJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *JsonDatabase) All() ([]*Contact, error) {
	if db.contacts == nil {
		return nil, nil
	}

	return db.contacts, nil
}

func (db *JsonDatabase) Search(q string) ([]*Contact, error) {
	if q == "" {
		return nil, errors.New("no query string given")
	}

	results := []*Contact{}
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

func (db *JsonDatabase) Delete(contact *Contact) error {
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

func (db *JsonDatabase) Update(contact *Contact) (*Contact, error) {
	for i, con := range db.contacts {
		if con.ID == contact.ID {
			db.contacts[i] = contact
			saveToFile(db.contacts)
			return contact, nil
		}
	}

	return nil, errors.New("could not update contact, no such contact found")
}

func (db *JsonDatabase) Find(id string) (*Contact, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id provided")
	}

	for _, c := range db.contacts {
		if c.ID == parsed {
			return c, nil
		}
	}

	return nil, nil
}
