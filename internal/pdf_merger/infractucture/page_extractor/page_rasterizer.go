package page_extractor

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	imagedraw "image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	webp "github.com/gen2brain/webp"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	filter "github.com/pdfcpu/pdfcpu/pkg/filter"
	pdfcpu "github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/matrix"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"golang.org/x/image/draw"
	"golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

var supportedFormats = map[string]struct{}{
	"jpg":  {},
	"jpeg": {},
	"png":  {},
	"tiff": {},
	"webp": {},
}

const defaultDPI = 144.0
const matrixTolerance = 1e-6

// RasterizeOptions represents configuration for page rasterization.
type RasterizeOptions struct {
	OutputDir string
	Format    string
	StartPage int
	EndPage   int
	DPI       float64
}

// PageRasterizer converts PDF pages into raster images.
type PageRasterizer struct{}

// NewPageRasterizer creates a new PageRasterizer instance.
func NewPageRasterizer() *PageRasterizer {
	return &PageRasterizer{}
}

// Rasterize renders PDF pages into images and returns the output file paths.
func (r *PageRasterizer) Rasterize(pdfPath string, opts RasterizeOptions) ([]string, error) {
	dpi := opts.DPI
	if dpi <= 0 {
		dpi = defaultDPI
	}

	format := strings.ToLower(opts.Format)
	if format == "" {
		format = "png"
	}
	if _, ok := supportedFormats[format]; !ok {
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output dir: %w", err)
	}

	ctx, err := api.ReadContextFile(pdfPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read pdf: %w", err)
	}

	if err := ctx.EnsurePageCount(); err != nil {
		return nil, fmt.Errorf("failed to determine page count: %w", err)
	}

	pageCount := ctx.PageCount
	if pageCount == 0 {
		return nil, fmt.Errorf("pdf has no pages")
	}

	start := opts.StartPage
	end := opts.EndPage
	if start <= 0 {
		start = 1
	}
	if end <= 0 || end > pageCount {
		end = pageCount
	}
	if start > end {
		return nil, fmt.Errorf("invalid page range: start=%d end=%d", start, end)
	}

	digits := len(strconv.Itoa(pageCount))
	if digits < 4 {
		digits = 4
	}

	baseName := strings.TrimSuffix(filepath.Base(pdfPath), filepath.Ext(pdfPath))
	scale := dpi / 72.0

	outputs := make([]string, 0, end-start+1)

	for pageNum := start; pageNum <= end; pageNum++ {
		pagePath := filepath.Join(opts.OutputDir, fmt.Sprintf("%s_%0*d.%s", baseName, digits, pageNum, format))

		if err := r.renderPage(ctx, pageNum, scale, pagePath, format); err != nil {
			return nil, fmt.Errorf("failed to rasterize page %d: %w", pageNum, err)
		}

		outputs = append(outputs, pagePath)
	}

	return outputs, nil
}

func (r *PageRasterizer) renderPage(ctx *model.Context, pageNum int, scale float64, outputPath, format string) error {
	pageDict, _, attrs, err := ctx.PageDict(pageNum, true)
	if err != nil {
		return fmt.Errorf("failed to obtain page dict: %w", err)
	}

	mediaBox := attrs.MediaBox
	if mediaBox == nil {
		return fmt.Errorf("page %d has no MediaBox", pageNum)
	}

	widthPt := mediaBox.Width()
	heightPt := mediaBox.Height()
	widthPx := int(math.Ceil(widthPt * scale))
	heightPx := int(math.Ceil(heightPt * scale))
	if widthPx <= 0 || heightPx <= 0 {
		return fmt.Errorf("invalid page dimensions: %dx%d", widthPx, heightPx)
	}

	canvas := image.NewRGBA(image.Rect(0, 0, widthPx, heightPx))
	imagedraw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.White}, image.Point{}, imagedraw.Src)

	placements, err := r.collectPlacements(ctx, pageDict, attrs)
	if err != nil {
		return err
	}

	for _, placement := range placements {
		if err := r.drawPlacement(ctx, canvas, placement, scale, heightPt); err != nil {
			return err
		}
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	switch format {
	case "jpg", "jpeg":
		err = jpeg.Encode(f, canvas, &jpeg.Options{Quality: 90})
	case "png":
		err = png.Encode(f, canvas)
	case "tiff":
		err = tiff.Encode(f, canvas, nil)
	case "webp":
		err = webp.Encode(f, canvas, webp.Options{Lossless: true})
	}

	if err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}
	return nil
}

