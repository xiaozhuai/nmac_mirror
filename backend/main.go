package main

import (
	"flag"
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"strings"
)

func main() {
	var configurationFile string
	flag.StringVar(&configurationFile, "config", "./config.yaml", `Configuration file path`)
	flag.Parse()

	configuration := LoadConfig(configurationFile)
	configuration.PrepareDirs()

	logOutput := configuration.GetLogFile()
	defer logOutput.Close()

	cache := NewCacheService(configuration.MaxCacheDbSize, configuration.CacheDbDir, configuration.CacheImageDir)
	defer cache.Close()

	ns := NewNMacService(configuration.Proxy, configuration.UserAgent, configuration.UseImageCache)

	app := iris.New()
	app.Logger().SetOutput(logOutput)
	app.Logger().SetLevel(configuration.LogLevel)
	app.Use(recover.New())
	app.Use(logger.New())
	if configuration.HttpsSupport && configuration.RedirectToHttps {
		app.Use(AutoRedirectToHttpsMiddleware(configuration.HttpsPort, configuration.RedirectToHttpsCode))
	}

	app.HandleDir("/", "public", AssetsDirOptions())

	app.Use(iris.CompressReader)
	app.Use(iris.Compress)

	app.ConfigureContainer(ApiBuilder(configuration, app.Logger(), cache, ns))

	if configuration.HttpsSupport {
		go app.Run(
			iris.Addr(fmt.Sprintf("%s:%d", configuration.ListenAddress, configuration.HttpPort)),
			iris.WithoutServerError(iris.ErrServerClosed),
		)

		app.Run(
			iris.TLS(
				fmt.Sprintf("%s:%d", configuration.ListenAddress, configuration.HttpsPort),
				configuration.CertFile,
				configuration.KeyFile,
				iris.TLSNoRedirect,
			),
			iris.WithoutServerError(iris.ErrServerClosed),
		)
	} else {
		app.Run(
			iris.Addr(fmt.Sprintf("%s:%d", configuration.ListenAddress, configuration.HttpPort)),
			iris.WithoutServerError(iris.ErrServerClosed),
		)
	}
}

func AssetsDirOptions() iris.DirOptions {
	return iris.DirOptions{
		IndexName: "/index.html",
		PushTargets: map[string][]string{
			"/": GetPushTargets(),
		},
		Compress:   false,
		ShowList:   false,
		Asset:      GzipAsset,
		AssetInfo:  GzipAssetInfo,
		AssetNames: GzipAssetNames,
		AssetValidator: func(ctx iris.Context, name string) bool {
			ctx.Header("Content-Encoding", "gzip")
			return true
		},
	}
}

func AutoRedirectToHttpsMiddleware(httpsPort int, redirectCode int) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		if ctx.Request().TLS == nil {
			h := ctx.Request().Host
			if pos := strings.Index(h, ":"); pos != -1 {
				h = h[0:pos]
			}
			uri := ctx.Request().RequestURI

			var httpsUrl string
			if httpsPort == 443 {
				httpsUrl = fmt.Sprintf("https://%s%s", h, uri)
			} else {
				httpsUrl = fmt.Sprintf("https://%s:%d%s", h, httpsPort, uri)
			}

			ctx.Redirect(httpsUrl, redirectCode)
			return
		}
		ctx.Next()
	}
}

func ApiBuilder(configuration *Configuration, logger *golog.Logger, cache CacheService, ns NMacService) func(api *iris.APIContainer) {
	return func(api *iris.APIContainer) {
		api.RegisterDependency(configuration)
		api.RegisterDependency(logger)
		api.RegisterDependency(cache)
		api.RegisterDependency(ns)

		api.Get("/api/categories", Categories)
		api.Get("/api/list", List)
		api.Get("/api/search", Search)
		api.Get("/api/detail", Detail)
		api.Get("/api/direct_url", DirectUrl)
		api.Get("/api/previous_version", PreviousVersion)
		api.Get("/image_cache", ImageCache)
	}
}

func GetPushTargets() []string {
	assetNames := GzipAssetNames()
	pushTargets := make([]string, 0)
	for _, name := range assetNames {
		// ignore *.map file
		if !strings.HasSuffix(name, ".map") && !strings.HasSuffix(name, ".ttf") && !strings.HasSuffix(name, ".woff") {
			target := strings.TrimPrefix(name, "public")

			// skip /index.html
			if target != "/index.html" {
				pushTargets = append(pushTargets, target)
			}
		}
	}
	return pushTargets
}
