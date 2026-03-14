package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/axelixlabs/axelix-cli/internal/models"
)

// APIError is returned for non-2xx HTTP responses.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Body)
}

// Client talks to an Axelix SBS-enabled Spring Boot application.
type Client struct {
	baseURL  string
	username string
	password string
	http     *http.Client
}

// New creates a new Client. baseURL is set to url + "/actuator".
func New(url, username, password string) *Client {
	return &Client{
		baseURL:  strings.TrimRight(url, "/") + "/actuator",
		username: username,
		password: password,
		http:     &http.Client{},
	}
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := c.baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if c.username != "" || c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}
	return req, nil
}

func (c *Client) doJSON(method, path string, bodyObj, out any) error {
	var bodyReader io.Reader
	if bodyObj != nil {
		data, err := json.Marshal(bodyObj)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := c.newRequest(method, path, bodyReader)
	if err != nil {
		return err
	}
	if bodyObj != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &APIError{StatusCode: resp.StatusCode, Body: strings.TrimSpace(string(respBody))}
	}
	if out != nil && len(respBody) > 0 {
		return json.Unmarshal(respBody, out)
	}
	return nil
}

func (c *Client) doDelete(path string) error {
	req, err := c.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return &APIError{StatusCode: resp.StatusCode, Body: strings.TrimSpace(string(body))}
	}
	return nil
}

func (c *Client) doRaw(path string) ([]byte, error) {
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{StatusCode: resp.StatusCode, Body: strings.TrimSpace(string(body))}
	}
	return body, nil
}

