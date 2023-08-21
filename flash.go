package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Flash struct {
	store *session.Store
}

func InitFlash() *Flash {
	return &Flash{
		store: session.New(),
	}
}

func (f *Flash) Get(c *fiber.Ctx) (string, error) {
	sess, err := f.store.Get(c)
	if err != nil {
		return "", err
	}

	message := sess.Get("flash")
	if message == nil {
		return "", nil
	}

	err = sess.Destroy()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", message), nil
}

func (f *Flash) Set(c *fiber.Ctx, msg string) error {
	sess, err := f.store.Get(c)
	if err != nil {
		return err
	}

	sess.Set("flash", msg)

	if err := sess.Save(); err != nil {
		return err
	}

	return nil
}
