package main

import (
	"log"
	"os"

	"github.com/Ivan-Jimenez/go-share-a-car/api/controllers"
	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App, l *log.Logger) {
	userController := controllers.NewUserController(l)
	app.Post("/api/v1/user/signup", userController.NewUser)

	app.Post("/api/v1/user/login", userController.Login)
}

func main() {
	l := log.New(os.Stdout, "[shareacar-api] ", log.LstdFlags)

	if err := database.Connect(); err != nil {
		l.Fatal("[ERROR] Faild to connect to database: ", err.Error())
	} else {
		l.Println("[INFO] Database connected!!!")
	}

	app := fiber.New()
	app.Use(logger.New())

	setupRoutes(app, l)
	log.Fatal(app.Listen(":5000"))
}
