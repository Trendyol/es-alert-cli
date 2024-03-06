package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type BaseClient struct {
	client      *fasthttp.Client
	baseURL     string
	AuthOptions *BasicAuth
}

type BasicAuth struct {
	Username string
	Password string
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func NewBaseClient(baseURL string, auth *BasicAuth) *BaseClient {
	return &BaseClient{
		client:      new(fasthttp.Client),
		baseURL:     baseURL,
		AuthOptions: auth,
	}
}

func (b *BaseClient) Post(url string, pv interface{}) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	res := fasthttp.AcquireResponse()
	req.SetRequestURI(b.baseURL + url)
	req.Header.SetMethod("POST")
	if b.AuthOptions != nil {
		req.Header.Set("Authorization", "Basic "+basicAuth(b.AuthOptions.Username, b.AuthOptions.Password))
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
	if b.AuthOptions != nil {
		req.Header.Set("Authorization", "Basic "+basicAuth(b.AuthOptions.Username, b.AuthOptions.Password))
	}
	if pv != nil {
		body, err := json.Marshal(pv)
		if err != nil {
			return nil, err
		}
		req.SetBody(body)
		req.Header.SetContentType("application/json")
	}

	if err := b.client.Do(req, res); err != nil {
		return nil, err
	}
	err := getBody(res)
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
