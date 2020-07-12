package main

import (
	"fmt"
	"github.com/prologic/bitcask"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"time"
)

type CacheService interface {
	GetDirectUrl(url string) (string, bool)
	SetDirectUrl(url string, directUrl string)
	GetImageCache(url string) (contentType string, data []byte, exists bool)
	SetImageCache(url string, contentType string, data []byte)
	Close()
}

type _CacheServiceImpl struct {
	imageCacheDir string
	db            *bitcask.Bitcask
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

func (_this *_CacheServiceImpl) newImageCacheFile() string {
	t := time.Now()
	return fmt.Sprintf("%04d_%02d_%02d__%s", t.Year(), t.Month(), t.Day(), _this.randString(32))
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

func (_this *_CacheServiceImpl) GetImageCache(url string) (string, []byte, bool) {
	key := _this.keyOfImageCacheFile(url)
	value, err := _this.db.Get(key)
	if err != nil {
		return "", nil, false
	}
	file := string(value)
	imageCachePath := path.Join(_this.imageCacheDir, file)
	imageCacheContentTypePath := imageCachePath + ".type"

	if !_this.fileExists(imageCachePath) || !_this.fileExists(imageCacheContentTypePath) {
		_ = _this.db.Delete(key)
		_ = os.Remove(imageCachePath)
		_ = os.Remove(imageCacheContentTypePath)
		return "", nil, false
	}

	contentType, err := ioutil.ReadFile(imageCacheContentTypePath)
	if err != nil {
		_ = _this.db.Delete(key)
		_ = os.Remove(imageCachePath)
		_ = os.Remove(imageCacheContentTypePath)
		return "", nil, false
	}

	data, err := ioutil.ReadFile(imageCachePath)
	if err != nil {
		_ = _this.db.Delete(key)
		_ = os.Remove(imageCachePath)
		_ = os.Remove(imageCacheContentTypePath)
		return "", nil, false
	}

	return string(contentType), data, true
}

func (_this *_CacheServiceImpl) SetImageCache(url string, contentType string, data []byte) {
	key := _this.keyOfImageCacheFile(url)
	var file string
	var imageCachePath string
	var imageCacheContentTypePath string
	for {
		file = _this.newImageCacheFile()
		imageCachePath = path.Join(_this.imageCacheDir, file)
		imageCacheContentTypePath = imageCachePath + ".type"
		if !_this.fileExists(imageCachePath) {
			break
		}
	}
	_ = ioutil.WriteFile(imageCachePath, data, 0644)
	_ = ioutil.WriteFile(imageCacheContentTypePath, []byte(contentType), 0644)
	_ = _this.db.Put(key, []byte(file))
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
