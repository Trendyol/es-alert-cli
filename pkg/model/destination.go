package model

type CustomWebhook struct {
	HeaderParams map[string]string `json:"header_params,omitempty" yaml:",omitempty"`
	QueryParams  map[string]string `json:"query_params,omitempty" yaml:",omitempty"`
	Path         string            `json:"path,omitempty" yaml:",omitempty"`
	Password     string            `json:"password,omitempty" yaml:",omitempty"`
	Scheme       string            `json:"scheme,omitempty" yaml:",omitempty"`
	Host         string            `json:"host,omitempty" yaml:",omitempty"`
	URL          string            `json:"url,omitempty" yaml:",omitempty"`
	Username     string            `json:"username,omitempty" yaml:",omitempty"`
	Port         int               `json:"port,omitempty" yaml:",omitempty"`
}

type ElasticFetchResponse struct {
	Hits struct {
		Hits []struct {
			ID     string `json:"_id"`
			Source struct {
				Destination Destination `json:"destination,omitempty"`
				Monitor     Monitor     `json:"monitor,omitempty"`
			} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

type Destination struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Type          string        `json:"type"`
	Slack         Slack         `json:"slack,omitempty" yaml:",omitempty"`
	CustomWebhook CustomWebhook `json:"custom_webhook,omitempty" yaml:",omitempty"`
}

type Slack struct {
	URL string `json:"url,omitempty" yaml:",omitempty"`
}

type UpdateMonitorResponse struct {
	ID          string  `json:"_id"`
	Monitor     Monitor `json:"monitor"`
	Version     int     `json:"_version"`
	SeqNo       int     `json:"_seq_no"`
	PrimaryTerm int     `json:"_primary_term"`
}
