package main

import (
	"log"
	"os"

	"github.com/Ivan-Jimenez/go-share-a-car/api/handlers"
	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App, l *log.Logger) {
	userHandlers := handlers.NewUsers(l)
	app.Post("/api/v1/user/signup", userHandlers.NewUser)
}

func main() {
	l := log.New(os.Stdout, "shareacar-api", log.LstdFlags)

	if err := database.Connect(); err != nil {
		l.Fatal("[ERROR] Faild to connect to database: %s", err)
	}
	l.Println("[DEBUG] Database connected!!!")

	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app, l)
	log.Fatal(app.Listen(":5000"))
}
