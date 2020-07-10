package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"os"
)

func main() {
	app := iris.New()

	app.Logger().SetOutput(os.Stdout)
	app.Logger().SetLevel("info")

	app.Use(recover.New())
	app.Use(logger.New())

	RegisterNMacService(
		"http://127.0.0.1:8118",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36",
	)

	app.Handle("GET", "/list", hero.Handler(List))
	app.Handle("GET", "/detail", hero.Handler(Detail))
	app.Handle("GET", "/direct_url", hero.Handler(DirectUrl))
	app.Handle("GET", "/previous_version", hero.Handler(PreviousVersion))
	app.Handle("GET", "/fetch_image", hero.Handler(FetchImage))

	err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
}
