package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
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

	app.Get("/api/v1/echo", func(ctx *fiber.Ctx) error {
		return ctx.SendString("ok")
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("ok")
	})

	app.Get("/api/v1/customer/:id", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"customer": "Acme Corp LLC",
		})
	})

	Strangle(app, "/*", "localhost:8000")

	log.Panic(app.Listen(":3000"))
}

func Strangle(app *fiber.App, route string, forwardHost string) {

	app.All(route, func(ctx *fiber.Ctx) error {

		forwardURI := fasthttp.AcquireURI()
		defer fasthttp.ReleaseURI(forwardURI)

		ctx.Request().URI().CopyTo(forwardURI)

		forwardURI.SetHost(forwardHost)

		addr := forwardURI.String()

		if err := proxy.DoTimeout(ctx, addr, time.Second*5); err != nil {
			return err
		}

		ctx.Response().Header.Del(fiber.HeaderServer)

		return nil
	})
}
