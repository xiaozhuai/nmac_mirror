package main

import (
	"flag"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	var configurationFile string
	flag.StringVar(&configurationFile, "config", "./config.yaml", `Configuration file path`)
	flag.Parse()

	configuration := LoadConfig(configurationFile)

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

	RegisterNMacService(configuration.Proxy, configuration.UserAgent, configuration.UseImageCache)

	app.HandleDir("/", "public", iris.DirOptions{
		Asset:      GzipAsset,
		AssetInfo:  GzipAssetInfo,
		AssetNames: GzipAssetNames,
		AssetValidator: func(ctx iris.Context, name string) bool {
			ctx.Header("Vary", "Accept-Encoding")
			ctx.Header("Content-Encoding", "gzip")
			return true
		},
	})

	app.Handle("GET", "/api/categories", hero.Handler(Categories))
	app.Handle("GET", "/api/list", hero.Handler(List))
	app.Handle("GET", "/api/detail", hero.Handler(Detail))
	app.Handle("GET", "/api/direct_url", hero.Handler(DirectUrl))
	app.Handle("GET", "/api/previous_version", hero.Handler(PreviousVersion))
	app.Handle("GET", "/api/fetch_image", hero.Handler(FetchImage))

	err := app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
	if err != nil {
		panic(err)
	}
}
