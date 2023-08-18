package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Contact struct {
	Name string
	Age  int
}

// TODO: implement contacts struct -> create json db?
// https://github.com/bigskysoftware/contact-app/blob/master/contacts_model.py#L92
type Contacts struct{}

func (c *Contacts) all() ([]Contact, error) {
	return []Contact{
		{
			Name: "Jonas",
			Age:  36,
		},
	}, nil
}

// func (c *Contacts) search() {}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/contacts", fiber.StatusTemporaryRedirect)
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		contacts_data := []Contact{}
		contacts := Contacts{}

		search := c.Query("q")
		if search != "" {
			all, err := contacts.all()
			if err != nil {
				return fmt.Errorf("error getting contacts: %s", err.Error())
			}

			contacts_data = all
		}

		return c.Render("contacts", fiber.Map{
			"Title":    "Contacts",
			"Contacts": contacts_data,
		}, "layouts/main")
	})

	app.Listen(":3000")
}
