package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	type Person struct {
		Name string
		Age  int
	}

	jonas := &Person{
		Name: "jonas",
		Age:  36,
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":  "Index!",
			"Person": jonas,
		}, "layouts/main")
	})

	app.Listen(":3000")
}
