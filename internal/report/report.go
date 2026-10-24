package report

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"math"
	"slices"
	"time"

	"github.com/ssqwet1324/load-test/internal/runner"
)

type Report struct {
	Total       int
	Success     int
	Errors      int
	TotalTime   time.Duration
	RPS         float64
	Min         time.Duration
	Max         time.Duration
	Avg         time.Duration
	P50         time.Duration
	P95         time.Duration
	P99         time.Duration
	StatusCodes map[int]int
}

// Build считает статистику по результатам
func Build(results []runner.Result, totalTime time.Duration) Report {
	report := Report{
		TotalTime:   totalTime,
		StatusCodes: make(map[int]int),
	}

	latencies := make([]time.Duration, 0, len(results))

	for _, result := range results {
		if errors.Is(result.Err, context.Canceled) {
			continue
		}

		if result.Err != nil {
			report.Errors++
			continue
		}

		report.Success++
		report.StatusCodes[result.StatusCode]++
		latencies = append(latencies, result.Duration)
	}

	report.Total = report.Success + report.Errors

	if totalTime > 0 {
		report.RPS = float64(report.Total) / totalTime.Seconds()
	}

	if len(latencies) == 0 {
		return report
	}

	slices.Sort(latencies)

	report.Min = latencies[0]
	report.Max = latencies[len(latencies)-1]

	var sum time.Duration
	for _, l := range latencies {
		sum += l
	}
	report.Avg = sum / time.Duration(len(latencies))

	report.P50 = percentile(latencies, 50)
	report.P95 = percentile(latencies, 95)
	report.P99 = percentile(latencies, 99)

	return report
}

// percentile возвращает p-й перцентиль
func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}

	idx := int(math.Ceil(p/100*float64(len(sorted)))) - 1
	idx = max(0, min(idx, len(sorted)-1))

	return sorted[idx]
}

// Print выводит отчёт в консоль
func (r Report) Print() {
	fmt.Println()
	fmt.Printf("Total:    %d (success %d, errors %d)\n", r.Total, r.Success, r.Errors)
	fmt.Printf("Time:     %s\n", round(r.TotalTime))
	fmt.Printf("RPS:      %.1f\n", r.RPS)

	if r.Success == 0 {
		fmt.Println("Latency:  no successful requests")
		return
	}

	fmt.Printf("Latency:  min %-10s avg %-10s max %s\n", round(r.Min), round(r.Avg), round(r.Max))
	fmt.Printf("          p50 %-10s p95 %-10s p99 %s\n", round(r.P50), round(r.P95), round(r.P99))

	fmt.Println("Status:")
	for _, code := range slices.Sorted(maps.Keys(r.StatusCodes)) {
		fmt.Printf("          %d × %d\n", code, r.StatusCodes[code])
	}
}

func round(d time.Duration) time.Duration {
	switch {
	case d >= time.Second:
		return d.Round(10 * time.Millisecond)
	case d >= time.Millisecond:
		return d.Round(100 * time.Microsecond)
	default:
		return d.Round(time.Microsecond)
	}
}
