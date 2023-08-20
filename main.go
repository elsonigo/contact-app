package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./static")

	db, err := OpenDatabase()
	if err != nil {
		panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/contacts", fiber.StatusTemporaryRedirect)
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		foundContacts := []Contact{}

		query := c.Query("q")
		if query != "" {
			found, err := db.Search(query)

			if err != nil {
				return fmt.Errorf("error getting contacts: %s", err.Error())
			}

			foundContacts = found
		}

		if query == "" {
			all, _ := db.All()
			foundContacts = all
		}

		return c.Render("contacts", fiber.Map{
			"Contacts": foundContacts,
			"Query":    query,
		}, "layouts/main")
	})

	app.Get("/contacts/new", func(c *fiber.Ctx) error {
		return c.Render("new", fiber.Map{}, "layouts/main")
	})

	app.Post("/contacts/new", func(c *fiber.Ctx) error {
		newContact := &Contact{
			Email: c.FormValue("email"),
			First: c.FormValue("first"),
			Last:  c.FormValue("last"),
			Phone: c.FormValue("phone"),
		}

		ct, err := db.Save(newContact)
		if err != nil {
			// implement flash message
			// https://www.alexedwards.net/blog/simple-flash-messages-in-golang
			return c.Render("new", fiber.Map{
				"Contact": ct,
			}, "layouts/main")
		}

		return c.Redirect("/contacts")
	})

	app.Listen(":3000")
}
