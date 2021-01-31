package main

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kataras/iris/v12"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ItemShortInfo struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	ImageUrl      string `json:"image_url"`
	DatePublished string `json:"date_published"`
	DetailPageUrl string `json:"detail_page_url"`
}

type DownloadUrl struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type ItemDetailInfo struct {
	Title           string         `json:"title"`
	Version         string         `json:"version"`
	Size            string         `json:"size"`
	DatePublished   string         `json:"date_published"`
	Content         string         `json:"content"`
	Urls            []*DownloadUrl `json:"urls"`
	PreviousPageUrl string         `json:"previous_page_url"`
}

type PreviousVersionInfo struct {
	Version string         `json:"version"`
	Urls    []*DownloadUrl `json:"urls"`
}

type CategoryInfo struct {
	Title    string `json:"title"`
	Category string `json:"category"`
}

type NMacService interface {
	AllowUrl(u string) bool
	UseImageCache() bool
	GetCategories() []*CategoryInfo
	GetList(category string, page int) (*iris.Map, error)
	Search(searchText string, page int) (*iris.Map, error)
	GetDetail(detailPageUrl string) (*ItemDetailInfo, error)
	GetDirectUrl(u string) (string, error)
	GetPreviousVersion(previousPageUrl string) []*PreviousVersionInfo
	FetchImage(u string) (contentType string, data []byte, err error)
}

type _NMacServiceImpl struct {
	proxy         *url.URL
	userAgent     string
	useImageCache bool
	categories    []*CategoryInfo
}

func (_this *_NMacServiceImpl) request(method string, u string, body io.Reader) (*http.Response, error) {
	transport := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if _this.proxy != nil {
		transport.Proxy = http.ProxyURL(_this.proxy)
	}
	client := http.Client{
		Transport: &transport,
		Timeout:   30 * time.Second,
	}
	req, _ := http.NewRequest(method, u, body)
	req.Header.Set("User-Agent", _this.userAgent)
	return client.Do(req)
}

func (_this *_NMacServiceImpl) refreshCategories() {
	ignoredCategories := []string{"apps", "get-vpn-subscription"}

	u := "https://nmac.to/"
	r, err := _this.request("GET", u, nil)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return
	}

	menu := make([]*CategoryInfo, 0)
	doc.Find("#main-menu-full li a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")
		title := strings.TrimSpace(selection.Text())
		category := ""
		if strings.HasSuffix(href, "/") {
			href = href[0 : len(href)-1]
		}
		if pos := strings.LastIndex(href, "/"); pos != -1 {
			category = href[pos+1:]
		}

		ignore := false
		for _, ic := range ignoredCategories {
			if ic == category {
				ignore = true
				break
			}
		}

		if !ignore {
			menu = append(menu, &CategoryInfo{
				Title:    title,
				Category: category,
			})
		}
	})

	doc.Find("#sub-menu-full li a").Each(func(i int, selection *goquery.Selection) {
		href := selection.AttrOr("href", "")
		title := strings.TrimSpace(selection.Text())
		category := ""
		if strings.HasSuffix(href, "/") {
			href = href[0 : len(href)-1]
		}
		if pos := strings.LastIndex(href, "/"); pos != -1 {
			category = href[pos+1:]
		}

		ignore := false
		for _, ic := range ignoredCategories {
			if ic == category {
				ignore = true
				break
			}
		}

		if !ignore {
			menu = append(menu, &CategoryInfo{
				Title:    title,
				Category: category,
			})
		}
	})

	_this.categories = menu
}

func imageCacheUrl(u string) string {
	return "/image_cache?url=" + url.QueryEscape(u)
}

