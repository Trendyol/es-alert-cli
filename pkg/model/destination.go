package model

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

// TODO: simplify model
type ElasticFetchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index       string  `json:"_index"`
			Type        string  `json:"_type"`
			ID          string  `json:"_id"`
			Version     int     `json:"_version"`
			SeqNo       int     `json:"_seq_no"`
			PrimaryTerm int     `json:"_primary_term"`
			Score       float64 `json:"_score"`
			Source      struct {
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
	Version     int     `json:"_version"`
	SeqNo       int     `json:"_seq_no"`
	PrimaryTerm int     `json:"_primary_term"`
	Monitor     Monitor `json:"monitor"`
}
