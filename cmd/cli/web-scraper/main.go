package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/landmaster135/devbox/internal/web_scraper/config"
	"github.com/landmaster135/devbox/internal/web_scraper/interfaces/fetchers"
	"github.com/landmaster135/devbox/internal/web_scraper/usecases"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}

func run(args []string, stdout, stderr io.Writer) exitCode {
	fs := flag.NewFlagSet("web-scraper", flag.ContinueOnError)
	fs.SetOutput(stderr)

	cfg := &config.Config{}
	fs.StringVar(&cfg.Operation, "operation", "", "実行するoperation (例: get_dom_tree)")
	fs.StringVar(&cfg.URL, "url", "", "対象のURL")
	fs.IntVar(&cfg.WaitSeconds, "wait-seconds", 0, "ページ読み込み後にDOMを取得するまでの待機秒数（秒）")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCodeError
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(stderr, "入力値が不正です: %v\n", err)
		fs.Usage()
		return exitCodeError
	}

	denySelectorsForReddit := []string{
		"faceplate-loader",
		"faceplate-tracker",
		"faceplate-perfmark",
		"shreddit-comments-page-ad",
		"shreddit-async-loader",
	}
	denySelectors := []string{}
	denySelectors = append(denySelectors, denySelectorsForReddit...)

	ctx := context.Background()
	domFetcher := fetchers.NewRodDOMFetcher()
	service := usecases.NewDOMService(domFetcher)

	switch cfg.OperationName() {
	case config.OperationGetDOMTree:
		html, err := service.GetDOMTree(ctx, cfg.URL, cfg.WaitDuration(), denySelectors)
		if err != nil {
			fmt.Fprintf(stderr, "DOM取得に失敗しました: %v\n", err)
			return exitCodeError
		}
		fmt.Fprintln(stdout, html)
		return exitCodeOK
	default:
		fmt.Fprintf(stderr, "未対応のoperationです: %s\n", cfg.Operation)
		return exitCodeError
	}
}
