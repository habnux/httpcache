package httpcache

import (
	cacheControl "github.com/bxcodec/httpcache/helper/cacheheader"
	"net/http"
	"time"
)

type Validator interface {
	ValidateCacheControl(req *http.Request, resp *http.Response) (validationResult cacheControl.ObjectResults, err error)
}

func NewDefaultValidator() Validator {
	return &defaultValidator{}
}

type defaultValidator struct{}

func (d *defaultValidator) ValidateCacheControl(req *http.Request, resp *http.Response) (validationResult cacheControl.ObjectResults, err error) {
	reqDir, err := cacheControl.ParseRequestCacheControl(req.Header.Get("Cache-Control"))
	if err != nil {
		return
	}

	resDir, err := cacheControl.ParseResponseCacheControl(resp.Header.Get("Cache-Control"))
	if err != nil {
		return
	}

	expiry := resp.Header.Get("Expires")
	expiresHeader, err := http.ParseTime(expiry)
	if err != nil && expiry != "" &&
		// https://stackoverflow.com/questions/11357430/http-expires-header-values-0-and-1
		expiry != "-1" && expiry != "0" {
		return
	}

	dateHeaderStr := resp.Header.Get("Date")
	dateHeader, err := http.ParseTime(dateHeaderStr)
	if err != nil && dateHeaderStr != "" {
		return
	}

	lastModifiedStr := resp.Header.Get("Last-Modified")
	lastModifiedHeader, err := http.ParseTime(lastModifiedStr)
	if err != nil && lastModifiedStr != "" {
		return
	}

	obj := cacheControl.Object{
		RespDirectives:         resDir,
		RespHeaders:            resp.Header,
		RespStatusCode:         resp.StatusCode,
		RespExpiresHeader:      expiresHeader,
		RespDateHeader:         dateHeader,
		RespLastModifiedHeader: lastModifiedHeader,
		ReqDirectives:          reqDir,
		ReqHeaders:             req.Header,
		ReqMethod:              req.Method,
		NowUTC:                 time.Now().UTC(),
	}

	validationResult = cacheControl.ObjectResults{}
	cacheControl.CachableObject(&obj, &validationResult)
	cacheControl.ExpirationObject(&obj, &validationResult)

	return validationResult, nil

}
