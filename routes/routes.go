package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matjam/grupal/database"
)

type CRUDRouter struct {
	*database.DB
	*fiber.App
}

func NewCRUDRouter(db database.DB) CRUDRouter {
	var crud CRUDRouter
	crud.App = fiber.New()

	api := crud.App.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/users/:id", crud.GetHandler(db.User, "id"))
	//v1.Get("/users", SearchHandler[api.User]())
	//v1.Post("/users", PostHandler[api.User]())
	//v1.Put("/users/:id", PutHandler[api.User]("id"))

	return crud
}

func (crud CRUDRouter) Start() {
	err := crud.Listen(":3000")
	if err != nil {
		panic("can't listen on port")
	}
}

// GetHandler returns a fiber.Handler configured to handle fetching a database based on the provided type parameter and
// the field name provided which must match the field in the database as well as the parameter passed into the request.
func (crud CRUDRouter) GetHandler(model database.Model, field string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		filter := map[string]any{field: c.Params(field)}
		row, err := model.Read(filter, 1, 0, nil)
		if err != nil {
			return err
		}
		return c.JSON(row)
	}
}

//
//// PostHandler returns a fiber.Handler configured to handle creating models provided via the type parameter.
//func PostHandler[T any]() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		row := new(T)
//		err := c.BodyParser(row)
//		if err != nil {
//			return err
//		}
//
//		newRow, err := database.Create[T](*row)
//		if err != nil {
//			return c.JSON(err)
//		}
//		return c.JSON(newRow)
//	}
//}
//
//func PutHandler[T any](field string) fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		row := new(T)
//		err := c.BodyParser(row)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
//}
//
//func SearchHandler[T any]() fiber.Handler {
//	return func(c *fiber.Ctx) error {
//		ft := map[string]any{}
//		ftString := c.Query("filter", "{}")
//		err := json.Unmarshal([]byte(ftString), &ft)
//		if err != nil {
//			log.Errorf("bad filter %v", ftString)
//			return c.Status(fiber.StatusBadRequest).SendString("Unacceptable filter")
//		}
//
//		var st []string
//		stString := c.Query("sort", "[]")
//		err = json.Unmarshal([]byte(stString), &st)
//		if err != nil || (len(st) != 2 && len(st) != 0) {
//			log.Errorf("bad sort %v", stString)
//			return c.Status(fiber.StatusBadRequest).SendString("Unacceptable sort")
//		}
//
//		var rg []int
//		rgString := c.Query("range", "[0,0]")
//		err = json.Unmarshal([]byte(rgString), &rg)
//		if err != nil || (len(rg) != 2 && len(rg) != 0) {
//			log.Errorf("bad range %v", rgString)
//			return c.Status(fiber.StatusBadRequest).SendString("Unacceptable range")
//		}
//
//		log.Infof("GET %v filter %v sort %v range %v", c.Request().URI().String(), ft, st, rg)
//
//		result, err := database.Search[T](ft, rg[0], rg[1], st)
//		if err != nil {
//			log.Errorf("database error: %v", err.Error())
//			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
//		}
//		return c.JSON(result)
//	}
//}
