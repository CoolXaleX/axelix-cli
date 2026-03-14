package models

// EnvironmentFeed is the response of the env actuator endpoint.
type EnvironmentFeed struct {
	ActiveProfiles  []string         `json:"activeProfiles"`
	DefaultProfiles []string         `json:"defaultProfiles"`
	PropertySources []PropertySource `json:"propertySources"`
}

// PropertySource represents a single property source.
type PropertySource struct {
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Properties  []EnvProperty `json:"properties"`
}

// EnvProperty represents a single property entry.
type EnvProperty struct {
	Name   string  `json:"name"`
	Value  *string `json:"value"`
	IsPrimary bool `json:"isPrimary"`
}
