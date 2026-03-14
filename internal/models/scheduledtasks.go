package models

// ServiceScheduledTasks is the response of the scheduled-tasks actuator endpoint.
type ServiceScheduledTasks struct {
	Cron        []CronTask        `json:"cron"`
	FixedDelay  []FixedDelayTask  `json:"fixedDelay"`
	FixedRate   []FixedRateTask   `json:"fixedRate"`
	Custom      []CustomTask      `json:"custom"`
}

// TaskRunnable holds the target of a scheduled task.
type TaskRunnable struct {
	Target string `json:"target"`
}

// NextExecution holds the next planned execution time.
type NextExecution struct {
	Time string `json:"time"`
}

// LastExecution holds the last execution result.
type LastExecution struct {
	Status    string             `json:"status"`
	Time      string             `json:"time"`
	Exception *LastExecException `json:"exception"`
}

// LastExecException holds exception details from a task execution.
type LastExecException struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// CronTask represents a cron-scheduled task.
type CronTask struct {
	Runnable      TaskRunnable   `json:"runnable"`
	Expression    string         `json:"expression"`
	NextExecution *NextExecution `json:"nextExecution"`
	LastExecution *LastExecution `json:"lastExecution"`
	Enabled       bool           `json:"enabled"`
}

// FixedDelayTask represents a fixed-delay scheduled task.
type FixedDelayTask struct {
	Runnable      TaskRunnable   `json:"runnable"`
	Interval      float64        `json:"interval"`
	InitialDelay  float64        `json:"initialDelay"`
	NextExecution *NextExecution `json:"nextExecution"`
	LastExecution *LastExecution `json:"lastExecution"`
	Enabled       bool           `json:"enabled"`
}

// FixedRateTask represents a fixed-rate scheduled task.
type FixedRateTask struct {
	Runnable      TaskRunnable   `json:"runnable"`
	Interval      float64        `json:"interval"`
	InitialDelay  float64        `json:"initialDelay"`
	NextExecution *NextExecution `json:"nextExecution"`
	LastExecution *LastExecution `json:"lastExecution"`
	Enabled       bool           `json:"enabled"`
}

// CustomTask represents a task with a user-defined trigger.
type CustomTask struct {
	Runnable      TaskRunnable   `json:"runnable"`
	Trigger       string         `json:"trigger"`
	NextExecution *NextExecution `json:"nextExecution"`
	LastExecution *LastExecution `json:"lastExecution"`
	Enabled       bool           `json:"enabled"`
}
