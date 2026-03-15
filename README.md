# axelix CLI

A command-line tool for managing Spring Boot applications via the [Axelix Spring Boot Suite](https://github.com/axelixlabs) actuator extension.

## Installation

```bash
go install github.com/axelixlabs/axelix-cli@latest
```

Or build from source:

```bash
git clone https://github.com/axelixlabs/axelix-cli
cd axelix-cli
go build -o axelix .
```

## Configuration

Services are stored in `~/.axelix/config.json`. The service name is always the first argument to every command.

```bash
axelix config add local   http://localhost:8080
axelix config add staging http://staging.internal:8080
axelix config add prod    http://prod.internal:8080

axelix config list
axelix config remove staging
```

## Command Structure

```
axelix <service> <command> [subcommand] [flags]
axelix config <subcommand>
```

**Global flag:** `--json` — output raw JSON instead of a table (works with any command).

```bash
axelix local beans
axelix prod  beans --json | jq '.beans[] | select(.isPrimary == true)'
```

## Commands

### Beans

```bash
axelix local beans
axelix prod  beans --json | jq '[.beans[] | select(.scope == "singleton") | .beanName]'
```

### Caches

```bash
axelix local caches list
axelix local caches get     --manager cacheManager --cache myCache
axelix local caches enable  --manager cacheManager [--cache myCache]
axelix local caches disable --manager cacheManager [--cache myCache]
axelix local caches clear  --all
axelix local caches clear  --manager cacheManager
axelix local caches clear  --manager cacheManager --cache myCache
axelix local caches clear  --manager cacheManager --cache myCache --key "user:42"
```

### Conditions

```bash
axelix local conditions            # positive matches (what auto-configured)
axelix local conditions --negative # negative matches (what did NOT configure)
```

### ConfigProps

```bash
axelix local configprops
axelix prod  configprops --json | jq '.beans[] | select(.prefix | startswith("spring.datasource"))'
```

### Details

```bash
axelix local details       # git, runtime, Spring, build, OS info
axelix prod  details --json
```

### Env

```bash
axelix local env
axelix local env --pattern "spring.datasource.*"
axelix prod  env --pattern "server.port"
```

### GC

```bash
axelix local gc status
axelix local gc trigger
axelix local gc log-enable  --level INFO
axelix local gc log-disable
axelix local gc log-file
axelix local gc log-file | grep "GC pause"
```

### Heap Dump

```bash
axelix local heap-dump                                              # all objects
axelix local heap-dump --live                                       # live objects only
axelix local heap-dump --live --out ./dumps/heap-$(date +%s).hprof
```

> The file is saved locally; its path is printed to stderr. `--json` has no effect on this command.

### Loggers

```bash
axelix local loggers list
axelix local loggers get --name com.example.MyService
axelix local loggers set --name com.example.MyService --level DEBUG
axelix local loggers set --name com.example.MyService --level ""   # reset to default
```

### Metadata

```bash
axelix local metadata
axelix prod  metadata --json
```

### Metrics

```bash
axelix local metrics list
axelix local metrics get --name jvm.memory.used
axelix local metrics get --name jvm.memory.used --tag area:heap
axelix prod  metrics get --name http.server.requests --tag uri:/api/v1/users
```

### Scheduled Tasks

```bash
axelix local scheduled-tasks list
axelix local scheduled-tasks enable      --trigger com.example.jobs.ReportJob.generate
axelix local scheduled-tasks disable     --trigger com.example.jobs.ReportJob.generate [--force]
axelix local scheduled-tasks execute     --trigger com.example.jobs.ReportJob.generate
axelix local scheduled-tasks set-cron    --trigger com.example.jobs.ReportJob.generate --cron "0 0 * * * *"
axelix local scheduled-tasks set-interval --trigger com.example.jobs.CleanupJob.run --interval 60000
```

### Thread Dump

```bash
axelix local thread-dump get
axelix prod  thread-dump get --json | jq '.threads[] | select(.threadState == "BLOCKED")'
axelix local thread-dump enable-contention
axelix local thread-dump disable-contention
```

### Transactions

```bash
axelix local transactions list
axelix local transactions clear
```

## Common Scenarios

```bash
# Memory diagnosis
axelix local metrics get --name jvm.memory.used --tag area:heap
axelix local gc trigger
axelix local heap-dump --live --out ./heap.hprof

# Enable DEBUG logging temporarily
axelix local loggers set --name com.example.api --level DEBUG
axelix local loggers set --name com.example.api --level INFO

# Compare environments
axelix local env --pattern "spring.datasource.url"
axelix prod  env --pattern "spring.datasource.url"

# Find blocked threads on prod
axelix prod thread-dump get --json | jq '[.threads[] | select(.threadState == "BLOCKED") | {name:.threadName, state:.threadState}]'

# Run a job on staging, check transactions after
axelix staging scheduled-tasks execute --trigger com.example.jobs.ReportJob.generate
axelix staging transactions list
```

## Requirements

- Go 1.22+
- Spring Boot application with Axelix SBS actuator extension enabled
