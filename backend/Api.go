package main

import (
	"github.com/kataras/iris/v12"
)

func Categories(ctx iris.Context, ns NMacService) {
	categories := ns.GetCategories()
	ctx.JSON(iris.Map{
		"code":    0,
		"message": "ok",
		"data":    categories,
	})
}

func List(ctx iris.Context, ns NMacService) {
	category := ctx.URLParamDefault("category", "")
	page := ctx.URLParamIntDefault("page", 1)
	data, err := ns.GetList(category, page)
	code := 0
	msg := "ok"
	if err != nil {
		code = 1
		msg = err.Error()
	}
	ctx.JSON(iris.Map{
		"code":    code,
		"message": msg,
		"data":    data,
	})
}

func Search(ctx iris.Context, ns NMacService) {
	searchText := ctx.URLParamDefault("s", "")
	page := ctx.URLParamIntDefault("page", 1)
	data, err := ns.Search(searchText, page)
	code := 0
	msg := "ok"
	if err != nil {
		code = 1
		msg = err.Error()
	}
	ctx.JSON(iris.Map{
		"code":    code,
		"message": msg,
		"data":    data,
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

	if !ns.AllowUrl(u) {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Query param "url" is not a nmac.to url`,
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

	if !ns.AllowUrl(u) {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Query param "url" is not a nmac.to url`,
			"data":    nil,
		})
		return
	}

	directUrl, exists := cache.GetDirectUrl(u)
	if exists {
		ctx.Header("Mirror-Cache", "hit")
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

	ctx.Header("Mirror-Cache", "miss")
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

	if !ns.AllowUrl(u) {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Query param "url" is not a nmac.to url`,
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

func ImageCache(ctx iris.Context, ns NMacService, cache CacheService) {
	u := ctx.URLParamDefault("url", "")
	if u == "" {
		ctx.StatusCode(404)
		return
	}

	if !ns.AllowUrl(u) {
		ctx.StatusCode(403)
		return
	}

	if !ns.UseImageCache() {
		ctx.StatusCode(403)
		return
	}

	contentType, data, exists := cache.GetImageCache(u)
	if exists {
		ctx.Header("Mirror-Cache", "hit")
		ctx.ContentType(contentType)
		ctx.Write(data)
		return
	}

	contentType, data, err := ns.FetchImage(u)
	if err != nil {
		ctx.StatusCode(404)
		return
	}

	cache.SetImageCache(u, contentType, data)

	ctx.Header("Mirror-Cache", "miss")
	ctx.ContentType(contentType)
	ctx.Write(data)
}