func (_this *_NMacServiceImpl) parseContent(theContent *goquery.Selection) (html string, version string, size string, previousPageUrl string, urls []*DownloadUrl) {
	urls = make([]*DownloadUrl, 0)
	previousPageUrl = ""
	version = ""
	size = ""

	theContent.Find("div").Each(func(i int, selection *goquery.Selection) {
		selection.Find("a.btn-small").Each(func(j int, selection2 *goquery.Selection) {
			urls = append(urls, &DownloadUrl{
				Title: strings.TrimSpace(selection2.Text()),
				Url:   selection2.AttrOr("href", ""),
			})
		})
		selection.Find("a.btn-danger").Each(func(j int, selection2 *goquery.Selection) {
			previousPageUrl = selection2.AttrOr("href", "")
		})
		selection.Remove()
	})
	theContent.Find("span.label").Each(func(i int, selection *goquery.Selection) {
		text := strings.TrimSpace(selection.Text())
		if strings.HasPrefix(text, "Size") {
			size = strings.ReplaceAll(text, "Size", "")
			size = strings.ReplaceAll(size, "–", "")
			size = strings.TrimSpace(size)
			selection.Remove()
		}
	})
	theContent.Find("script").Each(func(i int, selection *goquery.Selection) {
		selection.Remove()
	})
	theContent.Find("noscript").Each(func(i int, selection *goquery.Selection) {
		selection.Remove()
	})
	theContent.Find("hr").Each(func(i int, selection *goquery.Selection) {
		selection.Remove()
	})
	theContent.Find("br.clearer").Each(func(i int, selection *goquery.Selection) {
		selection.Remove()
	})
	theContent.Find("p").Each(func(i int, selection *goquery.Selection) {
		text := strings.TrimSpace(selection.Text())
		if text == "" && len(selection.Contents().Nodes) == 0 { // remove empty <p>
			selection.Remove()
		} else if strings.HasPrefix(strings.ToLower(text), "version") {
			version = strings.ReplaceAll(strings.ToLower(text), "version", "")
			version = strings.ReplaceAll(version, ":", "")
			version = strings.TrimSpace(version)
			selection.Remove()
		}
	})
	theContent.Find("img.lazyload").Each(func(i int, selection *goquery.Selection) {
		if strings.ToLower(selection.AttrOr("alt", "")) == "download" {
			selection.Remove()
		} else {
			selection.SetAttr("src", selection.AttrOr("data-src", ""))
			selection.RemoveAttr("data-src")
			selection.RemoveAttr("data-srcset")
			selection.RemoveAttr("data-srcset")
			selection.RemoveAttr("data-sizes")
		}
	})
	var firstImg string
	theContent.Find("img").Each(func(i int, selection *goquery.Selection) {
		selection.RemoveAttr("class")
		selection.RemoveAttr("width")
		selection.RemoveAttr("height")
		if _this.UseImageCache() {
			u := selection.AttrOr("src", "")
			if _this.AllowUrl(u) {
				selection.SetAttr("src", imageCacheUrl(u))
			}
		}
		if i == 0 {
			selection.SetAttr("class", "detail-icon")
			firstImg, _ = goquery.OuterHtml(selection)
			selection.Remove()
		}
	})

	html, _ = theContent.Html()
	return strings.TrimSpace(firstImg + html), version, size, previousPageUrl, urls
}

func (_this *_NMacServiceImpl) parseListPage(u string) (iris.Map, error) {
	r, err := _this.request("GET", u, nil)

	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, err
	}
	list := make([]*ItemShortInfo, 0)
	doc.Find(".main-loop-inner>div").Each(func(i int, selection *goquery.Selection) {
		if len(selection.Find(".article-image-wrapper").Nodes) > 0 {
			title := selection.Find(".article-excerpt-wrapper .article-excerpt a").Text()
			desc := selection.Find(".article-excerpt-wrapper .article-excerpt .excerpt").Text()
			imgUrl := selection.Find(".article-image-wrapper img").AttrOr("data-src", "")
			if _this.UseImageCache() && _this.AllowUrl(imgUrl) {
				imgUrl = imageCacheUrl(imgUrl)
			}
			detailPageUrl := selection.Find(".article-image-wrapper a").AttrOr("href", "")
			datePublished := selection.Find(".article-meta .meta-info").Text()
			dateRegexp := regexp.MustCompile(`\d+-\d+-\d+`)
			matches := dateRegexp.FindStringSubmatch(datePublished)
			if len(matches) > 0 {
				datePublished = matches[0]
			} else {
				datePublished = ""
			}

			list = append(list, &ItemShortInfo{
				Title:         title,
				Description:   desc,
				ImageUrl:      imgUrl,
				DatePublished: datePublished,
				DetailPageUrl: detailPageUrl,
			})
		}
	})

	maxPage, err := strconv.ParseInt(strings.TrimSpace(doc.Find(".pagination-inner a").Last().AttrOr("data-paginated", "1")), 10, 32)
	if err != nil {
		maxPage = 1
	}

	return iris.Map{
		"max_page": maxPage,
		"length":   len(list),
		"list":     list,
	}, nil
}

func (_this *_NMacServiceImpl) AllowUrl(u string) bool {
	uu, err := url.Parse(u)
	if err != nil {
		return false
	}

	host := uu.Host
	if pos := strings.Index(host, ":"); pos != -1 {
		host = host[0:pos]
	}

	if host == "nmac.to" {
		return true
	}

	return false
}

func (_this *_NMacServiceImpl) UseImageCache() bool {
	return _this.useImageCache
}

func (_this *_NMacServiceImpl) GetCategories() []*CategoryInfo {
	return _this.categories
}

