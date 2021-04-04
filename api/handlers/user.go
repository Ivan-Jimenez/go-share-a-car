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

type login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

// TODO(Ivan): Generate token.
func (users *Users) Login(c *fiber.Ctx) error {
	credentials := new(login)
	if err := c.BodyParser(credentials); err != nil {
		users.logger.Printf("[INFO][Login-Handler] %s", err.Error())
		return c.Status(400).SendString(err.Error())
	}

	filter := bson.D{{Key: "email", Value: credentials.Email}}
	user, err := users.data.FindUser(c.Context(), filter)
	if err != nil {
		users.logger.Printf("[INFO][Login-Handler] %s", err.Error())
		return c.Status(400).SendString("Email or password invalid")
	}

	if !user.DoPasswordMatch(credentials.Password) {
		users.logger.Printf("[INFO][Login-Handler] Invalid password")
		return c.Status(400).SendString("Email or password invalid")
	}
	c.Status(200).SendString("Login succeded")

	return nil
}
