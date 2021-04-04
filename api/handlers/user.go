package handlers

import (
	"log"

	"github.com/Ivan-Jimenez/go-share-a-car/api/data"
	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type Users struct {
	logger *log.Logger
}

func NewUsers(logger *log.Logger) *Users {
	return &Users{logger}
}

func (users *Users) NewUser(c *fiber.Ctx) error {
	collection := database.Instance.Database.Collection("users")

	user := new(data.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := user.Validate(); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if err := user.HashPassword(); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	// filter := bson.D{{Key: "email", Value: user.Email}}
	// searchEmail:= collection.FindOne(c.Context(), filter)

	user.ID = ""
	res, err := collection.InsertOne(c.Context(), user)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	filter := bson.D{{Key: "_id", Value: res.InsertedID}}
	createdRecord := collection.FindOne(c.Context(), filter)

	createdUser := &data.User{}
	createdRecord.Decode(createdUser)
	return c.Status(201).JSON(createdUser)
}
