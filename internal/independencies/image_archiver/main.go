// main.go  (Go 1.22+, go-jpeg-image-structure/v2, go-exif/v3)
package main

// import (
// 	"bufio"
// 	"bytes"
// 	"fmt"
// 	"io"
// 	"io/fs"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	exif "github.com/dsoprea/go-exif/v3"
// 	exifcommon "github.com/dsoprea/go-exif/v3/common"
// 	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
// )

// const (
// 	dirMovieEscaping     = "1-1_movie_escaping"
// 	dirCreateDateSetting = "1-2_create_date_setting"
// )

// var (
// 	logFile *os.File
// 	logger  *log.Logger
// )

// /* ---------- ログ ---------- */

// func initLog() {
// 	name := fmt.Sprintf("log_%s.txt", time.Now().Format("20060102"))
// 	var err error
// 	logFile, err = os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
// 	if err != nil {
// 		panic(err)
// 	}
// 	logger = log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)
// }

// func closeLog() { _ = logFile.Close() }

// /* ---------- 汎用 ---------- */

// func move(pattern, dst string) error {
// 	list, err := filepath.Glob(pattern)
// 	if err != nil {
// 		return err
// 	}
// 	if len(list) == 0 {
// 		return nil
// 	}
// 	if err := os.MkdirAll(dst, 0o755); err != nil {
// 		return err
// 	}
// 	for _, src := range list {
// 		dest := filepath.Join(dst, filepath.Base(src))
// 		logger.Printf("mv %s → %s", src, dest)
// 		if err := os.Rename(src, dest); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func isJPEG(p string) bool {
// 	ext := strings.ToLower(filepath.Ext(p))
// 	return ext == ".jpg" || ext == ".jpeg"
// }

// /* ---------- EXIF 書き込み ---------- */

// func setCreateDate(path string, t time.Time) error {
// 	jmp := jpegstructure.NewJpegMediaParser()
// 	intfc, err := jmp.ParseFile(path)
// 	if err != nil {
// 		return err
// 	}
// 	sl := intfc.(*jpegstructure.SegmentList)

// 	rootIb, err := sl.ConstructExifBuilder()
// 	if err != nil {
// 		return err
// 	}
// 	exifIb, err := exif.GetOrCreateIbFromRootIb(rootIb, "IFD/Exif")
// 	if err != nil {
// 		return err
// 	}
// 	stamp := exifcommon.ExifFullTimestampString(t)
// 	if err := exifIb.SetStandardWithName("CreateDate", stamp); err != nil {
// 		return err
// 	}

// 	if err := sl.SetExif(rootIb); err != nil {
// 		return err
// 	}

// 	var buf bytes.Buffer
// 	if err := sl.Write(&buf); err != nil {
// 		return err
// 	}
// 	return os.WriteFile(path, buf.Bytes(), 0o644)
// }

// func mirrorFSDate(path string) error {
// 	info, err := os.Stat(path)
// 	if err != nil {
// 		return err
// 	}
// 	return setCreateDate(path, info.ModTime())
// }

// /* ---------- CreateDate → ファイル名 ---------- */

// func renameByCreateDate(root string) error {
// 	return filepath.WalkDir(root, func(p string, d fs.DirEntry, err error) error {
// 		if err != nil || d.IsDir() || !isJPEG(p) {
// 			return err
// 		}

// 		jmp := jpegstructure.NewJpegMediaParser()
// 		intfc, err := jmp.ParseFile(p)
// 		if err != nil { return err }

// 		sl := intfc.(*jpegstructure.SegmentList)
// 		_, _, tags, err := sl.DumpExif()
// 		if err != nil { return err }

// 		var stamp string
// 		for _, t := range tags {
// 			if t.TagName == "CreateDate" {
// 				stamp = t.FormattedFirst
// 				break
// 			}
// 		}
// 		if stamp == "" { return nil }

// 		const exifLayout = "2006:01:02 15:04:05"
// 		ts, err := time.ParseInLocation(exifLayout, stamp, time.Local)
// 		if err != nil { return err }

// 		newName := ts.Format("20060102150405") + filepath.Ext(p)
// 		newPath := filepath.Join(filepath.Dir(p), newName)
// 		if p == newPath {
// 			return nil
// 		}
// 		if _, err := os.Stat(newPath); err == nil { // 既に存在 → _01 など付与
// 			newName = ts.Format("20060102150405_02") + filepath.Ext(p)
// 			newPath = filepath.Join(filepath.Dir(p), newName)
// 		}
// 		logger.Printf("rename %s → %s", p, newName)
// 		return os.Rename(p, newPath)
// 	})
// }


// /* ---------- main ---------- */

// func main() {
// 	initLog()
// 	defer closeLog()

// 	// 1) *.mp4 を退避
// 	if err := move("*.mp4", dirMovieEscaping); err != nil {
// 		log.Fatal(err)
// 	}

// 	// 2-A) 1-2_create_date_setting : FileCreateDate → CreateDate
// 	filepath.WalkDir(dirCreateDateSetting, func(p string, d fs.DirEntry, err error) error {
// 		if err == nil && !d.IsDir() && isJPEG(p) {
// 			_ = mirrorFSDate(p)
// 		}
// 		return nil
// 	})

// 	// 2-B) 1-1_movie_escaping : FileModifyDate → CreateDate
// 	filepath.WalkDir(dirMovieEscaping, func(p string, d fs.DirEntry, err error) error {
// 		if err == nil && !d.IsDir() && isJPEG(p) {
// 			_ = mirrorFSDate(p)
// 		}
// 		return nil
// 	})

// 	// 3) mp4 を戻す
// 	_ = move(filepath.Join(dirMovieEscaping, "*.mp4"), ".")

// 	// 4) CreateDate でリネーム
// 	if err := renameByCreateDate("."); err != nil {
// 		log.Fatal(err)
// 	}

// 	// 5) 画像を戻す
// 	for _, pat := range []string{"*.jpg", "*.png", "*.webp"} {
// 		_ = move(filepath.Join(dirCreateDateSetting, pat), ".")
// 	}

// 	// pause
// 	fmt.Print("Input 'y' to terminate this process......")
// 	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
// }
