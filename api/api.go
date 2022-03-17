package api

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"github.com/matjam/grupal/model"
)

func Start() {
	app := fiber.New()
	apiRoot := app.Group("/api")

	v1 := apiRoot.Group("/v1")
	v1.Get("/users/:id", GetHandler[model.User]("id"))
	v1.Get("/users", SearchHandler[model.User]())
	v1.Post("/users", PostHandler[model.User]())
	v1.Put("/users/:id", PutHandler[model.User]("id"))
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

func PutHandler[T any](field string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		row := new(T)
		err := c.BodyParser(row)
		if err != nil {
			return err
		}

		return nil
	}
}

func SearchHandler[T any]() fiber.Handler {
	return func(c *fiber.Ctx) error {
		ft := map[string]any{}
		ftString := c.Query("filter", "{}")
		err := json.Unmarshal([]byte(ftString), &ft)
		if err != nil {
			log.Errorf("bad filter %v", ftString)
			return c.Status(fiber.StatusBadRequest).SendString("Unacceptable filter")
		}

		var st []string
		stString := c.Query("sort", "[]")
		err = json.Unmarshal([]byte(stString), &st)
		if err != nil || (len(st) != 2 && len(st) != 0) {
			log.Errorf("bad sort %v", stString)
			return c.Status(fiber.StatusBadRequest).SendString("Unacceptable sort")
		}

		var rg []int
		rgString := c.Query("range", "[0,0]")
		err = json.Unmarshal([]byte(rgString), &rg)
		if err != nil || (len(rg) != 2 && len(rg) != 0) {
			log.Errorf("bad range %v", rgString)
			return c.Status(fiber.StatusBadRequest).SendString("Unacceptable range")
		}

		log.Infof("GET %v filter %v sort %v range %v", c.Request().URI().String(), ft, st, rg)

		result, err := model.Search[T](ft, rg[0], rg[1], st)
		if err != nil {
			log.Errorf("database error: %v", err.Error())
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		return c.JSON(result)
	}
}
