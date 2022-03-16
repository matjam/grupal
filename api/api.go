package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matjam/grupal/model"
)

func Start() {
	app := fiber.New()
	apiRoot := app.Group("/api")

	v1 := apiRoot.Group("/v1")
	v1.Get("/user/:id", GetHandler[model.User]("id"))
	v1.Get("/user/email/:email", GetHandler[model.User]("email"))
	v1.Post("/user", PostHandler[model.User]())
	app.Listen(":3000")
}

// GetHandler returns a fiber.Handler configured to handle fetching a model based on the provided type parameter and
// the field name provided which must match the field in the model as well as the parameter passed into the request.
func GetHandler[T any](field string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		row, err := model.Get[T](field, c.Params(field))
		if err != nil {
			return err
		}
		return c.JSON(row)
	}
}

// PostHandler returns a fiber.Handler configured to handle creating models provided via the type parameter.
func PostHandler[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		row := new(T)
		err := c.BodyParser(row)
		if err != nil {
			return err
		}

		newRow, err := model.Create[T](*row)
		if err != nil {
			return c.JSON(err)
		}
		return c.JSON(newRow)
	}
}
