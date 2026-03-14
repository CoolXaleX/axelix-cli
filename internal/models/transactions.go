package models

// TransactionMonitoringFeed is the response of the transactions-monitoring actuator endpoint.
type TransactionMonitoringFeed struct {
	Entrypoints []TransactionalEntrypoint `json:"entrypoints"`
}

// TransactionalEntrypoint represents a monitored transactional method.
type TransactionalEntrypoint struct {
	ClassName      string           `json:"className"`
	MethodName     string           `json:"methodName"`
	Executions     []TransactionExecution `json:"executions"`
	ExecutionStats ExecutionStats   `json:"executionStats"`
}

// TransactionExecution holds a single transaction execution record.
type TransactionExecution struct {
	DurationMs int64 `json:"durationMs"`
	Timestamp  int64 `json:"timestamp"`
}

// ExecutionStats holds aggregated statistics for a transactional entrypoint.
type ExecutionStats struct {
	AverageDurationMs int64 `json:"averageDurationMs"`
	MaxDurationMs     int64 `json:"maxDurationMs"`
	MedianDurationMs  int64 `json:"medianDurationMs"`
}
