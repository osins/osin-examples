package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/django"
	"github.com/osins/osin-storage/storage"
	"github.com/osins/osins-examples/auth"
	"github.com/osins/osins-examples/route"
	"github.com/osins/osins-examples/util"
)

var (
	_, f, _, _ = runtime.Caller(0)
	BasePATH   = filepath.Dir(f)
	ENVFile    = BasePATH + "/.env"
)

func main() {
	now := time.Now()
	tz, _ := time.LoadLocation("Asia/Shanghai")
	parisTime := now.In(tz)
	fmt.Printf("Local time: %s\nParis time: %s\n", now, parisTime)

	err := util.NewEnv().Load(ENVFile)
	if err != nil {
		return
	}

	fmt.Println(pg.GetPostgresDSN())
	pg.DB()
	pg.Migrate()

	auth.NewClient().Init()

	engine := django.New(BasePATH+"/template/view", ".django")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New(cors.Config{
		Next:             nil,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	}))

	route := route.New()

	app.Static("/oauth/assets", BasePATH+"/template/assets")

	// Authorization code endpoint
	app.Get("/oauth/authorize", route.Authorize)
	app.Post("/oauth/authorize", route.Authorize)

	// Access token endpoint
	app.Post("/oauth/token", route.Token)

	app.Get("/api/user", route.Info)

	if err := app.Listen(os.Getenv("APP_LISTEN")); err != nil {
		fmt.Printf("fiber boot error: %s", err.Error())
	}
}
