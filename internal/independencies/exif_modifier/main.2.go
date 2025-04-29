// main.go  (Go 1.22+, go-jpeg-image-structure/v2, go-exif/v3)
package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"
// 	"time"

// 	exif "github.com/dsoprea/go-exif/v3"
// 	exifcommon "github.com/dsoprea/go-exif/v3/common"
// 	jpeg "github.com/dsoprea/go-jpeg-image-structure/v2"
// )

// // ---------- 必要なユーティリティ ----------

// // must : err が nil なら値だけ返し、nil でなければ panic
// func must[T any](v T, err error) T {
// 	if err != nil {
// 		panic(err)
// 	}
// 	return v
// }

// // writeExif : 1 枚の JPEG に 3 つの日時タグを書き込む
// func writeExif(jpegPath, exifDate string) error {
// 	parser := jpeg.NewJpegMediaParser()
// 	intfc, err := parser.ParseFile(jpegPath)
// 	if err != nil {
// 		return err
// 	}
// 	sl := intfc.(*jpeg.SegmentList)

// 	// 既存 EXIF を読み込み／なければ空で作成
// 	rootIb, err := sl.ConstructExifBuilder()
// 	if err != nil { // 画像に EXIF が無い
// 		im := must(exifcommon.NewIfdMappingWithStandard())
// 		ti := exif.NewTagIndex()
// 		rootIb = exif.NewIfdBuilder(im, ti,
// 			exifcommon.IfdStandardIfdIdentity,
// 			exifcommon.TestDefaultByteOrder) // = Intel endian
// 	}

// 	// IFD0 → DateTime
// 	ifd0, err := exif.GetOrCreateIbFromRootIb(rootIb, "IFD")
// 	if err != nil {
// 		return err
// 	}
// 	if err := ifd0.SetStandardWithName("DateTime", exifDate); err != nil {
// 		return err
// 	}

// 	// ExifIFD → DateTimeOriginal, CreateDate
// 	exifIfd, err := exif.GetOrCreateIbFromRootIb(rootIb, "IFD/Exif")
// 	if err != nil {
// 		return err
// 	}
// 	for _, tag := range []string{"DateTimeOriginal", "CreateDate"} {
// 		if err := exifIfd.SetStandardWithName(tag, exifDate); err != nil {
// 			return err
// 		}
// 	}

// 	// EXIF を JPEG へ反映
// 	if err := sl.SetExif(rootIb); err != nil { // :contentReference[oaicite:0]{index=0}
// 		return err
// 	}

// 	out := must(os.Create(jpegPath))
// 	defer out.Close()
// 	return sl.Write(out)
// }

// // ---------- メイン処理 ----------

// func main() {
// 	relDir := "."
// 	now := time.Now()
// 	logFile := must(os.OpenFile(
// 		filepath.Join(relDir, fmt.Sprintf("log_%s.txt", now.Format("20060102"))),
// 		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644))
// 	defer logFile.Close()

// 	writeLog := func(s string) { _, _ = logFile.WriteString(s) }

// 	reader := bufio.NewReader(os.Stdin)
// 	var prefix string
// 	for {
// 		fmt.Print("Input prefix of pictures (yyyyMMddHH): ")
// 		in, _ := reader.ReadString('\n')
// 		in = strings.TrimSpace(in)
// 		if len(in) == 10 {
// 			prefix = in
// 			break
// 		}
// 	}
// 	targetTime := must(time.Parse("2006010215", prefix))
// 	exifDate := targetTime.Format("2006:01:02 15:04:05")
// 	fmt.Println(exifDate)
// 	writeLog(exifDate + "\n")

// 	files := must(os.ReadDir(relDir))
// 	counter := 1
// 	for _, entry := range files {
// 		name := entry.Name()
// 		if entry.IsDir() || skip(name) {
// 			continue
// 		}
// 		jpegPath := filepath.Join(relDir, name)
// 		if err := writeExif(jpegPath, exifDate); err != nil {
// 			fmt.Fprintln(os.Stderr, "EXIF update error:", err)
// 			continue
// 		}
// 		newName := fmt.Sprintf("%s%04d%s", prefix, counter, filepath.Ext(name))
// 		_ = os.Rename(jpegPath, filepath.Join(relDir, newName))
// 		fmt.Println("LOG-INFO", newName)
// 		writeLog("LOG-INFO," + newName + "\n")
// 		counter++
// 	}
// }

// func skip(n string) bool {
// 	l := strings.ToLower(n)
// 	switch filepath.Ext(l) {
// 	case ".txt", ".exe", ".bat", ".ps1":
// 		return true
// 	}
// 	return strings.Contains(l, ".jpg_original")
// }
