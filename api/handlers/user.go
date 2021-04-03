package handlers

import (
	"github.com/Ivan-Jimenez/go-share-a-car/api/data"
	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func NewUser(c *fiber.Ctx) error {
	collection := database.Instance.Database.Collection("users")

	user := new(data.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

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
