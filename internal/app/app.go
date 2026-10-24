package app

import (
	"fmt"
	"os"

	"github.com/ssqwet1324/load-test/internal/config"
)

func Run() {
	cfg, err := config.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", cfg)
}