func (_this *_NMacServiceImpl) GetList(category string, page int) (*iris.Map, error) {
	u := "https://nmac.to/"
	if category != "" {
		if category == "tutorial" {
			u += "tutorial/"
		} else {
			u += fmt.Sprintf("category/%s/", category)
		}
	}
	if page > 1 {
		u += fmt.Sprintf("page/%d/", page)
	}

	m, err := _this.parseListPage(u)
	if err != nil {
		return nil, err
	}

	return &iris.Map{
		"category": category,
		"page":     page,
		"max_page": m["max_page"],
		"length":   m["length"],
		"list":     m["list"],
	}, nil
}

func (_this *_NMacServiceImpl) Search(searchText string, page int) (*iris.Map, error) {
	u := "https://nmac.to/"
	if page > 1 {
		u += fmt.Sprintf("page/%d/", page)
	}

	if searchText != "" {
		u += fmt.Sprintf("?s=%s", url.QueryEscape(searchText))
	}

	m, err := _this.parseListPage(u)
	if err != nil {
		return nil, err
	}

	return &iris.Map{
		"searchText": searchText,
		"page":       page,
		"max_page":   m["max_page"],
		"length":     m["length"],
		"list":       m["list"],
	}, nil
}

func (_this *_NMacServiceImpl) GetDetail(detailPageUrl string) (*ItemDetailInfo, error) {
	r, err := _this.request("GET", detailPageUrl, nil)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return nil, err
	}

	title := doc.Find(".main-content h1").Text()
	datePublished := doc.Find(".sortbar .date").Text()
	dateRegexp := regexp.MustCompile(`\d+-\d+-\d+`)
	matches := dateRegexp.FindStringSubmatch(datePublished)
	if len(matches) > 0 {
		datePublished = matches[0]
	} else {
		datePublished = ""
	}

	content, version, size, previousPageUrl, urls := _this.parseContent(doc.Find(".the-content"))
	if strings.HasPrefix(previousPageUrl, "/") {
		previousPageUrl = "https://nmac.to" + previousPageUrl
	}

	detail := &ItemDetailInfo{
		Title:           title,
		Version:         version,
		Size:            size,
		DatePublished:   datePublished,
		Content:         content,
		Urls:            urls,
		PreviousPageUrl: previousPageUrl,
	}

	return detail, nil
}

func (_this *_NMacServiceImpl) GetDirectUrl(u string) (string, error) {
	r, err := _this.request("GET", u, nil)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return "", err
	}
	return doc.Find("a.btn-block").AttrOr("href", ""), nil
}

func (_this *_NMacServiceImpl) GetPreviousVersion(previousPageUrl string) []*PreviousVersionInfo {
	versions := make([]*PreviousVersionInfo, 0)
	r, err := _this.request("GET", previousPageUrl, nil)
	if err != nil {
		return versions
	}
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		return versions
	}

	doc.Find(".accordion .accordion-group").Each(func(i int, selection *goquery.Selection) {
		version := strings.TrimSpace(selection.Find(".accordion-heading a").Text())
		version = strings.ReplaceAll(version, "Version", "")
		version = strings.ReplaceAll(version, "version", "")
		version = strings.ReplaceAll(version, ":", "")
		version = strings.TrimSpace(version)
		urls := make([]*DownloadUrl, 0)

		selection.Find(".accordion-inner a.btn-block").Each(func(j int, selection2 *goquery.Selection) {
			urls = append(urls, &DownloadUrl{
				Title: strings.TrimSpace(selection2.Text()),
				Url:   selection2.AttrOr("href", ""),
			})
		})

		versions = append(versions, &PreviousVersionInfo{
			Version: version,
			Urls:    urls,
		})
	})

	return versions
}

func (_this *_NMacServiceImpl) FetchImage(u string) (contentType string, data []byte, err error) {
	r, err := _this.request("GET", u, nil)
	if err != nil {
		return "", nil, err
	}

	data, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return "", nil, err
	}

	contentType = r.Header.Get("Content-Type")
	return contentType, data, nil
}

func NewNMacService(proxy string, userAgent string, useImageCache bool) NMacService {
	var proxyUrl *url.URL
	var err error
	if proxy != "" {
		proxyUrl, err = url.Parse(proxy)
		if err != nil {
			panic(err)
		}
	}
	ns := &_NMacServiceImpl{
		proxy:         proxyUrl,
		userAgent:     userAgent,
		useImageCache: useImageCache,
	}

	refreshCategoriesTicker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			<-refreshCategoriesTicker.C
			ns.refreshCategories()
		}
	}()
	ns.refreshCategories()
	return ns
}
