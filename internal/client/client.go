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
	baseURL string
	http    *http.Client
}

// New creates a new Client. baseURL is set to url + "/actuator".
func New(url string) *Client {
	return &Client{
		baseURL: strings.TrimRight(url, "/") + "/actuator",
		http:    &http.Client{},
	}
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
	req, err := http.NewRequest(method, c.baseURL+path, bodyReader)
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
	req, err := http.NewRequest(http.MethodDelete, c.baseURL+path, nil)
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
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
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

func (c *Client) GetBeans() (*models.BeansFeed, error) {
	out := &models.BeansFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-beans", nil, out)
}

func (c *Client) GetCaches() (*models.CachesFeed, error) {
	out := &models.CachesFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-caches", nil, out)
}

func (c *Client) GetSingleCache(manager, cache string) (json.RawMessage, error) {
	var out json.RawMessage
	return out, c.doJSON(http.MethodGet, "/axelix-caches/"+manager+"/"+cache, nil, &out)
}

func (c *Client) EnableCache(manager, cache string) error {
	if cache == "" {
		return c.doJSON(http.MethodPost, "/axelix-caches/"+manager+"/enable", nil, nil)
	}
	return c.doJSON(http.MethodPost, "/axelix-caches/"+manager+"/"+cache+"/enable", nil, nil)
}

func (c *Client) DisableCache(manager, cache string) error {
	if cache == "" {
		return c.doJSON(http.MethodPost, "/axelix-caches/"+manager+"/disable", nil, nil)
	}
	return c.doJSON(http.MethodPost, "/axelix-caches/"+manager+"/"+cache+"/disable", nil, nil)
}

// ClearAllCaches clears every cache across all managers.
func (c *Client) ClearAllCaches() error {
	return c.doDelete("/axelix-caches")
}

// ClearManagerCaches clears all caches in the given manager.
func (c *Client) ClearManagerCaches(manager string) error {
	return c.doDelete("/axelix-caches/" + manager + "/clear-all")
}

// ClearCache clears a specific cache, optionally by key.
func (c *Client) ClearCache(manager, cache, key string) error {
	path := "/axelix-caches/" + manager + "/" + cache + "/clear"
	if key != "" {
		path += "?key=" + key
	}
	return c.doDelete(path)
}

func (c *Client) GetConditions() (*models.ConditionsFeed, error) {
	out := &models.ConditionsFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-conditions", nil, out)
}

func (c *Client) GetConfigProps() (*models.ConfigPropsFeed, error) {
	out := &models.ConfigPropsFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-configprops", nil, out)
}

func (c *Client) GetDetails() (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-details", nil, &out)
}

func (c *Client) GetEnv(pattern string) (*models.EnvironmentFeed, error) {
	path := "/axelix-env"
	if pattern != "" {
		path += "?pattern=" + pattern
	}
	out := &models.EnvironmentFeed{}
	return out, c.doJSON(http.MethodGet, path, nil, out)
}

func (c *Client) GetGCLogStatus() (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-gc/log/status", nil, &out)
}

func (c *Client) GetGCLogFile() (string, error) {
	return c.doText("/axelix-gc/log/file")
}

func (c *Client) TriggerGC() error {
	return c.doJSON(http.MethodPost, "/axelix-gc/trigger", nil, nil)
}

func (c *Client) EnableGCLog(level string) error {
	return c.doJSON(http.MethodPost, "/axelix-gc/log/enable", map[string]string{"level": level}, nil)
}

func (c *Client) DisableGCLog() error {
	return c.doJSON(http.MethodPost, "/axelix-gc/log/disable", nil, nil)
}

// DownloadHeapDump downloads a heap dump. Pass live=true to include only live objects.
func (c *Client) DownloadHeapDump(live bool) ([]byte, error) {
	path := "/axelix-heap-dump"
	if live {
		path += "?live=true"
	}
	return c.doRaw(path)
}

func (c *Client) GetLoggers() (*models.ServiceLoggers, error) {
	out := &models.ServiceLoggers{}
	return out, c.doJSON(http.MethodGet, "/axelix-loggers", nil, out)
}

func (c *Client) GetLogger(name string) (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-loggers/"+name, nil, &out)
}

func (c *Client) SetLogLevel(name, level string) error {
	return c.doJSON(http.MethodPost, "/axelix-loggers/"+name, map[string]string{"configuredLevel": level}, nil)
}

func (c *Client) GetMetadata() (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-metadata", nil, &out)
}

func (c *Client) GetMetrics() (map[string]any, error) {
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, "/axelix-metrics", nil, &out)
}

func (c *Client) GetMetric(name, tag string) (map[string]any, error) {
	path := "/axelix-metrics/" + name
	if tag != "" {
		path += "?tag=" + tag
	}
	out := map[string]any{}
	return out, c.doJSON(http.MethodGet, path, nil, &out)
}

func (c *Client) GetScheduledTasks() (*models.ServiceScheduledTasks, error) {
	out := &models.ServiceScheduledTasks{}
	return out, c.doJSON(http.MethodGet, "/axelix-scheduled-tasks", nil, out)
}

func (c *Client) EnableScheduledTask(trigger string, force bool) error {
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/enable", map[string]any{"trigger": trigger, "force": force}, nil)
}

func (c *Client) DisableScheduledTask(trigger string, force bool) error {
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/disable", map[string]any{"trigger": trigger, "force": force}, nil)
}

func (c *Client) ExecuteScheduledTask(trigger string) error {
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/execute", map[string]string{"trigger": trigger}, nil)
}

func (c *Client) SetTaskCronExpression(trigger, expr string) error {
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/modify/cron-expression", map[string]string{"trigger": trigger, "cronExpression": expr}, nil)
}

func (c *Client) SetTaskInterval(trigger string, intervalMs int64) error {
	return c.doJSON(http.MethodPost, "/axelix-scheduled-tasks/modify/interval", map[string]any{"trigger": trigger, "interval": intervalMs}, nil)
}

func (c *Client) GetThreadDump() (*models.ThreadDumpFeed, error) {
	out := &models.ThreadDumpFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-thread-dump", nil, out)
}

func (c *Client) EnableThreadContention() error {
	return c.doJSON(http.MethodPost, "/axelix-thread-dump/enable", nil, nil)
}

func (c *Client) DisableThreadContention() error {
	return c.doJSON(http.MethodPost, "/axelix-thread-dump/disable", nil, nil)
}

func (c *Client) GetTransactions() (*models.TransactionMonitoringFeed, error) {
	out := &models.TransactionMonitoringFeed{}
	return out, c.doJSON(http.MethodGet, "/axelix-transactions-monitoring", nil, out)
}

func (c *Client) ClearTransactions() error {
	return c.doDelete("/axelix-transactions-monitoring")
}
