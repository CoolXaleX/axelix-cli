# axelix CLI — Usage Examples

## Configuration

The CLI stores named services in `~/.axelix/config.json`. One service is marked as *current* and used by default.

```bash
# Add a service (first one added becomes current automatically)
axelix config add local http://localhost:8080
axelix config add staging http://staging.internal:8080
axelix config add prod http://prod.internal:8080

# Switch the active service
axelix config use prod

# List all services (✓ marks the active one)
axelix config list

# Remove a service
axelix config remove staging
```

---

## Global Flags (available on all commands)

```bash
--url http://localhost:8080   # One-off URL override (bypasses saved services)
--service staging             # Use a specific named service for this call
--json                        # Output raw JSON (useful for jq)
```

**Priority:** `--url` → `AXELIX_URL` env var → `--service` → current service in config

---

## Beans

```bash
# List all beans: Name | Scope | Class | ProxyType | Primary | Lazy
axelix beans
axelix beans --json | jq '.beans[] | select(.isPrimary == true)'
```

---

## Caches

```bash
# List all cache managers and caches
axelix caches list

# Get a specific cache
axelix caches get --manager cacheManager --cache myCache

# Enable/disable an entire manager
axelix caches enable  --manager cacheManager
axelix caches disable --manager cacheManager

# Enable/disable a specific cache
axelix caches enable  --manager cacheManager --cache myCache
axelix caches disable --manager cacheManager --cache myCache

# Clear all caches across all managers (requires explicit --all)
axelix caches clear --all

# Clear all caches in a manager
axelix caches clear --manager cacheManager

# Clear a specific cache
axelix caches clear --manager cacheManager --cache myCache

# Evict a single key from a cache
axelix caches clear --manager cacheManager --cache myCache --key "user:42"
```

---

## Conditions

```bash
# Positive matches (what auto-configured and why)
axelix conditions

# Negative matches (what did NOT configure and why)
axelix conditions --negative
```

---

## ConfigProps

```bash
# All @ConfigurationProperties beans with prefixes and values
axelix configprops
axelix configprops --json | jq '.beans[] | select(.prefix | startswith("spring.datasource"))'
```

---

## Details

```bash
# Git, Runtime, Spring, Build, OS info
axelix details
axelix details --json
```

---

## Env

```bash
# All property sources with their properties
axelix env

# Filter by pattern
axelix env --pattern "spring.datasource.*"
axelix env --pattern "server.port"
```

---

## GC

```bash
# GC logging status
axelix gc status

# Trigger GC manually
axelix gc trigger

# Enable GC logging
axelix gc log-enable --level INFO
axelix gc log-enable --level DEBUG

# Disable GC logging
axelix gc log-disable

# Print GC log file content
axelix gc log-file
axelix gc log-file | grep "GC pause"
```

---

## Heap Dump

```bash
# Download heap dump (all objects, live and dead)
axelix heap-dump

# Download heap dump with only live objects (smaller file, no GC-able objects)
axelix heap-dump --live

# Specify output path
axelix heap-dump --out /tmp/myapp.hprof
axelix heap-dump --live --out ./dumps/heap-$(date +%s).hprof
```

> File is saved locally; path is printed to stderr. `--json` has no effect on this command.

---

## Loggers

```bash
# List all loggers: Logger | Configured Level | Effective Level
axelix loggers list
axelix loggers list --json | jq '.[] | select(.effectiveLevel == "DEBUG")'

# Get a specific logger
axelix loggers get --name com.example.MyService

# Change log level
axelix loggers set --name com.example.MyService --level DEBUG
axelix loggers set --name ROOT --level WARN
axelix loggers set --name com.example --level TRACE

# Reset to default level (pass empty string)
axelix loggers set --name com.example.MyService --level ""
```

---

## Metadata

```bash
# Version, git commit, build time, artifact info
axelix metadata
axelix metadata --json
```

---

## Metrics

```bash
# List all metric groups
axelix metrics list

# Get a specific metric
axelix metrics get --name jvm.memory.used
axelix metrics get --name jvm.memory.used --tag area:heap
axelix metrics get --name http.server.requests --tag uri:/api/v1/users
```

---

## Scheduled Tasks

```bash
# List all tasks: Type | Target | Schedule | Enabled | Next Exec | Last Status
axelix scheduled-tasks list

# Enable a task (trigger = fully qualified class + method name)
axelix scheduled-tasks enable --trigger com.example.jobs.ReportJob.generate

# Disable a task
axelix scheduled-tasks disable --trigger com.example.jobs.ReportJob.generate

# Force disable (even if the task is currently running)
axelix scheduled-tasks disable --trigger com.example.jobs.ReportJob.generate --force

# Execute a task immediately
axelix scheduled-tasks execute --trigger com.example.jobs.ReportJob.generate

# Change cron expression
axelix scheduled-tasks set-cron --trigger com.example.jobs.ReportJob.generate --cron "0 0 * * * *"

# Change interval (in milliseconds)
axelix scheduled-tasks set-interval --trigger com.example.jobs.CleanupJob.run --interval 60000
```

---

## Thread Dump

```bash
# List threads: Name | ID | State | Daemon | Priority | Blocked Count
axelix thread-dump get
axelix thread-dump get --json | jq '.threads[] | select(.threadState == "BLOCKED")'

# Enable contention monitoring (tracks lock wait/block time per thread)
axelix thread-dump enable-contention

# Disable contention monitoring
axelix thread-dump disable-contention
```

---

## Transactions

```bash
# List @Transactional methods: Class | Method | Executions | Avg ms | Max ms | Median ms
axelix transactions list

# Clear accumulated statistics
axelix transactions clear
```

---

## Common Scenarios

```bash
# Work with multiple environments
axelix config add local   http://localhost:8080
axelix config add staging http://staging:8080
axelix config add prod    http://prod:8080
axelix config use staging

axelix beans                          # hits staging
axelix beans --service prod           # hits prod for this one call
axelix beans --url http://other:8080  # hits arbitrary URL

# Memory diagnosis
axelix metrics get --name jvm.memory.used --tag area:heap
axelix gc trigger
axelix heap-dump --live --out ./heap.hprof

# Enable DEBUG logging for a package, then restore
axelix loggers set --name com.example.api --level DEBUG
# ... check logs ...
axelix loggers set --name com.example.api --level INFO

# Find blocked threads
axelix thread-dump get --json | jq '[.threads[] | select(.threadState == "BLOCKED") | {name:.threadName, state:.threadState}]'

# Check what auto-configured
axelix conditions | grep DataSource

# Pipe into jq
axelix beans --json | jq '[.beans[] | select(.scope == "singleton") | .beanName]'
```
