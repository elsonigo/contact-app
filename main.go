package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	flash := InitFlash()

	app.Static("/static", "./static")

	db, err := OpenDatabase()
	if err != nil {
		panic(err)
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/contacts", fiber.StatusTemporaryRedirect)
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		foundContacts := []*Contact{}

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

		f, _ := flash.Get(c)

		return c.Render("contacts", fiber.Map{
			"Contacts": foundContacts,
			"Query":    query,
			"Flash":    f,
		})
	})

	app.Get("/contacts/new", func(c *fiber.Ctx) error {
		return c.Render("new", fiber.Map{})
	})

	app.Post("/contacts/new", func(c *fiber.Ctx) error {
		newContact := &Contact{
			Email: c.FormValue("email"),
			First: c.FormValue("first_name"),
			Last:  c.FormValue("last_name"),
			Phone: c.FormValue("phone"),
		}

		ct, err := db.Save(newContact)
		if err != nil {
			return c.Render("new", fiber.Map{
				"Contact": ct,
			})
		}

		flash.Set(c, "new contact created")

		return c.Redirect("/contacts")
	})

	app.Get("/contacts/:id", func(c *fiber.Ctx) error {
		contact, err := db.Find(c.Params("id"))
		if contact.Email == "" || err != nil {
			flash.Set(c, "could not find contact")
			return c.Redirect("/contacts")
		}

		return c.Render("show", fiber.Map{
			"Contact": contact,
		})
	})

	app.Get("/contacts/:id/edit", func(c *fiber.Ctx) error {
		contact, err := db.Find(c.Params("id"))
		if contact.Email == "" || err != nil {
			flash.Set(c, "could not find contact")
			return c.Redirect("/contacts")
		}

		return c.Render("edit", fiber.Map{
			"Contact": contact,
		})
	})

	app.Post("/contacts/:id/edit", func(c *fiber.Ctx) error {
		contact, err := db.Find(c.Params("id"))
		if contact.Email == "" || err != nil {
			return c.Redirect("/contacts")
		}

		contact.Email = c.FormValue("email")
		contact.First = c.FormValue("first_name")
		contact.Last = c.FormValue("last_name")
		contact.Phone = c.FormValue("phone")

		ct, err := db.Update(contact)
		if err != nil {
			return c.Render("edit", fiber.Map{
				"Contact": ct,
			})
		}

		flash.Set(c, "contact updated")

		return c.Redirect("/contacts")
	})

	app.Post("/contacts/:id/delete", func(c *fiber.Ctx) error {
		contact, err := db.Find(c.Params("id"))
		if contact.Email == "" || err != nil {

			if err != nil {
				flash.Set(c, err.Error())
			}

			return c.Redirect("/contacts")
		}

		err = db.Delete(contact)
		if err != nil {
			flash.Set(c, err.Error())
			return c.Redirect("/contacts")
		}

		flash.Set(c, "contact deleted")

		return c.Redirect("/contacts")
	})

	app.Listen(":3000")
}
