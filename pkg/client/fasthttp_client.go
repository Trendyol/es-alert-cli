package client

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
	"strings"
)

type BaseClient struct {
	client  *fasthttp.Client
	baseUrl string
}

func NewBaseClient(baseUrl string) *BaseClient {
	return &BaseClient{
		client:  new(fasthttp.Client),
		baseUrl: baseUrl,
	}
}
func (b *BaseClient) GET(url string, token ...string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	req.SetRequestURI(b.baseUrl + url)
	req.Header.SetMethod("GET")
	if len(token) > 0 && strings.TrimSpace(token[0]) != "" {
		req.Header.Set("Authorization", token[0])
	}
	req.SetRequestURIBytes(req.RequestURI())
	if err := b.client.Do(req, res); err != nil {
		return nil, err
	}
	err := b.getBody(res)
	if err != nil {
		return nil, err
	}
	return res, err
}
func (b *BaseClient) POST(url string, pv interface{}, opts ...map[string]string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	req.SetRequestURI(b.baseUrl + url)
	req.Header.SetMethod("POST")
	for _, opt := range opts {
		b.setOptsHeader(opt, req)
	}
	body, err := json.Marshal(pv)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)
	req.Header.SetContentType("application/json")
	if err := b.client.Do(req, res); err != nil {
		return nil, err
	}
	err = b.getBody(res)
	if err != nil {
		return nil, err
	}
	return res, err
}
func (b *BaseClient) PUT(url string, pv interface{}, opts ...map[string]string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	res := fasthttp.AcquireResponse()
	req.SetRequestURI(b.baseUrl + url)
	req.Header.SetMethod("PUT")
	for _, opt := range opts {
		b.setOptsHeader(opt, req)
	}
	body, err := json.Marshal(pv)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)
	req.Header.SetContentType("application/json")
	if err := b.client.Do(req, res); err != nil {
		return nil, err
	}
	err = b.getBody(res)
	if err != nil {
		return nil, err
	}
	return res, err
}
func (b BaseClient) getBody(res *fasthttp.Response) error {
	contentEncoding := res.Header.Peek("Content-Encoding")
	switch {
	case bytes.EqualFold(contentEncoding, []byte("gzip")):
		body, err := res.BodyGunzip()
		if err != nil {
			return err
		}
		res.SetBody(body)
		return nil
	case bytes.EqualFold(contentEncoding, []byte("brotli")):
		body, err := res.BodyUnbrotli()
		if err != nil {
			return err
		}
		res.SetBody(body)
		return nil
	default:
		return nil
	}
}
func (b BaseClient) Bind(body []byte, rv interface{}) error {
	if len(body) > 0 || body != nil {
		if err := json.Unmarshal(body, &rv); err != nil {
			return err
		}
	}
	return nil
}
func (b BaseClient) setOptsHeader(opts map[string]string, req *fasthttp.Request) {
	if opts != nil {
		for k, v := range opts {
			req.Header.Set(k, v)
		}
	}
}
