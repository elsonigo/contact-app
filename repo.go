package main

import "github.com/google/uuid"

type ContactRepo interface {
	Save(*Contact) (*Contact, error)
	Update(*Contact) (*Contact, error)
	All() ([]*Contact, error)
	Search(string) ([]*Contact, error)
	Delete(*Contact) error
	Find(uuid.UUID) (*Contact, error)
}
