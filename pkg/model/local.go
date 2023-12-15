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

/*
type MatchParam struct {
	Field                string `yaml:"field" json:"field"`
	AutoGenerateSynonyms bool   `yaml:"auto_generate_synonyms_phrase_query" json:"auto_generate_synonyms"`
	Boost                int    `yaml:"boost" json:"boost"`
	FuzzyTranspositions  bool   `yaml:"fuzzy_transpositions" json:"fuzzy_transpositions"`
	Lenient              bool   `yaml:"lenient" json:"lenient"`
	MaxExpansions        int    `yaml:"max_expansions" json:"max_expansions"`
	Operator             string `yaml:"operator" json:"operator"`
	PrefixLength         int    `yaml:"prefix_length" json:"prefix_length"`
	Query                string `yaml:"query" json:"query"`
	ZeroTermsQuery       string `yaml:"zero_terms_query" json:"zero_terms_query"`
}*/

/*
type MatchParam struct {
	Field                string `yaml:"field" json:"field"`
	AutoGenerateSynonyms bool   `yaml:"auto_generate_synonyms_phrase_query" json:"auto_generate_synonyms"`
	Boost                int    `yaml:"boost" json:"boost"`
	FuzzyTranspositions  bool   `yaml:"fuzzy_transpositions" json:"fuzzy_transpositions"`
	Lenient              bool   `yaml:"lenient" json:"lenient"`
	MaxExpansions        int    `yaml:"max_expansions" json:"max_expansions"`
	Operator             string `yaml:"operator" json:"operator"`
	PrefixLength         int    `yaml:"prefix_length" json:"prefix_length"`
	Query                string `yaml:"query" json:"query"`
	ZeroTermsQuery       string `yaml:"zero_terms_query" json:"zero_terms_query"`
}*/

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
	Id        string    `yaml:"-" json:"id"`
	Name      string    `yaml:"name" json:"name"`
	Severity  string    `yaml:"severity" json:"severity"`
	Condition Condition `yaml:"condition" json:"condition"`
	Actions   []Action  `yaml:"actions" json:"actions"`
}

type Action struct {
	Name            string `yaml:"name" json:"name"`
	DestinationId   string `yaml:"destinationId" json:"destination_id"`
	Id              string `yaml:"-" json:"id"`
	DestinationName string `json:"destination_name,omitempty"`
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
