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

type Database struct {
	contacts []*Contact
}

func OpenDatabase() (*Database, error) {
	data, err := os.ReadFile("db.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return &Database{}, nil
		}

		return nil, err
	}

	contacts := []*Contact{}
	err = json.Unmarshal(data, &contacts)
	if err != nil {
		return nil, err
	}

	return &Database{
		contacts: contacts,
	}, nil
}

func (db *Database) Save(contact *Contact) (*Contact, error) {
	contact.ID = uuid.New()

	validated := db.validate(contact)

	if validated.Errors != nil {
		return validated, errors.New("invalid contact")
	}

	db.contacts = append(db.contacts, validated)
	err := saveToFile(db.contacts)
	if err != nil {
		return validated, err
	}

	return validated, nil
}

func (db *Database) validate(contact *Contact) *Contact {
	if contact.Email == "" {
		contact.Errors = map[string]string{
			"email": "Email required.",
		}

		return contact
	}

	for _, c := range db.contacts {
		if c.Email == contact.Email {
			contact.Errors = map[string]string{
				"email": "Email already taken.",
			}
		}
	}

	// reset errors on struct just in case
	contact.Errors = nil

	return contact
}

func saveToFile(contacts []*Contact) error {
	marshalledJson, _ := json.MarshalIndent(contacts, "", "  ")
	err := os.WriteFile("db.json", marshalledJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) All() ([]*Contact, error) {
	if db.contacts == nil {
		return nil, nil
	}

	return db.contacts, nil
}

func (db *Database) Search(q string) ([]*Contact, error) {
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

func (db *Database) Delete(contact *Contact) error {
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

func (db *Database) Update(contact *Contact) (*Contact, error) {
	validated := db.validate(contact)

	if validated.Errors != nil {
		return validated, errors.New("invalid contact")
	}

	for i, con := range db.contacts {
		if con.ID == contact.ID {
			db.contacts[i] = contact
			saveToFile(db.contacts)
			return contact, nil
		}
	}

	return nil, errors.New("could not update contact, no such contact found")
}

func (db *Database) Find(id string) (*Contact, error) {
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
