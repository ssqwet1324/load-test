package app

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/ssqwet1324/load-test/internal/config"
	"github.com/ssqwet1324/load-test/internal/runner"
)

func Run() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", cfg)

	client := &http.Client{Timeout: cfg.Timeout}
	fmt.Printf("%+v\n", runner.DoRequest(context.Background(), client, cfg.URL))
}