func (c *Client) doText(path string) (string, error) {
	data, err := c.doRaw(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetBeans fetches all Spring beans.
func (c *Client) GetBeans() (*models.BeansFeed, error) {
	out := &models.BeansFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-beans", nil, out)
}

// GetCaches fetches all cache managers and their caches.
func (c *Client) GetCaches() (*models.CachesFeed, error) {
	out := &models.CachesFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-caches", nil, out)
}

// GetSingleCache fetches a single cache by manager and cache name.
func (c *Client) GetSingleCache(manager, cache string) (json.RawMessage, error) {
	var out json.RawMessage
	return out, c.doJSON(http.MethodGet, "/axelix-caches/"+manager+"/"+cache, nil, &out)
}

// EnableCache enables a cache (or all caches in a manager if cache is empty).
func (c *Client) EnableCache(manager, cache string) error {
	if cache == "" {
		return c.doJSON(http.MethodPost, "/axelix-caches/"+manager+"/enable", nil, nil)
	}
	return c.doJSON(http.MethodPost, "/axelix-caches/"+manager+"/"+cache+"/enable", nil, nil)
}

// DisableCache disables a cache (or all caches in a manager if cache is empty).
func (c *Client) DisableCache(manager, cache string) error {
	if cache == "" {
		return c.doJSON(http.MethodPost, "/axelix-caches/"+manager+"/disable", nil, nil)
	}
	return c.doJSON(http.MethodPost, "/axelix-caches/"+manager+"/"+cache+"/disable", nil, nil)
}

// ClearCaches clears caches. Arguments narrow the scope; empty strings mean "all".
func (c *Client) ClearCaches(manager, cache, key string) error {
	path := "/axelix-caches"
	if manager != "" {
		path += "/" + manager
		if cache != "" {
			path += "/" + cache
			if key != "" {
				path += "/" + key
			}
		}
	}
	return c.doDelete(path)
}

// GetConditions fetches auto-configuration conditions.
func (c *Client) GetConditions() (*models.ConditionsFeed, error) {
	out := &models.ConditionsFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-conditions", nil, out)
}

// GetConfigProps fetches configuration properties.
func (c *Client) GetConfigProps() (*models.ConfigPropsFeed, error) {
	out := &models.ConfigPropsFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-configprops", nil, out)
}

// GetDetails fetches instance details.
func (c *Client) GetDetails() (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-details", nil, &out)
}

// GetEnv fetches environment properties, optionally filtered by pattern.
func (c *Client) GetEnv(pattern string) (*models.EnvironmentFeed, error) {
	path := "/axelix-env"
	if pattern != "" {
		path += "?pattern=" + pattern
	}
	out := &models.EnvironmentFeed{}
	return out, c.doJSON(http.MethodGet, path, nil, out)
}

// GetGCLogStatus fetches GC log status.
func (c *Client) GetGCLogStatus() (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-gc/log/status", nil, &out)
}

// GetGCLogFile retrieves the GC log file content.
func (c *Client) GetGCLogFile() (string, error) {
	return c.doText("/axelix-gc/log/file")
}

// TriggerGC triggers garbage collection.
func (c *Client) TriggerGC() error {
	return c.doJSON(http.MethodPost, "/axelix-gc/trigger", nil, nil)
}

// EnableGCLog enables GC logging at the given level.
func (c *Client) EnableGCLog(level string) error {
	body := map[string]string{"level": level}
	return c.doJSON(http.MethodPost, "/axelix-gc/log/enable", body, nil)
}

// DisableGCLog disables GC logging.
func (c *Client) DisableGCLog() error {
	return c.doJSON(http.MethodPost, "/axelix-gc/log/disable", nil, nil)
}

// DownloadHeapDump downloads a heap dump. If live is true, only live objects are included.
func (c *Client) DownloadHeapDump(live bool) ([]byte, error) {
	path := "/axelix-heap-dump"
	if live {
		path += "?live=true"
	}
	return c.doRaw(path)
}

// GetLoggers fetches all loggers.
func (c *Client) GetLoggers() (*models.ServiceLoggers, error) {
	out := &models.ServiceLoggers{}
	return out, c.doJSON(http.MethodGet, "/axelix-loggers", nil, out)
}

// GetLogger fetches a single logger by name.
func (c *Client) GetLogger(name string) (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-loggers/"+name, nil, &out)
}

// SetLogLevel sets the log level for the given logger.
func (c *Client) SetLogLevel(name, level string) error {
	body := map[string]string{"configuredLevel": level}
	return c.doJSON(http.MethodPost, "/axelix-loggers/"+name, body, nil)
}

// GetMetadata fetches instance metadata.
func (c *Client) GetMetadata() (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-metadata", nil, &out)
}

// GetMetrics fetches all metric groups.
func (c *Client) GetMetrics() (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-metrics", nil, &out)
}

// GetMetric fetches a single metric by name, optionally filtered by tag (e.g. "key:value").
func (c *Client) GetMetric(name, tag string) (map[string]any, error) {
	path := "/axelix-metrics/" + name
	if tag != "" {
		path += "?tag=" + tag
	}
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, path, nil, &out)
}

// GetScheduledTasks fetches all scheduled tasks.
func (c *Client) GetScheduledTasks() (*models.ServiceScheduledTasks, error) {
	out := &models.ServiceScheduledTasks{}
	return out, c.doJSON(http.MethodGet, "/axelix-scheduled-tasks", nil, out)
}

// EnableScheduledTask enables a scheduled task identified by trigger.
func (c *Client) EnableScheduledTask(trigger string, force bool) error {
	body := map[string]any{"trigger": trigger, "force": force}
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/enable", body, nil)
}

// DisableScheduledTask disables a scheduled task identified by trigger.
func (c *Client) DisableScheduledTask(trigger string, force bool) error {
	body := map[string]any{"trigger": trigger, "force": force}
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/disable", body, nil)
}

// ExecuteScheduledTask immediately executes a scheduled task.
func (c *Client) ExecuteScheduledTask(trigger string) error {
	body := map[string]string{"trigger": trigger}
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/execute", body, nil)
}

// SetTaskCronExpression modifies the cron expression of a scheduled task.
func (c *Client) SetTaskCronExpression(trigger, expr string) error {
	body := map[string]string{"trigger": trigger, "cronExpression": expr}
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/modify/cron-expression", body, nil)
}

// SetTaskInterval modifies the interval (in ms) of a scheduled task.
func (c *Client) SetTaskInterval(trigger string, intervalMs int64) error {
	body := map[string]any{"trigger": trigger, "interval": intervalMs}
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/modify/interval", body, nil)
}

// GetThreadDump fetches a thread dump.
func (c *Client) GetThreadDump() (*models.ThreadDumpFeed, error) {
	out := &models.ThreadDumpFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-thread-dump", nil, out)
}

// EnableThreadContention enables thread contention monitoring.
func (c *Client) EnableThreadContention() error {
	return c.doJSON(http.MethodPost, "/axelix-thread-dump/enable", nil, nil)
}

// DisableThreadContention disables thread contention monitoring.
func (c *Client) DisableThreadContention() error {
	return c.doJSON(http.MethodPost, "/axelix-thread-dump/disable", nil, nil)
}

// GetTransactions fetches all monitored transaction entrypoints.
func (c *Client) GetTransactions() (*models.TransactionMonitoringFeed, error) {
	out := &models.TransactionMonitoringFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-transactions-monitoring", nil, out)
}

// ClearTransactions clears all recorded transaction data.
func (c *Client) ClearTransactions() error {
	return c.doDelete("/axelix-transactions-monitoring")
}
