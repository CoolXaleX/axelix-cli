package models

// MetricsGroupsFeed is the response of the metrics list endpoint.
type MetricsGroupsFeed struct {
	MetricsGroups []MetricsGroup `json:"metricsGroups"`
}

// MetricsGroup groups related metrics together.
type MetricsGroup struct {
	GroupName string              `json:"groupName"`
	Metrics   []MetricDescription `json:"metrics"`
}

// MetricDescription describes a single metric.
type MetricDescription struct {
	MetricName  string `json:"metricName"`
	Description string `json:"description"`
}
