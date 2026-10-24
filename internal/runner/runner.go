package runner

import (
	"context"
	"io"
	"net/http"
	"time"
)

type Result struct {
	Duration   time.Duration
	StatusCode int
	Err        error
}

// DoRequest делает один GET-запрос и меряет время
func DoRequest(ctx context.Context, client *http.Client, url string) Result {
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
