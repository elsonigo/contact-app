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
		panic(err.Error())
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
			"Title":    "Contacts",
			"Contacts": foundContacts,
			"Query":    query,
		}, "layouts/main")
	})

	app.Listen(":3000")
}
