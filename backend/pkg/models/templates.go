package models

type MarketTemplate struct {
	ID               string `json:"id"`
	TitlePattern     string `json:"title_pattern"`
	ResolutionMethod string `json:"resolution_method"`
	DefaultRules     string `json:"default_rules"`
}

var Templates = []MarketTemplate{
	{
		ID:               "traffic_count",
		TitlePattern:     "Will %d or more cars cross the line in the next %d seconds?",
		ResolutionMethod: "traffic_count",
		DefaultRules:     `{"targetCount": 20, "window": 60}`,
	},
	{
		ID:               "pedestrian_count",
		TitlePattern:     "Will %d or more pedestrians be detected in the next %d seconds?",
		ResolutionMethod: "pedestrian_count",
		DefaultRules:     `{"targetCount": 5, "window": 60}`,
	},
}
