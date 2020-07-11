package main

import (
	"flag"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"strings"
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

	// Auto redirect to https
	app.Use(func(ctx iris.Context) {
		if configuration.HttpsSupport && configuration.RedirectToHttps && ctx.Request().TLS == nil {
			host := ctx.Request().Host
			if pos := strings.Index(host, ":"); pos != -1 {
				host = host[0:pos]
			}
			uri := ctx.Request().RequestURI

			var httpsUrl string
			if configuration.HttpsPort == 443 {
				httpsUrl = fmt.Sprintf("https://%s%s", host, uri)
			} else {
				httpsUrl = fmt.Sprintf("https://%s:%d%s", host, configuration.HttpsPort, uri)
			}

			ctx.Redirect(httpsUrl, configuration.RedirectToHttpsCode)
			return
		}
		ctx.Next()
	})

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
	app.Handle("GET", "/api/search", hero.Handler(Search))
	app.Handle("GET", "/api/detail", hero.Handler(Detail))
	app.Handle("GET", "/api/direct_url", hero.Handler(DirectUrl))
	app.Handle("GET", "/api/previous_version", hero.Handler(PreviousVersion))
	app.Handle("GET", "/api/fetch_image", hero.Handler(FetchImage))

	if configuration.HttpsSupport {
		go func() {
			err := app.Run(
				iris.Addr(fmt.Sprintf("%s:%d", configuration.ListenAddress, configuration.HttpPort)),
				iris.WithoutServerError(iris.ErrServerClosed),
			)
			if err != nil {
				panic(err)
			}
		}()

		err := app.Run(
			iris.TLS(
				fmt.Sprintf("%s:%d", configuration.ListenAddress, configuration.HttpsPort),
				configuration.CertFile,
				configuration.KeyFile,
			),
			iris.WithoutServerError(iris.ErrServerClosed),
		)
		if err != nil {
			panic(err)
		}
	} else {
		err := app.Run(
			iris.Addr(fmt.Sprintf("%s:%d", configuration.ListenAddress, configuration.HttpPort)),
			iris.WithoutServerError(iris.ErrServerClosed),
		)
		if err != nil {
			panic(err)
		}
	}

}