func (r *PageRasterizer) collectPlacements(ctx *model.Context, pageDict types.Dict, attrs *model.InheritedPageAttrs) ([]imagePlacement, error) {
	resources := attrs.Resources
	if resources == nil {
		return nil, nil
	}

	xObjectObj := resources.DictEntry("XObject")
	if xObjectObj == nil {
		return nil, nil
	}

	xObjectDict, err := ctx.DereferenceDict(xObjectObj)
	if err != nil {
		return nil, fmt.Errorf("failed to dereference XObject dict: %w", err)
	}

	imageStreams := map[string]imageResource{}
	for name, obj := range xObjectDict {
		objNr := 0
		if indRef, ok := obj.(types.IndirectRef); ok {
			objNr = indRef.ObjectNumber.Value()
		}
		sd, ok, err := ctx.DereferenceStreamDict(obj)
		if err != nil || !ok || sd == nil {
			continue
		}
		if subtype := sd.NameEntry("Subtype"); subtype == nil || *subtype != "Image" {
			continue
		}
		imageStreams[name] = imageResource{
			stream:       sd,
			objectNumber: objNr,
		}
	}

	if len(imageStreams) == 0 {
		return nil, nil
	}

	content, err := ctx.PageContent(pageDict)
	if err != nil && err != model.ErrNoContent {
		return nil, fmt.Errorf("failed to read page content: %w", err)
	}

	tokens := tokenizeContent(content)
	placements := parsePlacements(tokens)

	filtered := make([]imagePlacement, 0, len(placements))
	for _, p := range placements {
		if res, ok := imageStreams[p.name]; ok {
			p.resource = res
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

type imageResource struct {
	stream       *types.StreamDict
	objectNumber int
}

type imagePlacement struct {
	name     string
	matrix   matrix.Matrix
	resource imageResource
}

func (r *PageRasterizer) drawPlacement(ctx *model.Context, canvas *image.RGBA, place imagePlacement, scale float64, pageHeight float64) error {
	m := place.matrix

	if !(math.Abs(m[0][1]) < matrixTolerance && math.Abs(m[1][0]) < matrixTolerance) {
		return fmt.Errorf("unsupported rotation or shear in image %s", place.name)
	}

	minX := m[2][0]
	minY := m[2][1]
	maxX := minX + m[0][0]
	maxY := minY + m[1][1]

	if maxX < minX {
		minX, maxX = maxX, minX
	}
	if maxY < minY {
		minY, maxY = maxY, minY
	}

	left := int(math.Round(minX * scale))
	right := int(math.Round(maxX * scale))
	top := int(math.Round((pageHeight - maxY) * scale))
	bottom := int(math.Round((pageHeight - minY) * scale))

	if left == right || top == bottom {
		return nil
	}

	bounds := image.Rect(left, top, right, bottom)

	if err := place.resource.stream.Decode(); err != nil && err != filter.ErrUnsupportedFilter {
		return fmt.Errorf("failed to decode image stream %s: %w", place.name, err)
	}

	reader, _, err := pdfcpu.RenderImage(ctx.XRefTable, place.resource.stream, false, place.name, place.resource.objectNumber)
	if err != nil {
		return fmt.Errorf("failed to render image %s: %w", place.name, err)
	}

	buf, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	img, _, err := image.Decode(bytes.NewReader(buf))
	if err != nil {
		return err
	}

	dst := image.NewRGBA(bounds.Sub(bounds.Min))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	imagedraw.Draw(canvas, bounds, dst, image.Point{}, imagedraw.Over)

	return nil
}

func tokenizeContent(content []byte) []string {
	tokens := []string{}
	i := 0
	for i < len(content) {
		b := content[i]
		switch {
		case isWhitespace(b):
			i++
		case b == '%':
			for i < len(content) && content[i] != '\n' && content[i] != '\r' {
				i++
			}
		case b == '(':
			i = skipLiteralString(content, i+1)
		case b == '<':
			if i+1 < len(content) && content[i+1] == '<' {
				i = skipDictionary(content, i+2)
			} else {
				i = skipHexString(content, i+1)
			}
		case b == '[' || b == ']':
			tokens = append(tokens, string(b))
			i++
		case b == '/':
			start := i + 1
			i++
			for i < len(content) && !isDelimiter(content[i]) && !isWhitespace(content[i]) {
				i++
			}
			tokens = append(tokens, "/"+string(content[start:i]))
		default:
			start := i
			for i < len(content) && !isDelimiter(content[i]) && !isWhitespace(content[i]) {
				i++
			}
			if start < i {
				tokens = append(tokens, string(content[start:i]))
			} else {
				i++
			}
		}
	}
	return tokens
}

func parsePlacements(tokens []string) []imagePlacement {
	placements := []imagePlacement{}
	var stack []string
	current := matrix.IdentMatrix
	var states []matrix.Matrix
	skipInline := false

	for _, token := range tokens {
		if skipInline {
			if token == "EI" {
				skipInline = false
			}
			continue
		}

		switch token {
		case "q":
			states = append(states, current)
			stack = stack[:0]
		case "Q":
			if len(states) > 0 {
				current = states[len(states)-1]
				states = states[:len(states)-1]
			} else {
				current = matrix.IdentMatrix
			}
			stack = stack[:0]
		case "cm":
			if len(stack) < 6 {
				stack = stack[:0]
				continue
			}
			vals := make([]float64, 6)
			for i := 5; i >= 0; i-- {
				val, err := strconv.ParseFloat(stack[len(stack)-1], 64)
				stack = stack[:len(stack)-1]
				if err != nil {
					val = 0
				}
				vals[i] = val
			}
			current = current.Multiply(matrixFromComponents(vals))
		case "Do":
			if len(stack) < 1 {
				stack = stack[:0]
				continue
			}
			name := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if strings.HasPrefix(name, "/") {
				name = name[1:]
			}
			placements = append(placements, imagePlacement{name: name, matrix: current})
			stack = stack[:0]
		case "BI":
			skipInline = true
			stack = stack[:0]
		default:
			stack = append(stack, token)
			if len(stack) > 32 {
				stack = stack[len(stack)-32:]
			}
		}
	}

	return placements
}

func matrixFromComponents(vals []float64) matrix.Matrix {
	m := matrix.IdentMatrix
	m[0][0] = vals[0]
	m[0][1] = vals[1]
	m[1][0] = vals[2]
	m[1][1] = vals[3]
	m[2][0] = vals[4]
	m[2][1] = vals[5]
	return m
}

func isWhitespace(b byte) bool {
	return b == 0 || b == 9 || b == 10 || b == 12 || b == 13 || b == 32
}

func isDelimiter(b byte) bool {
	switch b {
	case '(', ')', '<', '>', '[', ']', '{', '}', '/', '%':
		return true
	}
	return false
}

func skipLiteralString(content []byte, i int) int {
	depth := 1
	for i < len(content) && depth > 0 {
		switch content[i] {
		case '\\':
			i += 2
			continue
		case '(':
			depth++
		case ')':
			depth--
		}
		i++
	}
	return i
}

func skipHexString(content []byte, i int) int {
	for i < len(content) && content[i] != '>' {
		i++
	}
	if i < len(content) {
		i++
	}
	return i
}

func skipDictionary(content []byte, i int) int {
	depth := 1
	for i < len(content) && depth > 0 {
		if content[i] == '<' && i+1 < len(content) && content[i+1] == '<' {
			depth++
			i += 2
			continue
		}
		if content[i] == '>' && i+1 < len(content) && content[i+1] == '>' {
			depth--
			i += 2
			continue
		}
		if content[i] == '(' {
			i = skipLiteralString(content, i+1)
			continue
		}
		if content[i] == '<' {
			i = skipHexString(content, i+1)
			continue
		}
		i++
	}
	return i
}
