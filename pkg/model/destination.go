package model

type Destination struct {
	ID            string
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	Slack         Slack         `json:"slack,omitempty" yaml:",omitempty"`
	CustomWebhook CustomWebhook `json:"custom_webhook,omitempty" yaml:",omitempty"`
}

type Slack struct {
	URL string `json:"url,omitempty" yaml:",omitempty"`
}

type CustomWebhook struct {
	Path         string            `json:"path,omitempty" yaml:",omitempty"`
	HeaderParams map[string]string `json:"header_params,omitempty" yaml:",omitempty"`
	Password     string            `json:"password,omitempty" yaml:",omitempty"`
	Port         int               `json:"port,omitempty" yaml:",omitempty"`
	Scheme       string            `json:"scheme,omitempty" yaml:",omitempty"`
	QueryParams  map[string]string `json:"query_params,omitempty" yaml:",omitempty"`
	Host         string            `json:"host,omitempty" yaml:",omitempty"`
	URL          string            `json:"url,omitempty" yaml:",omitempty"`
	Username     string            `json:"username,omitempty" yaml:",omitempty"`
}
