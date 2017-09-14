package fit

import (
	"net/http"
)

type Cookie map[string]*http.Cookie

func (ck Cookie) Get(key string) *http.Cookie {
	if v, ok := ck[key]; ok {
		return v
	}
	return nil
}

func (ck Cookie) Set(key string, value string) {
	ck[key] = &http.Cookie{
		Name:  key,
		Value: value,
		Path:  "/",
	}
}
func (ck Cookie) SetWithPath(key, value, path string) {
	ck[key] = &http.Cookie{
		Name:  key,
		Value: value,
		Path:  path,
	}
}
func (ck Cookie) Del(key string) {
	delete(ck, key)
}

func (ck Cookie) SetHttpCookie(httpCookit *http.Cookie) {
	ck[httpCookit.Name] = httpCookit
}

func (ck Cookie) SetResponseCookie(w http.ResponseWriter) {
	for _, v := range ck {
		http.SetCookie(w, v)
	}
}
