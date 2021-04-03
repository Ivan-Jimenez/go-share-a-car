package main

import (
	"log"

	"github.com/Ivan-Jimenez/go-share-a-car/api/handlers"
	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	app.Post("/api/v1/user/signup", handlers.NewUser)
}

func main() {
	if err := database.Connect(); err != nil {
		log.Fatal("[ERROR] Faild to connect to database: %s", err)
	}
	log.Println("[DEBUG] Database connected!!!")
	app := fiber.New()
	setupRoutes(app)
	log.Fatal(app.Listen(":5000"))
}
