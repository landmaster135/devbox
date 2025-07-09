// main.go  (Go 1.22+, go-jpeg-image-structure/v2, go-exif/v3)
package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	exif "github.com/dsoprea/go-exif/v3"
	jpegstructure "github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/gen2brain/webp"
)

// -----------------------------------------------------------------------------
// Image conversion (gen2brain/webp)
// -----------------------------------------------------------------------------

func convertImageToWebp(srcPath string, logger *log.Logger) error {
	// Open and decode source image (JPEG/PNG)
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer srcFile.Close()

	img, _, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	destPath := strings.TrimSuffix(srcPath, filepath.Ext(srcPath)) + ".webp"
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("create dest: %w", err)
	}
	defer destFile.Close()

	opts := &webp.Options{Quality: 75, Lossless: false, Method: webp.DefaultMethod}
	if err := webp.Encode(destFile, img, *opts); err != nil {
		return fmt.Errorf("encode webp: %w", err)
	}

	logger.Printf("+ webp.Encode %s -> %s (q=%d)", filepath.Base(srcPath), filepath.Base(destPath), opts.Quality)
	return nil
}

func convertImagesToWebp(dir string, logger *log.Logger) error {
	patterns := []string{"*.jpg", "*.jpeg", "*.png"}
	converted := 0
	for _, pat := range patterns {
		files, _ := filepath.Glob(filepath.Join(dir, pat))
		for _, f := range files {
			logger.Printf("<Source File Name: %s>", filepath.Base(f))
			if err := convertImageToWebp(f, logger); err != nil {
				logger.Print(err)
			} else {
				converted++
			}
			logger.Print("----------------")
		}
	}
	logger.Printf("%d .webp images converted.", converted)
	return nil
}

// -----------------------------------------------------------------------------
// EXIF helpers (unchanged)
// -----------------------------------------------------------------------------

func extractExifFromJpeg(path string) ([]byte, error) {
	mp := jpegstructure.NewJpegMediaParser()
	intfc, err := mp.ParseFile(path)
	if err != nil {
		return nil, err
	}
	sl := intfc.(*jpegstructure.SegmentList)
	_, rawExif, err := sl.Exif()
	if err != nil {
		if errors.Is(err, exif.ErrNoExif) {
			return nil, nil
		}
		return nil, err
	}
	return rawExif, nil
}

func embedExifIntoWebp(webpPath string, exifPayload []byte) error {
	if len(exifPayload) == 0 {
		return nil
	}
	const fourCCEXIF = "EXIF"

	data, err := os.ReadFile(webpPath)
	if err != nil {
		return err
	}
	if !bytes.HasPrefix(data, []byte("RIFF")) || !bytes.Contains(data[:16], []byte("WEBP")) {
		return fmt.Errorf("%s is not a WebP file", webpPath)
	}

	cur := 12 // after RIFF header
	var cleaned bytes.Buffer
	cleaned.Write(data[:cur])

	for cur < len(data) {
		if cur+8 > len(data) {
			return fmt.Errorf("truncated chunk in %s", webpPath)
		}
		fourCC := string(data[cur : cur+4])
		size := int(binary.LittleEndian.Uint32(data[cur+4 : cur+8]))
		chunkEnd := cur + 8 + size
		if size%2 == 1 {
			chunkEnd++
		}
		if chunkEnd > len(data) {
			return fmt.Errorf("invalid chunk size in %s", webpPath)
		}
		if fourCC != fourCCEXIF {
			cleaned.Write(data[cur:chunkEnd])
		}
		cur = chunkEnd
	}

	var exifChunk bytes.Buffer
	exifChunk.WriteString(fourCCEXIF)
	if err := binary.Write(&exifChunk, binary.LittleEndian, uint32(len(exifPayload))); err != nil {
		return err
	}
	exifChunk.Write(exifPayload)
	if len(exifPayload)%2 == 1 {
		exifChunk.WriteByte(0)
	}

	cleaned.Write(exifChunk.Bytes())
	final := cleaned.Bytes()
	binary.LittleEndian.PutUint32(final[4:8], uint32(len(final)-8))

	return os.WriteFile(webpPath, final, 0644)
}

func copyExifFromSrc(dir, srcExt string, logger *log.Logger) error {
	destExt := ".webp"
	files, _ := filepath.Glob(filepath.Join(dir, "*"+destExt))
	fixed := 0
	for _, webpFile := range files {
		src := strings.TrimSuffix(webpFile, destExt) + srcExt
		if _, err := os.Stat(src); err != nil {
			continue
		}
		logger.Printf("Processing EXIF copy: %s -> %s", src, webpFile)
		exifPayload, err := extractExifFromJpeg(src)
		if err != nil {
			logger.Printf("extract error: %v", err)
			continue
		}
		if len(exifPayload) == 0 {
			logger.Printf("No EXIF found in %s", src)
			continue
		}
		if err := embedExifIntoWebp(webpFile, exifPayload); err != nil {
			logger.Printf("embed error: %v", err)
			continue
		}
		fixed++
	}
	logger.Printf("EXIF info of %d %s images are fixed with %s images.", fixed, destExt, srcExt)
	return nil
}

// -----------------------------------------------------------------------------
// Utility: move originals
// -----------------------------------------------------------------------------

func moveOriginalFiles(dir string, logger *log.Logger) error {
	destDir := filepath.Join(dir, "5_original_files")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	patterns := []string{"*_original", "*.jpg", "*.jpeg", "*.png"}
	for _, pat := range patterns {
		files, _ := filepath.Glob(filepath.Join(dir, pat))
		for _, f := range files {
			dest := filepath.Join(destDir, filepath.Base(f))
			if err := os.Rename(f, dest); err != nil {
				logger.Printf("move %s: %v", f, err)
			} else {
				logger.Printf("moved %s to %s", f, dest)
			}
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// main
// -----------------------------------------------------------------------------

func main() {
	logFileName := fmt.Sprintf("log_%s.txt", time.Now().Format("20060102"))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("open log file: %v", err)
	}
	defer logFile.Close()

	logger := log.New(io.MultiWriter(os.Stdout, logFile), "", log.LstdFlags)

	logger.Println("3-1_convert_webp_from_jpg ++++++++++++++++++++++++++++++++")
	convertImagesToWebp(".", logger)

	logger.Println("4-1_copy_exif_from_jpg ++++++++++++++++++++++++++++++++")
	copyExifFromSrc(".", ".jpg", logger)

	logger.Println("4-1-b_copy_exif_from_png ++++++++++++++++++++++++++++++++")
	// Skipped (PNG rarely stores EXIF). Implement similarly if needed.

	logger.Println("5-1_escape_original_files ++++++++++++++++++++++++++++++++")
	moveOriginalFiles(".", logger)
}
