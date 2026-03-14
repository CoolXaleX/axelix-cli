package models

// GcLogStatusResponse is the response of the GC log status endpoint.
type GcLogStatusResponse struct {
	Enabled         bool     `json:"enabled"`
	Level           *string  `json:"level"`
	AvailableLevels []string `json:"availableLevels"`
}
