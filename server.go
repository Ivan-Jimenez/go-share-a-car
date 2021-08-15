package main

import (
	"log"
	"os"

	"github.com/Ivan-Jimenez/go-share-a-car/api/routes"
	"github.com/Ivan-Jimenez/go-share-a-car/database"
	"github.com/Ivan-Jimenez/go-share-a-car/pkg/user"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	l := log.New(os.Stdout, "[shareacar-api] ", log.LstdFlags)

	if err := database.Connect(); err != nil {
		l.Fatal("[ERROR] Faild to connect to database: ", err.Error())
	}
	l.Println("[INFO] Database connected!!!")

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Send([]byte("Share a car API"))
	})
	api := app.Group("/api")
	v1 := api.Group("/v1")

	routes.UserRouter(v1, user.NewService())

	// setupRoutes(app, l)
	log.Fatal(app.Listen(":5000"))
}
