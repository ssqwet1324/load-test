package app

import (
	"context"
	"time"

	"github.com/ssqwet1324/load-test/internal/config"
	"github.com/ssqwet1324/load-test/internal/report"
	"github.com/ssqwet1324/load-test/internal/runner"
)

func Run(ctx context.Context) error {
	cfg, err := config.Parse()
	if err != nil {
		return err
	}

	start := time.Now()
	results := runner.Run(ctx, cfg)
	rep := report.Build(results, time.Since(start))
	rep.Print()

	return nil
}
