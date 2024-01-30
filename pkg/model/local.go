package model

type Monitor struct {
	ID       string    `yaml:"-" json:"-"`
	Type     string    `yaml:"type" json:"type"`
	Name     string    `yaml:"name" json:"name"`
	Enabled  bool      `yaml:"enabled" json:"enabled"`
	Schedule Schedule  `yaml:"schedule" json:"schedule"`
	Inputs   []Input   `yaml:"inputs" json:"inputs"`
	Triggers []Trigger `yaml:"triggers" json:"triggers"`
}

type Schedule struct {
	Period Period `yaml:"period" json:"period"`
}

type Period struct {
	Interval int    `yaml:"interval" json:"interval"`
	Unit     string `yaml:"unit" json:"unit"`
}

type Input struct {
	Search Search `yaml:"search" json:"search"`
}

type Search struct {
	Indices []string   `yaml:"indices" json:"indices"`
	Query   QueryParam `yaml:"query" json:"query"`
}

type QueryParam struct {
	Query InnerQuery `yaml:"query" json:"query"`
}

type InnerQuery struct {
	Bool BoolParam `yaml:"bool" json:"bool"`
}

type BoolParam struct {
	AdjustPureNegative bool        `yaml:"adjust_pure_negative" json:"adjust_pure_negative"`
	Boost              float64     `yaml:"boost" json:"boost"`
	Must               []MustParam `yaml:"must" json:"must"`
	MustNot            []MustParam `yaml:"must_not" json:"must_not"`
}

type MustParam struct {
	Match map[string]any `yaml:"match" json:"match,omitempty"`
	Range map[string]any `yaml:"range" json:"range,omitempty"`
}

type RangeParam struct {
	Field        string `yaml:"field" json:"field"`
	Boost        int    `yaml:"boost" json:"boost"`
	From         string `yaml:"from" json:"from"`
	IncludeLower bool   `yaml:"include_lower" json:"include_lower"`
	IncludeUpper bool   `yaml:"include_upper" json:"includeUpper"`
	TimeZone     string `yaml:"time_zone" json:"timeZone"`
	To           string `yaml:"to" json:"to"`
}

type Trigger struct {
	ID        string    `yaml:"-" json:"id"`
	Name      string    `yaml:"name" json:"name"`
	Severity  string    `yaml:"severity" json:"severity"`
	Condition Condition `yaml:"condition" json:"condition"`
	Actions   []Action  `yaml:"actions" json:"actions"`
}

type Action struct {
	Name            string `yaml:"name" json:"name"`
	DestinationName string `json:"destination_name,omitempty"`
	DestinationID   string `yaml:"destinationID" json:"destination_id"`
	SubjectTemplate Script `yaml:"subject" json:"subject_template"`
	MessageTemplate Script `yaml:"message" json:"message_template"`
}
type Condition struct {
	Script Script `yaml:"script" json:"script"`
}

type Script struct {
	Source string `json:"source"`
	Lang   string `json:"lang"`
}
