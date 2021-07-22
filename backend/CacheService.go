package main

import (
	"encoding/json"
	"fmt"
	"git.mills.io/prologic/bitcask"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"
)

type CacheService interface {
	GetDirectUrl(url string) (directUrl string, exists bool)
	SetDirectUrl(url string, directUrl string)
	GetImageCache(url string) (contentType string, cachePath string, exists bool)
	SetImageCache(url string, contentType string, data []byte)
	Close()
}

type _CacheServiceImpl struct {
	imageCacheDir string
	db            *bitcask.Bitcask
}

type ImageCacheInfo struct {
	File        string `json:"file"`
	ContentType string `json:"content_type"`
}

func (_this *_CacheServiceImpl) keyOfDirectUrl(u string) []byte {
	return []byte(fmt.Sprintf("direct_url_of__%s", u))
}

func (_this *_CacheServiceImpl) keyOfImageCacheFile(u string) []byte {
	return []byte(fmt.Sprintf("image_cache_file_of__%s", u))
}

func (_this *_CacheServiceImpl) randString(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (_this *_CacheServiceImpl) fileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (_this *_CacheServiceImpl) newImageCacheFile() (cacheFile string, cachePath string) {
	for {
		t := time.Now()
		cacheFile = fmt.Sprintf("%04d_%02d_%02d__%s", t.Year(), t.Month(), t.Day(), _this.randString(32))
		cachePath = path.Join(_this.imageCacheDir, cacheFile)
		if !_this.fileExists(cachePath) {
			return cacheFile, cachePath
		}
	}
}

func (_this *_CacheServiceImpl) GetDirectUrl(url string) (string, bool) {
	key := _this.keyOfDirectUrl(url)
	value, err := _this.db.Get(key)
	if err != nil {
		return "", false
	}
	return string(value), true
}

func (_this *_CacheServiceImpl) SetDirectUrl(url string, directUrl string) {
	key := _this.keyOfDirectUrl(url)
	_ = _this.db.Put(key, []byte(directUrl))
}

func (_this *_CacheServiceImpl) GetImageCache(url string) (string, string, bool) {
	key := _this.keyOfImageCacheFile(url)
	value, err := _this.db.Get(key)
	if err != nil {
		return "", "", false
	}

	var info ImageCacheInfo
	_ = json.Unmarshal(value, &info)

	cachePath := path.Join(_this.imageCacheDir, info.File)
	contentType := info.ContentType

	if !_this.fileExists(cachePath) {
		_ = _this.db.Delete(key)
		_ = os.Remove(cachePath)
		return "", "", false
	}

	return contentType, cachePath, true
}

func (_this *_CacheServiceImpl) SetImageCache(url string, contentType string, data []byte) {
	key := _this.keyOfImageCacheFile(url)
	cacheFile, cachePath := _this.newImageCacheFile()

	info, _ := json.Marshal(ImageCacheInfo{
		File:        cacheFile,
		ContentType: contentType,
	})

	_ = ioutil.WriteFile(cachePath, data, 0644)
	_ = _this.db.Put(key, info)
}

func (_this *_CacheServiceImpl) Close() {
	if _this.db != nil {
		_ = _this.db.Close()
	}
}

func NewCacheService(maxDataSize int, cacheDbDir string, imageCacheDir string) CacheService {
	db, err := bitcask.Open(
		cacheDbDir,
		bitcask.WithMaxKeySize(8192),
		bitcask.WithMaxDatafileSize(maxDataSize),
	)
	if err != nil {
		panic(err)
	}

	return &_CacheServiceImpl{
		db:            db,
		imageCacheDir: imageCacheDir,
	}
}
