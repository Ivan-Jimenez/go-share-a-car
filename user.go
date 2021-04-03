package main

import (
	"context"
	"log"

	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func NewUser(c *fiber.Ctx) error {
	collection := database.Instance.Collection("testCollection")
	res, err := collection.InsertOne(context.Background(), bson.M{"hello": "mother fuckers"})
	if err != nil {
		// log.Panic("[ERROR] Faild to save document")
		return err
	}

	log.Printf("[DEBUG] Inserted: %s", res)

	return c.SendString("Test user endpoint")
}
