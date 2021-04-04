package handlers

import (
	"log"

	"github.com/Ivan-Jimenez/go-share-a-car/api/data"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type Users struct {
	data   *data.UserData
	logger *log.Logger
}

func NewUsers(logger *log.Logger) *Users {
	return &Users{
		data.NewUserData(logger),
		logger,
	}
}

func (users *Users) NewUser(c *fiber.Ctx) error {
	user := new(data.User)
	if err := c.BodyParser(user); err != nil {
		users.logger.Printf("[INFO][NewUser-Handler] %s", err.Error())
		return c.Status(400).SendString(err.Error())
	}

	if err := user.Validate(); err != nil {
		users.logger.Printf("[INFO][NewUser-Handler] %s", err.Error())
		return c.Status(400).SendString(err.Error())
	}

	filter := bson.D{{Key: "email", Value: user.Email}}
	if _, err := users.data.FindUser(c.Context(), filter); err == nil {
		users.logger.Printf("[INFO][NewUser-Handler] The email %s is already in the database.", user.Email)
		return c.Status(400).SendString("Email is already in the database. Use another one.")
	}

	if err := user.HashPassword(); err != nil {
		users.logger.Printf("[ERROR][NewUser-Handler] %s", err.Error())
		return c.Status(500).SendString(err.Error())
	}

	user, err := users.data.SaveUser(c.Context(), user)
	if err != nil {
		users.logger.Printf("[ERROR][NewUser-Handler] Fail to save user %s", err.Error())
		return c.Status(500).SendString(err.Error())
	}
	c.Status(201).JSON(user)

	return nil
}
