# load-test

A simple HTTP load testing tool written in Go, similar to `hey` or `ab`.

Sends a configurable number of concurrent requests to a target URL and reports throughput, latency percentiles and status code distribution.

## Features

- Concurrent requests via a worker pool
- Per-request timeout
- Latency stats: min, avg, max, p50, p95, p99
- Requests per second (RPS)
- Status code breakdown
- Graceful shutdown on Ctrl+C – prints a report for completed requests

## Installation

```bash
git clone https://github.com/ssqwet1324/load-test.git
cd load-test
go build -o loadtest ./cmd
```

Requires Go 1.25+.

## Usage

```bash
./cmd -url http://localhost:8080 -n 10000 -c 100 -timeout 2s
```

| Flag       | Description                   | Default  |
|------------|-------------------------------|----------|
| `-url`     | Target URL (http or https)    | required |
| `-n`       | Total number of requests      | `100`    |
| `-c`       | Number of concurrent workers  | `10`     |
| `-timeout` | Timeout per request           | `5s`     |

## Example output

```
Total:    10000 (success 10000, errors 0)
Time:     2.58s
RPS:      3868.6
Latency:  min 0s         avg 25.6ms     max 90.8ms
          p50 25.3ms     p95 48.1ms     p99 50.1ms
Status:
          200 × 9503
          500 × 497
```

## How it works

```
flags → config → runner (worker pool) → []Result → report → output
```

- `internal/config` – parses and validates CLI flags
- `internal/runner` – worker pool: a dispatcher sends jobs to a channel, workers execute requests and send results back
- `internal/report` – calculates statistics and prints the report
- `internal/app` – wires everything together