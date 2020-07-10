package main

import (
	"fmt"
	"github.com/kataras/golog"
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

func escapeUrl(s string) (result string) {
	for _, c := range s {
		if c <= 0x7f { // single byte
			result += fmt.Sprintf("%X", c)
		} else if c > 0x1fffff { // quaternary byte
			result += fmt.Sprintf("%X%X%X%X",
				0xf0+((c&0x1c0000)>>18),
				0x80+((c&0x3f000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else if c > 0x7ff { // triple byte
			result += fmt.Sprintf("%X%X%X",
				0xe0+((c&0xf000)>>12),
				0x80+((c&0xfc0)>>6),
				0x80+(c&0x3f),
			)
		} else { // double byte
			result += fmt.Sprintf("%X%X",
				0xc0+((c&0x7c0)>>6),
				0x80+(c&0x3f),
			)
		}
	}

	return result
}

func DirectUrl(ctx iris.Context, ns NMacService, cache CacheService, logger *golog.Logger) {
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

func FetchImage(ctx iris.Context, ns NMacService, cache CacheService, logger *golog.Logger) {
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

	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		ctx.JSON(iris.Map{
			"code":    1,
			"message": `Read image data failed!`,
			"data":    nil,
		})
		return
	}

	contentType = r.Header.Get("Content-Type")

	cache.SetImageCache(u, contentType, data)

	ctx.ContentType(contentType)
	ctx.Write(data)
}

// TODO 默认user-agent
