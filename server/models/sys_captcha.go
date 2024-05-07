package models

import (
	"dcss/global"
	"github.com/patrickmn/go-cache"
)

type RespCaptcha struct {
	CaptchaID   string `json:"captcha_id"`
	CaptchaPath string `json:"captcha_path"`
}

type LocalCacheStore struct{}

func (l *LocalCacheStore) Set(id string, value string) error {
	global.CaptchaCache.Set(id, value, cache.DefaultExpiration)

	return nil
}

func (l *LocalCacheStore) Get(id string, clear bool) string {
	get, b := global.CaptchaCache.Get(id)
	if !b {
		return ""
	}
	if clear {
		defer global.CaptchaCache.Delete(id)
	}

	s, ok := get.(string)
	if ok {
		return s
	}
	return ""
}

func (l *LocalCacheStore) Verify(id string, answer string, clear bool) bool {
	get, b := global.CaptchaCache.Get(id)
	if !b {
		return b
	}

	if clear {
		defer global.CaptchaCache.Delete(id)
	}

	s, ok := get.(string)
	if !ok {
		return ok
	}

	if s == answer {
		return true
	}

	return false
}
