// main.go
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	usecases "github.com/landmaster135/devbox/internal/independencies/image_trimmer/usecases"
)

func main() {
	// 既存オプション
	x1 := flag.Int("x1", 0, "X-coordinate on upper left")
	y1 := flag.Int("y1", 0, "Y-coordinate on upper left")
	x2 := flag.Int("x2", 0, "X-coordinate on lower right")
	y2 := flag.Int("y2", 0, "Y-coordinate on lower right")
	suffix := flag.String("suffix", "trimmed", "suffix to attach to file name to save")
	move := flag.Bool("move", false, "move originals instead of copying (effective only with -archive)")
	srcDir := flag.String("src", ".", "source directory to scan")
	outDir := flag.String("out", "", "output directory (output into src if empty)")
	arcDir := flag.String("arc", "./5_original_files", "move processed originals to this directory")
	recursive := flag.Bool("r", false, "recursively scan sub-directories")

	flag.Parse()

	if *outDir == "" {
		*outDir = *srcDir
	}

	// --- ディレクトリを列挙 ---
	err := filepath.WalkDir(*srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 再帰しない場合はサブディレクトリへ入らない
		if d.IsDir() && path != *srcDir && !*recursive {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}

		switch ext := strings.ToLower(filepath.Ext(d.Name())); ext {
		case ".png", ".jpg", ".jpeg":
			if err := usecases.CropAndSave(path, *outDir, *x1, *y1, *x2, *y2, *suffix); err != nil {
				return err
			}
			if *move {
				if err := usecases.MoveOriginal(path, *arcDir); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
