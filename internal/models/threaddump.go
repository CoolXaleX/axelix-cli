package models

// ThreadDumpFeed is the response of the thread-dump actuator endpoint.
type ThreadDumpFeed struct {
	ThreadContentionMonitoringEnabled bool         `json:"threadContentionMonitoringEnabled"`
	Threads                           []ThreadInfo `json:"threads"`
}

// ThreadInfo holds information about a single thread.
type ThreadInfo struct {
	ThreadName  string  `json:"threadName"`
	ThreadID    int64   `json:"threadId"`
	BlockedTime int64   `json:"blockedTime"`
	BlockedCount int64  `json:"blockedCount"`
	WaitedTime  int64   `json:"waitedTime"`
	WaitedCount int64   `json:"waitedCount"`
	Daemon      bool    `json:"daemon"`
	Suspended   bool    `json:"suspended"`
	ThreadState string  `json:"threadState"`
	Priority    int     `json:"priority"`
}
