package models

// ServiceLoggers is the response of the loggers actuator endpoint.
type ServiceLoggers struct {
	Levels  []string                `json:"levels"`
	Loggers map[string]LoggerLevels `json:"loggers"`
	Groups  map[string]LoggerGroup  `json:"groups"`
}

// LoggerLevels holds configured and effective log levels for a logger.
type LoggerLevels struct {
	ConfiguredLevel *string `json:"configuredLevel"`
	EffectiveLevel  string  `json:"effectiveLevel"`
}

// LoggerGroup holds the members of a logger group.
type LoggerGroup struct {
	ConfiguredLevel *string  `json:"configuredLevel"`
	EffectiveLevel  string   `json:"effectiveLevel"`
	Members         []string `json:"members"`
}
