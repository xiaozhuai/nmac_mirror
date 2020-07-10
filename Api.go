package main

import (
	"github.com/kataras/iris/v12"
)

func List(ctx iris.Context, ns NMacService) {
	category := ctx.URLParamDefault("category", "")
	page := ctx.URLParamIntDefault("page", 1)
	list, err := ns.GetList(category, page)
	code := 0
	msg := "ok"
	if err != nil {
		code = 1
		msg = err.Error()
	}
	ctx.JSON(iris.Map{
		"code":    code,
		"message": msg,
		"data": iris.Map{
			"category": category,
			"page":     page,
			"size":     len(list),
			"list":     list,
		},
	})
}

func Detail(ctx iris.Context, ns NMacService) {
	u := ctx.URLParamDefault("url", "")
	if u == "" {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Query param "url" cannot be empty!`,
			"data":    nil,
		})
		return
	}

	detail, err := ns.GetDetail(u)
	code := 0
	msg := "ok"
	if err != nil {
		code = 1
		msg = err.Error()
	}
	ctx.JSON(iris.Map{
		"code":    code,
		"message": msg,
		"data":    detail,
	})
}

func DirectUrl(ctx iris.Context, ns NMacService, cache CacheService) {
	u := ctx.URLParamDefault("url", "")
	if u == "" {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Query param "url" cannot be empty!`,
			"data":    nil,
		})
		return
	}

	directUrl, exists := cache.GetDirectUrl(u)
	if exists {
		ctx.JSON(iris.Map{
			"code":    0,
			"message": "ok",
			"data":    directUrl,
		})
		return
	}

	directUrl, err := ns.GetDirectUrl(u)
	cache.SetDirectUrl(u, directUrl)

	code := 0
	msg := "ok"
	if err != nil {
		code = 1
		msg = err.Error()
	}
	ctx.JSON(iris.Map{
		"code":    code,
		"message": msg,
		"data":    directUrl,
	})
}

func PreviousVersion(ctx iris.Context, ns NMacService) {
	u := ctx.URLParamDefault("url", "")
	if u == "" {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Query param "url" cannot be empty!`,
			"data":    nil,
		})
		return
	}

	versions := ns.GetPreviousVersion(u)
	ctx.JSON(iris.Map{
		"code":    0,
		"message": "ok",
		"data":    versions,
	})
}

func FetchImage(ctx iris.Context, ns NMacService, cache CacheService) {
	u := ctx.URLParamDefault("url", "")
	if u == "" {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Query param "url" cannot be empty!`,
			"data":    nil,
		})
		return
	}

	contentType, data, exists := cache.GetImageCache(u)
	if exists {
		ctx.ContentType(contentType)
		ctx.Write(data)
		return
	}

	contentType, data, err := ns.FetchImage(u)
	if err != nil {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Fetch image failed!`,
			"data":    nil,
		})
		return
	}

	cache.SetImageCache(u, contentType, data)

	ctx.ContentType(contentType)
	ctx.Write(data)
}
