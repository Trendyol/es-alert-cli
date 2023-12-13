package model

type Monitor struct {
	Id       string    `yaml:"-" json:"-"`
	Type     string    `yaml:"type" json:"type"`
	Name     string    `yaml:"name" json:"name"`
	Enabled  bool      `yaml:"enabled" json:"enabled"`
	Schedule Schedule  `yaml:"schedule" json:"schedule"`
	Inputs   []Input   `yaml:"inputs" json:"inputs"`
	Triggers []Trigger `yaml:"triggers" json:"triggers"`
}

type Schedule struct {
	Period Period      `yaml:"period" json:"period"`
	Cron   interface{} `yaml:"cron" json:"cron"`
}

type Period struct {
	Interval int    `yaml:"interval" json:"interval"`
	Unit     string `yaml:"unit" json:"unit"`
}

type Input struct {
	Search Search `yaml:"search" json:"search"`
}

type Search struct {
	Indices []string   `yaml:"indices"`
	Query   QueryParam `yaml:"query"`
}

type QueryParam struct {
	Query BoolParam `yaml:"query"`
}

type BoolParam struct {
	AdjustPureNegative bool        `yaml:"adjust_pure_negative"`
	Boost              int         `yaml:"boost"`
	Must               []MustParam `yaml:"must"`
	MustNot            []MustParam `yaml:"must_not"`
}

type MustParam struct {
	Match MatchParam `yaml:"match"`
	Range RangeParam `yaml:"range"`
}

type MatchParam struct {
	Field                string `yaml:"field"`
	AutoGenerateSynonyms bool   `yaml:"auto_generate_synonyms_phrase_query"`
	Boost                int    `yaml:"boost"`
	FuzzyTranspositions  bool   `yaml:"fuzzy_transpositions"`
	Lenient              bool   `yaml:"lenient"`
	MaxExpansions        int    `yaml:"max_expansions"`
	Operator             string `yaml:"operator"`
	PrefixLength         int    `yaml:"prefix_length"`
	Query                string `yaml:"query"`
	ZeroTermsQuery       string `yaml:"zero_terms_query"`
}

type RangeParam struct {
	Field        string `yaml:"field"`
	Boost        int    `yaml:"boost"`
	From         string `yaml:"from"`
	IncludeLower bool   `yaml:"include_lower"`
	IncludeUpper bool   `yaml:"include_upper"`
	TimeZone     string `yaml:"time_zone"`
	To           string `yaml:"to"`
}

type Trigger struct {
	Name      string    `yaml:"name" json:"name"`
	Severity  string    `yaml:"severity" json:"severity"`
	Condition Condition `yaml:"condition" json:"condition"`
	Actions   []Action  `yaml:"actions" json:"actions"`
}

type Action struct {
	Name            string `yaml:"name"`
	DestinationId   string `yaml:"destinationId"`
	DestinationName string
	SubjectTemplate Script `yaml:"subject" json:"subject_template"`
	MessageTemplate Script `yaml:"message" json:"message_template"`
}
type Condition struct {
	Script Script `json:"script" yaml:"script"`
}

type Script struct {
	Source string `json:"source"`
	Lang   string `json:"lang"`
}
