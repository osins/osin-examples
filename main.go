package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django"
	"github.com/osins/osin-storage/storage/pg"
	"sso.humanrisk.cn/auth"
	"sso.humanrisk.cn/route"
	"sso.humanrisk.cn/util"
)

func main() {
	err := util.NewEnv().Load()
	if err != nil {
		return
	}

	fmt.Println(pg.GetPostgresDSN())
	pg.DB()
	pg.Migrate()

	auth.NewClient().Init()

	engine := django.New("./template/view", ".django")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	route := route.New()

	// Authorization code endpoint
	app.Get("/oauth/authorize", route.Authorize)
	app.Post("/oauth/authorize", route.Authorize)

	// Access token endpoint
	app.Post("/oauth/token", route.Token)

	app.Get("/api/user", route.Info)

	if err := app.Listen(":14000"); err != nil {
		fmt.Printf("fiber boot error: %s", err.Error())
	}
}
