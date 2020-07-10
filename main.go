package main

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"os"
)

func main() {
	configuration := LoadConfig("config.yaml")

	app := iris.New()

	if configuration.Log == "stdout" {
		app.Logger().SetOutput(os.Stdout)
	} else {
		logFile, err := os.Open(configuration.Log)
		if err != nil {
			panic(err)
		}
		app.Logger().SetOutput(logFile)
	}
	app.Logger().SetLevel(configuration.LogLevel)

	app.Use(recover.New())
	app.Use(logger.New())

	hero.Register(func(ctx iris.Context) *golog.Logger {
		return app.Logger()
	})

	_ = os.MkdirAll(configuration.CacheDbDir, 0777)
	_ = os.MkdirAll(configuration.CacheImageDir, 0777)
	cache := RegisterCacheService(configuration.MaxCacheDbSize, configuration.CacheDbDir, configuration.CacheImageDir)
	defer cache.Close()

	RegisterNMacService(configuration.Proxy, configuration.UserAgent)

	app.Handle("GET", "/list", hero.Handler(List))
	app.Handle("GET", "/detail", hero.Handler(Detail))
	app.Handle("GET", "/direct_url", hero.Handler(DirectUrl))
	app.Handle("GET", "/previous_version", hero.Handler(PreviousVersion))
	app.Handle("GET", "/fetch_image", hero.Handler(FetchImage))

	err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		panic(err)
	}
}
