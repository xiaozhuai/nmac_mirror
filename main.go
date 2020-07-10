package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	configuration := LoadConfig("config.yaml")

	app := iris.New()

	logOutput := configuration.GetLogFile()
	defer logOutput.Close()
	app.Logger().SetOutput(logOutput)
	app.Logger().SetLevel(configuration.LogLevel)

	app.Use(recover.New())
	app.Use(logger.New())

	hero.Register(app.Logger())

	configuration.PrepareDirs()
	cache := RegisterCacheService(configuration.MaxCacheDbSize, configuration.CacheDbDir, configuration.CacheImageDir)
	defer cache.Close()

	RegisterNMacService(configuration.Proxy, configuration.UserAgent)

	app.HandleDir("/", "./public", iris.DirOptions{
		Asset:      GzipAsset,
		AssetInfo:  GzipAssetInfo,
		AssetNames: GzipAssetNames,
		AssetValidator: func(ctx iris.Context, name string) bool {
			ctx.Header("Content-Encoding", "gzip")
			return true
		},
	})

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
