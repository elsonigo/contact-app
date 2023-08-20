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
	First  string            `json:"first"`
	Last   string            `json:"last"`
	Phone  string            `json:"phone"`
	Email  string            `json:"email"`
	Errors map[string]string `json:"errors"`
}

// https://github.com/bigskysoftware/contact-app/blob/master/contacts_model.py#L92

type Database struct {
	contacts []Contact
}

func OpenDatabase() (*Database, error) {
	data, err := os.ReadFile("db.json")
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return &Database{}, nil
		}

		return nil, err
	}

	contacts := []Contact{}
	err = json.Unmarshal(data, &contacts)
	if err != nil {
		return nil, err
	}

	return &Database{
		contacts: contacts,
	}, nil
}

func (db *Database) Save(contact *Contact) error {
	db.contacts = append(db.contacts, *contact)
	saveToFile(db.contacts)

	return nil
}

func saveToFile(contacts []Contact) error {
	marshalledJson, _ := json.MarshalIndent(contacts, "", "  ")
	err := os.WriteFile("db.json", marshalledJson, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) All() ([]Contact, error) {
	if db.contacts == nil {
		return []Contact{}, nil
	}

	return db.contacts, nil
}

func (db *Database) Search(q string) ([]Contact, error) {
	if q == "" {
		return nil, errors.New("no query string given")
	}

	results := []Contact{}
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

func (db *Database) Update(contact *Contact) error {
	for i, con := range db.contacts {
		if con.ID == contact.ID {
			db.contacts[i] = *contact
			saveToFile(db.contacts)
			return nil
		}
	}

	return errors.New("could not update contact, no such contact found")
}
