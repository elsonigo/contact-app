package main

type ContactRepo interface {
	Save(*Contact) (*Contact, error)
	Update(*Contact) (*Contact, error)
	All() ([]*Contact, error)
	Search(string) ([]*Contact, error)
	Delete(*Contact) error
	Find(string) (*Contact, error)
}
