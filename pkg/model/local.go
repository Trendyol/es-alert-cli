package model

type Monitor struct {
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

type InputConfig struct {
	Search SearchConfig `yaml:"search"`
}

type SearchConfig struct {
	Indices []string    `yaml:"indices"`
	Query   QueryConfig `yaml:"query"`
}

type QueryConfig struct {
	Query BoolConfig `yaml:"query"`
}

type BoolConfig struct {
	AdjustPureNegative bool         `yaml:"adjust_pure_negative"`
	Boost              int          `yaml:"boost"`
	Must               []MustConfig `yaml:"must"`
	MustNot            []MustConfig `yaml:"must_not"`
}

type MustConfig struct {
	Match MatchConfig `yaml:"match"`
	Range RangeConfig `yaml:"range"`
}

type MatchConfig struct {
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

type RangeConfig struct {
	Field        string `yaml:"field"`
	Boost        int    `yaml:"boost"`
	From         string `yaml:"from"`
	IncludeLower bool   `yaml:"include_lower"`
	IncludeUpper bool   `yaml:"include_upper"`
	TimeZone     string `yaml:"time_zone"`
	To           string `yaml:"to"`
}

type TriggerConfig struct {
	Name      string         `yaml:"name"`
	Severity  string         `yaml:"severity"`
	Condition string         `yaml:"condition"`
	Actions   []ActionConfig `yaml:"actions"`
}

type ActionConfig struct {
	Name          string `yaml:"name"`
	DestinationId string `yaml:"destinationId"`
	Subject       string `yaml:"subject"`
	Message       string `yaml:"message"`
}
