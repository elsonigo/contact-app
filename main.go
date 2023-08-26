package main

import (
	"github.com/elsonigo/contact-app/domain"
	"github.com/elsonigo/contact-app/repositories/json_db"
	"github.com/elsonigo/contact-app/services/contactsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Static("/static", "./static")

	flash := InitFlash()

	json_db, err := json_db.OpenJsonDatabase()
	if err != nil {
		panic(err)
	}
	cs := contactsrv.NewContactService(json_db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/contacts", fiber.StatusTemporaryRedirect)
	})

	app.Get("/contacts", func(c *fiber.Ctx) error {
		foundContacts := []*domain.Contact{}

		query := c.Query("q")
		if query != "" {
			found, err := cs.Search(query)

			if err != nil {
				return err
			}

			foundContacts = found
		}

		if query == "" {
			all := cs.All()
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
		newContact := &domain.Contact{
			Email: c.FormValue("email"),
			First: c.FormValue("first_name"),
			Last:  c.FormValue("last_name"),
			Phone: c.FormValue("phone"),
		}

		ct, err := cs.Save(newContact)
		if err != nil {
			return c.Render("new", fiber.Map{
				"Contact": ct,
			})
		}

		flash.Set(c, "new contact created")

		return c.Redirect("/contacts")
	})

	app.Get("/contacts/:id", func(c *fiber.Ctx) error {
		contact := cs.Find(c.Params("id"))
		if contact == nil {
			flash.Set(c, "could not find contact")
			return c.Redirect("/contacts")
		}

		return c.Render("show", fiber.Map{
			"Contact": contact,
		})
	})

	app.Get("/contacts/:id/edit", func(c *fiber.Ctx) error {
		contact := cs.Find(c.Params("id"))
		if contact == nil {
			flash.Set(c, "could not find contact")
			return c.Redirect("/contacts")
		}

		return c.Render("edit", fiber.Map{
			"Contact": contact,
		})
	})

	app.Post("/contacts/:id/edit", func(c *fiber.Ctx) error {
		contact := cs.Find(c.Params("id"))
		if contact == nil {
			return c.Redirect("/contacts")
		}

		contact.Email = c.FormValue("email")
		contact.First = c.FormValue("first_name")
		contact.Last = c.FormValue("last_name")
		contact.Phone = c.FormValue("phone")

		ct, err := cs.Update(contact)
		if err != nil {
			return c.Render("edit", fiber.Map{
				"Contact": ct,
			})
		}

		flash.Set(c, "contact updated")

		return c.Redirect("/contacts")
	})

	app.Delete("/contacts/:id", func(c *fiber.Ctx) error {
		contact := cs.Find(c.Params("id"))
		if contact == nil {
			return c.Redirect("/contacts", 303)
		}

		err = cs.Delete(contact)
		if err != nil {
			flash.Set(c, err.Error())
			return c.Redirect("/contacts", 303)
		}

		flash.Set(c, "contact deleted")

		return c.Redirect("/contacts", 303)
	})

	app.Get("/contacts/:id/email", func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Redirect("contacts")
		}

		contact := cs.Find(c.Params("id"))
		if contact == nil {
			return c.Redirect("contacts")
		}

		email := c.FormValue("email")
		validationError := cs.ValidateEmail(email, contact.ID)
		if validationError != nil {
			return c.SendString(validationError.Error())
		}

		return nil
	})

	app.Listen(":3000")
}
