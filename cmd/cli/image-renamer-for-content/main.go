package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	cfg "github.com/landmaster135/devbox/internal/image_renamer_for_content/config"
	usecases "github.com/landmaster135/devbox/internal/image_renamer_for_content/usecases"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeError
)

func parseFlags(args []string, stderr io.Writer) (cfg.Config, error) {
	var (
		srcDir     string
		sortByName bool
		sortByTime bool
		contentID  string
		suffix     string
		delimiter  string
		digits     int
		start      int
		recursive  bool
		workers    int
	)

	defaultWorkers := cfg.DefaultWorkers()

	flagSet := flag.NewFlagSet("image-renamer-for-content", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	flagSet.StringVar(&srcDir, "src", ".", "リネーム対象のソースディレクトリ")
	flagSet.BoolVar(&sortByName, "name", false, "ファイル名順で並べ替え")
	flagSet.BoolVar(&sortByTime, "time", false, "更新日時順で並べ替え")
	flagSet.StringVar(&contentID, "content", "", "コンテンツID (必須)")
	flagSet.StringVar(&suffix, "suffix", "01", "シリアルの後ろに付加するサフィックス")
	flagSet.StringVar(&delimiter, "delimiter", "", "コンテンツIDとシリアルの間に挟む区切り文字")
	flagSet.IntVar(&digits, "digits", 4, "シリアル番号の桁数")
	flagSet.IntVar(&start, "start", 1, "シリアル番号の開始値")
	flagSet.BoolVar(&recursive, "r", false, "サブディレクトリを再帰的に処理")
	flagSet.IntVar(&workers, "workers", defaultWorkers, fmt.Sprintf("並行ワーカー数 (デフォルト: %d)", defaultWorkers))

	if err := flagSet.Parse(args); err != nil {
		return cfg.Config{}, err
	}

	if workers <= 0 {
		workers = cfg.DefaultWorkers()
	}

	return cfg.Config{
		SrcDir:     srcDir,
		SortByName: sortByName,
		SortByTime: sortByTime,
		ContentID:  contentID,
		Suffix:     suffix,
		Delimiter:  delimiter,
		Digits:     digits,
		Start:      start,
		Recursive:  recursive,
		Workers:    workers,
	}, nil
}

func run(args []string, stdout, stderr io.Writer) exitCode {
	config, err := parseFlags(args, stderr)
	if err != nil {
		return exitCodeError
	}

	_, errorCount, err := usecases.ProcessContentImageRename(config, stdout, stderr)
	if err != nil {
		return exitCodeError
	}

	if errorCount > 0 {
		return exitCodeError
	}

	return exitCodeOK
}

func main() {
	code := run(os.Args[1:], os.Stdout, os.Stderr)
	os.Exit(int(code))
}
