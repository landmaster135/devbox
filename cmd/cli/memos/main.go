package main

import (
	"os"

	runner "github.com/landmaster135/devbox/cmd/cli/memos/runner"
	cfg "github.com/landmaster135/devbox/internal/memos/config"
	usecases "github.com/landmaster135/devbox/internal/memos/usecases"
)

type serviceFactory func(conf *cfg.Config) usecases.MemoService

func main() {
	os.Exit(runner.Run(os.Args[1:], os.Stdout, os.Stderr, runner.NewServiceFromConfig))
}
