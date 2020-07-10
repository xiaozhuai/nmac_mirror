package main

import (
	"github.com/kataras/iris/v12"
	"io/ioutil"
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

func DirectUrl(ctx iris.Context, ns NMacService) {
	u := ctx.URLParamDefault("url", "")
	if u == "" {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Query param "url" cannot be empty!`,
			"data":    nil,
		})
		return
	}

	directUrl, err := ns.GetDirectUrl(u)
	// TODO 缓存direct url

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

func FetchImage(ctx iris.Context, ns NMacService) {
	u := ctx.URLParamDefault("url", "")
	if u == "" {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Query param "url" cannot be empty!`,
			"data":    nil,
		})
		return
	}

	client := ns.GetHttpClient()
	r, err := client.Get(u)
	if err != nil {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Fetch image failed!`,
			"data":    nil,
		})
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Read image data failed!`,
			"data":    nil,
		})
		return
	}

	// TODO 缓存图片

	ctx.ContentType(r.Header.Get("Content-Type"))
	ctx.Write(data)
}

// TODO 默认user-agent
