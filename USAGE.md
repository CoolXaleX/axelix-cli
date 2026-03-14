# axelix CLI — Usage Examples

## Configuration

Services are stored in `~/.axelix/config.json`. The service name is always the first argument to every command.

```bash
# Add named services
axelix config add local   http://localhost:8080
axelix config add staging http://staging.internal:8080
axelix config add prod    http://prod.internal:8080

# List all configured services
axelix config list

# Remove a service
axelix config remove staging
```

---

## Command Structure

```
axelix <service> <command> [subcommand] [flags]
axelix config <subcommand>
```

**Global flag:** `--json` — output raw JSON instead of a table (works with any command)

```bash
axelix local  beans
axelix prod   beans --json | jq '.beans[] | select(.isPrimary == true)'
axelix --json prod metrics get --name jvm.memory.used
```

---

## Beans

```bash
# List all beans: Name | Scope | Class | ProxyType | Primary | Lazy
axelix local beans
axelix prod  beans --json | jq '[.beans[] | select(.scope == "singleton") | .beanName]'
```

---

## Caches

```bash
# List all cache managers and caches
axelix local caches list

# Get a specific cache
axelix local caches get --manager cacheManager --cache myCache

# Enable/disable an entire manager
axelix local caches enable  --manager cacheManager
axelix local caches disable --manager cacheManager

# Enable/disable a specific cache
axelix local caches enable  --manager cacheManager --cache myCache
axelix local caches disable --manager cacheManager --cache myCache

# Clear all caches across all managers (requires explicit --all)
axelix local caches clear --all

# Clear all caches in a manager
axelix local caches clear --manager cacheManager

# Clear a specific cache
axelix local caches clear --manager cacheManager --cache myCache

# Evict a single key from a cache
axelix local caches clear --manager cacheManager --cache myCache --key "user:42"
```

---

## Conditions

```bash
# Positive matches (what auto-configured and why)
axelix local conditions

# Negative matches (what did NOT configure and why)
axelix local conditions --negative
```

---

## ConfigProps

```bash
# All @ConfigurationProperties beans with prefixes and values
axelix local configprops
axelix prod  configprops --json | jq '.beans[] | select(.prefix | startswith("spring.datasource"))'
```

---

## Details

```bash
# Git, Runtime, Spring, Build, OS info
axelix local details
axelix prod  details --json
```

---

## Env

```bash
# All property sources with their properties
axelix local env

# Filter by pattern
axelix local env --pattern "spring.datasource.*"
axelix prod  env --pattern "server.port"
```

---

## GC

```bash
# GC logging status
axelix local gc status

# Trigger GC manually
axelix local gc trigger

# Enable GC logging
axelix local gc log-enable --level INFO
axelix local gc log-enable --level DEBUG

# Disable GC logging
axelix local gc log-disable

# Print GC log file content
axelix local gc log-file
axelix local gc log-file | grep "GC pause"
```

---

## Heap Dump

```bash
# Download heap dump (all objects, live and dead)
axelix local heap-dump

# Download heap dump with only live objects (smaller file, no GC-able objects)
axelix local heap-dump --live

# Specify output path
axelix local heap-dump --out /tmp/myapp.hprof
axelix local heap-dump --live --out ./dumps/heap-$(date +%s).hprof
```

> File is saved locally; path is printed to stderr. `--json` has no effect on this command.

---

## Loggers

```bash
# List all loggers: Logger | Configured Level | Effective Level
axelix local loggers list
axelix prod  loggers list --json | jq '.loggers | to_entries[] | select(.value.effectiveLevel == "DEBUG")'

# Get a specific logger
axelix local loggers get --name com.example.MyService

# Change log level
axelix local loggers set --name com.example.MyService --level DEBUG
axelix local loggers set --name ROOT --level WARN
axelix prod  loggers set --name com.example --level TRACE

# Reset to default level (pass empty string)
axelix local loggers set --name com.example.MyService --level ""
```

---

## Metadata

```bash
# Version, git commit, build time, artifact info
axelix local metadata
axelix prod  metadata --json
```

---

## Metrics

```bash
# List all metric groups
axelix local metrics list

# Get a specific metric
axelix local metrics get --name jvm.memory.used
axelix local metrics get --name jvm.memory.used --tag area:heap
axelix prod  metrics get --name http.server.requests --tag uri:/api/v1/users
```

---

## Scheduled Tasks

```bash
# List all tasks: Type | Target | Schedule | Enabled | Next Exec | Last Status
axelix local scheduled-tasks list

# Enable a task (trigger = fully qualified class + method name)
axelix local scheduled-tasks enable --trigger com.example.jobs.ReportJob.generate

# Disable a task
axelix local scheduled-tasks disable --trigger com.example.jobs.ReportJob.generate

# Force disable (even if the task is currently running)
axelix local scheduled-tasks disable --trigger com.example.jobs.ReportJob.generate --force

# Execute a task immediately
axelix local scheduled-tasks execute --trigger com.example.jobs.ReportJob.generate

# Change cron expression
axelix local scheduled-tasks set-cron --trigger com.example.jobs.ReportJob.generate --cron "0 0 * * * *"

# Change interval (in milliseconds)
axelix local scheduled-tasks set-interval --trigger com.example.jobs.CleanupJob.run --interval 60000
```

---

## Thread Dump

```bash
# List threads: Name | ID | State | Daemon | Priority | Blocked Count
axelix local thread-dump get
axelix prod  thread-dump get --json | jq '.threads[] | select(.threadState == "BLOCKED")'

# Enable contention monitoring (tracks lock wait/block time per thread)
axelix local thread-dump enable-contention

# Disable contention monitoring
axelix local thread-dump disable-contention
```

---

## Transactions

```bash
# List @Transactional methods: Class | Method | Executions | Avg ms | Max ms | Median ms
axelix local transactions list

# Clear accumulated statistics
axelix local transactions clear
```

---

## Common Scenarios

```bash
# Memory diagnosis
axelix local metrics get --name jvm.memory.used --tag area:heap
axelix local gc trigger
axelix local heap-dump --live --out ./heap.hprof

# Enable DEBUG logging for a package, then restore
axelix local loggers set --name com.example.api --level DEBUG
# ... check logs ...
axelix local loggers set --name com.example.api --level INFO

# Compare environments
axelix local env --pattern "spring.datasource.url"
axelix prod  env --pattern "spring.datasource.url"

# Find blocked threads on prod
axelix prod thread-dump get --json | jq '[.threads[] | select(.threadState == "BLOCKED") | {name:.threadName, state:.threadState}]'

# Check what auto-configured
axelix local conditions | grep DataSource

# Run a job on staging, check transactions after
axelix staging scheduled-tasks execute --trigger com.example.jobs.ReportJob.generate
axelix staging transactions list
```
