package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	fconfig := fiber.Config{
		AppName:               "ProxyDemo",
		DisableStartupMessage: false,
		EnablePrintRoutes:     true,
		JSONDecoder:           json.Unmarshal,
		JSONEncoder:           json.Marshal,
		Prefork:               false,
	}

	app := fiber.New(fconfig)

	app.Use(logger.New())

	app.All("/*", func(ctx *fiber.Ctx) error {

		fmt.Println("=================")
		fmt.Printf("%s\n", ctx.Request().String())
		fmt.Println("=================")

		return ctx.JSON(fiber.Map{
			"message": "ok",
		})
	})

	log.Panic(app.Listen(":8000"))
}
