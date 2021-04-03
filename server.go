package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDatabase() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost/shareacardb"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
}

func setupRoutes(a *fiber.App) {

}

func main() {
	app := fiber.New()

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.SendString("Welcome to the rice fields mother fucker!")
	// })
	initDatabase()

	app.Listen(":5000")
}
