package client

import (
	"bytes"
	"encoding/json"
	"github.com/valyala/fasthttp"
)

type BaseClient struct {
	client  *fasthttp.Client
	baseURL string
}

func NewBaseClient(baseURL string) *BaseClient {
	return &BaseClient{
		client:  new(fasthttp.Client),
		baseURL: baseURL,
	}
}

func (b *BaseClient) Post(url string, pv interface{}) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	req.SetRequestURI(b.baseURL + url)
	req.Header.SetMethod("POST")

	body, err := json.Marshal(pv)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)
	req.Header.SetContentType("application/json")

	if err := b.client.Do(req, res); err != nil {
		return nil, err
	}
	err = getBody(res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (b *BaseClient) Put(url string, pv interface{}) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	req.SetRequestURI(b.baseURL + url)
	req.Header.SetMethod("PUT")

	body, err := json.Marshal(pv)
	if err != nil {
		return nil, err
	}
	req.SetBody(body)
	req.Header.SetContentType("application/json")

	if err := b.client.Do(req, res); err != nil {
		return nil, err
	}
	err = getBody(res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func getBody(res *fasthttp.Response) error {
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

func (b *BaseClient) Bind(body []byte, rv interface{}) error {
	if len(body) > 0 || body != nil {
		if err := json.Unmarshal(body, &rv); err != nil {
			return err
		}
	}
	return nil
}
