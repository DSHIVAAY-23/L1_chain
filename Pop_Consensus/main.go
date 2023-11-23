package main

import (
	"fmt"
	"log"
	"pop_v1/config"
	"pop_v1/router"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	router.MainRoute(app)
	serverAddr := "0.0.0.0:"+config.Config("PORT")
	fmt.Printf("Starting server on %s\n", serverAddr)
	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}

}
