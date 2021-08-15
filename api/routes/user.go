package routes

import (
	"github.com/Ivan-Jimenez/go-share-a-car/pkg/entities"
	"github.com/Ivan-Jimenez/go-share-a-car/pkg/user"
	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, service user.Service) {
	router := app.Group("/user")
	router.Post("/signup", newUser(service))
	router.Post("/login", login(service))
}

func newUser(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.User
		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.JSON(&fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}
		result, dberr := service.NewUser(&requestBody)
		if dberr != nil {
			return c.JSON(&fiber.Map{
				"error":   true,
				"message": dberr.Error(),
			})
		}
		return c.JSON(&fiber.Map{
			"error": false,
			"data":  result,
		})
	}
}

func login(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.LoginCredentials
		err := c.BodyParser(&requestBody)
		if err != nil {
			return c.JSON(&fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		}
		result, token, dberr := service.Login(&requestBody)
		if dberr != nil {
			return c.JSON(&fiber.Map{
				"error":   true,
				"message": dberr.Error(),
			})
		}
		return c.JSON(&fiber.Map{
			"error": false,
			"data": fiber.Map{
				"token": token,
				"user":  result,
			},
		})
	}
}
