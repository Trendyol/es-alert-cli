package model

type Monitor struct {
	Schedule Schedule  `yaml:"schedule" json:"schedule"`
	ID       string    `yaml:"-" json:"-"`
	Type     string    `yaml:"type" json:"type"`
	Name     string    `yaml:"name" json:"name"`
	Inputs   []Input   `yaml:"inputs" json:"inputs"`
	Triggers []Trigger `yaml:"triggers" json:"triggers"`
	Enabled  bool      `yaml:"enabled" json:"enabled"`
}

type Schedule struct {
	Period Period `yaml:"period" json:"period"`
}

type Period struct {
	Unit     string `yaml:"unit" json:"unit"`
	Interval int    `yaml:"interval" json:"interval"`
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
	Must               []MustParam `yaml:"must" json:"must"`
	MustNot            []MustParam `yaml:"must_not" json:"must_not"`
	Boost              float64     `yaml:"boost" json:"boost"`
	AdjustPureNegative bool        `yaml:"adjust_pure_negative" json:"adjust_pure_negative"`
}

type MustParam struct {
	Match map[string]any `yaml:"match" json:"match,omitempty"`
	Range map[string]any `yaml:"range" json:"range,omitempty"`
}

type RangeParam struct {
	Field        string `yaml:"field" json:"field"`
	From         string `yaml:"from" json:"from"`
	TimeZone     string `yaml:"time_zone" json:"timeZone"`
	To           string `yaml:"to" json:"to"`
	Boost        int    `yaml:"boost" json:"boost"`
	IncludeLower bool   `yaml:"include_lower" json:"include_lower"`
	IncludeUpper bool   `yaml:"include_upper" json:"includeUpper"`
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
