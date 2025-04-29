package main

import (
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	usecases "github.com/landmaster135/devbox/internal/independencies/image_converter/usecases"
)

func main() {
	srcDir := flag.String("src", ".", "source directory to scan")
	outDir := flag.String("out", "./999_converted_images", "output directory")
	archiveDir := flag.String("archive", "", "move processed originals to this directory (disabled if empty)")
	moveOrig := flag.Bool("move", false, "move originals instead of copying (effective only with -archive)")
	outExt := flag.String("ext", "png", "target extension (png|jpg|webp|avif)")
	quality := flag.Int("q", 80, "quality for lossy formats (1-100)")
	workers := flag.Int("workers", runtime.NumCPU(), "number of concurrent workers")
	recursive := flag.Bool("R", false, "recursively scan sub-directories")

	flag.Parse()
	*outExt = strings.ToLower(strings.TrimPrefix(*outExt, "."))
	codecs := usecases.MakeCodecTable(*quality)

	if _, ok := codecs[*outExt]; !ok {
		log.Fatalf("unsupported target format: %s", *outExt)
	}

	// 1. collect paths
	paths := make(chan string, 512)
	go func() {
		if *recursive {
			filepath.WalkDir(*srcDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(path)), ".")
				if _, ok := codecs[ext]; ok {
					paths <- path
				}
				return nil
			})
		} else {
			entries, err := os.ReadDir(*srcDir)
			if err != nil {
				log.Fatalf("read dir: %v", err)
			}
			for _, e := range entries {
				if e.IsDir() { // サブフォルダは完全にスキップ
					continue
				}
				ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(e.Name())), ".")
				if _, ok := codecs[ext]; ok {
					paths <- filepath.Join(*srcDir, e.Name())
				}
			}
		}
		close(paths)
	}()

	// 2. worker pool
	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range paths {
				if err := usecases.ConvertFile(p, *srcDir, *outDir, *outExt, codecs); err != nil {
					log.Printf("warn: %v", err)
				}
				if err := usecases.ArchiveOriginal(p, *srcDir, *archiveDir, *moveOrig); err != nil {
					log.Printf("warn (archive): %v", err)
				}
			}
		}()
	}
	wg.Wait()
}
