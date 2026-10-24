package runner

import (
	"context"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/ssqwet1324/load-test/internal/config"
)

type Result struct {
	Duration   time.Duration
	StatusCode int
	Err        error
}

// doRequest делает один GET-запрос и меряет время
func doRequest(ctx context.Context, client *http.Client, url string) Result {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{Err: err}
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return Result{Duration: time.Since(start), Err: err}
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		return Result{Duration: time.Since(start), Err: err}
	}
	elapsed := time.Since(start)

	return Result{Duration: elapsed, StatusCode: resp.StatusCode}
}

func Run(ctx context.Context, cfg config.Config) []Result {
	client := &http.Client{Timeout: cfg.Timeout}
	jobs := make(chan struct{})
	results := make(chan Result)

	var wg sync.WaitGroup

	for range cfg.Concurrency {
		wg.Go(func() {
			for range jobs {
				results <- doRequest(ctx, client, cfg.URL)
			}
		})
	}

	go func() {
		defer close(jobs)
		for range cfg.Requests {
			select {
			case jobs <- struct{}{}:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	collected := make([]Result, 0, cfg.Requests)
	for r := range results {
		collected = append(collected, r)
	}

	return collected
}
