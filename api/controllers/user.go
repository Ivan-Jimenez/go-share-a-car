package controllers

import (
	"log"

	"github.com/Ivan-Jimenez/go-share-a-car/api/models"
	"github.com/Ivan-Jimenez/go-share-a-car/api/util"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type Users struct {
	data   *models.UserData
	logger *log.Logger
}

func NewUserController(logger *log.Logger) *Users {
	return &Users{
		models.NewUserData(logger),
		logger,
	}
}

func (users *Users) NewUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		users.logger.Printf("[INFO][NewUser-Handler] %s", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := user.Validate(); err != nil {
		users.logger.Printf("[INFO][NewUser-Handler] %s", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	filter := bson.D{{Key: "email", Value: user.Email}}
	if _, err := users.data.FindUser(c.Context(), filter); err == nil {
		users.logger.Printf("[INFO][NewUser-Handler] The email %s is already in the database.", user.Email)
		return c.Status(fiber.StatusBadRequest).SendString("Email is already in the database. Use another one.")
	}

	if err := user.HashPassword(); err != nil {
		users.logger.Printf("[ERROR][NewUser-Handler] %s", err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	// TODO: Don't return password field
	user, err := users.data.SaveUser(c.Context(), user)
	if err != nil {
		users.logger.Printf("[ERROR][NewUser-Handler] Fail to save user %s", err.Error())
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	c.Status(fiber.StatusOK).JSON(user)

	return nil
}

// TODO(Ivan): Generate token.
func (users *Users) Login(c *fiber.Ctx) error {
	credentials := new(models.LoginCredentials)
	if err := c.BodyParser(credentials); err != nil {
		users.logger.Printf("[INFO][Login-Handler] %s", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	filter := bson.D{{Key: "email", Value: credentials.Email}}
	user, err := users.data.FindUser(c.Context(), filter)
	if err != nil {
		users.logger.Printf("[INFO][Login-Handler] %s", err.Error())
		return c.Status(fiber.StatusBadRequest).SendString("Email or password invalid")
	}

	if !user.DoPasswordMatch(credentials.Password) {
		users.logger.Printf("[INFO][Login-Handler] Invalid password")
		return c.Status(fiber.StatusBadRequest).SendString("Email or password invalid")
	}

	accToken, refToken := util.GenerateTokens(user.ID)
	accCookie, refCookie := util.GetAuthCookies(accToken, refToken)
	c.Cookie(accCookie)
	c.Cookie(refCookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  accToken,
		"refreshToken": refToken,
	})
}

func (users *Users) EmailVefification(c *fiber.Ctx) error {
	// TODO(Ivan): Emal vefication
	return nil
}
