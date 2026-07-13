package main

import (
	"bytes"
	"testing"

	usecases "github.com/landmaster135/devbox/internal/pdf_merger/usecases"
)

func TestParseFlags_Recursive(t *testing.T) {
	var stderr bytes.Buffer

	opts, err := parseFlags([]string{"-operation", usecases.OperationMergeIntoNew, "-src-dir", "images", "-output-dir", "output", "-recursive"}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags()でエラーが発生: %v", err)
	}

	if opts.Operation != usecases.OperationMergeIntoNew {
		t.Errorf("Operation: 期待 %s, 実際 %s", usecases.OperationMergeIntoNew, opts.Operation)
	}
	if opts.Dir != "images" {
		t.Errorf("Dir: 期待 %s, 実際 %s", "images", opts.Dir)
	}
	if !opts.Recursive {
		t.Errorf("Recursive: 期待 true, 実際 false")
	}
	if opts.OutputDir != "output" {
		t.Errorf("OutputDir: 期待 %s, 実際 %s", "output", opts.OutputDir)
	}
}

func TestParseFlags_RecursiveDefault(t *testing.T) {
	var stderr bytes.Buffer

	opts, err := parseFlags([]string{"-operation", usecases.OperationMergeIntoNew, "-src-dir", "images", "-output-dir", "output"}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags()でエラーが発生: %v", err)
	}

	if opts.Recursive {
		t.Errorf("Recursive: 期待 false, 実際 true")
	}
}

func TestParseFlags_OutFlagRemoved(t *testing.T) {
	var stderr bytes.Buffer

	_, err := parseFlags([]string{"-operation", usecases.OperationMergeIntoNew, "-src-dir", "images", "-out", "output.pdf"}, &stderr)
	if err == nil {
		t.Fatal("廃止された-outフラグではエラーが期待されます")
	}
}

func TestParseFlags_OutputDirRequired(t *testing.T) {
	var stderr bytes.Buffer

	_, err := parseFlags([]string{"-operation", usecases.OperationMergeIntoNew, "-src-dir", "images"}, &stderr)
	if err == nil {
		t.Fatal("-output-dir未指定ではエラーが期待されます")
	}
}

func TestParseFlags_DirFlagRemoved(t *testing.T) {
	var stderr bytes.Buffer

	_, err := parseFlags([]string{"-operation", usecases.OperationMergeIntoNew, "-dir", "images", "-output-dir", "output"}, &stderr)
	if err == nil {
		t.Fatal("廃止された-dirフラグではエラーが期待されます")
	}
}

func TestParseFlags_OperationValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "Operation未指定",
			args: []string{"-src-dir", "images", "-output-dir", "output"},
		},
		{
			name: "未知Operation",
			args: []string{"-operation", "unknown", "-src-dir", "images", "-output-dir", "output"},
		},
		{
			name: "MergeIntoNewでAdd指定",
			args: []string{"-operation", usecases.OperationMergeIntoNew, "-src-dir", "images", "-add", "existing.pdf", "-output-dir", "output"},
		},
		{
			name: "MergeIntoNewでSrcFile指定",
			args: []string{"-operation", usecases.OperationMergeIntoNew, "-src-dir", "images", "-src-file", "input.pdf", "-output-dir", "output"},
		},
		{
			name: "AddIntoExistでAdd未指定",
			args: []string{"-operation", usecases.OperationAddIntoExist, "-src-dir", "images", "-output-dir", "output"},
		},
		{
			name: "AddIntoExistでSrcFile指定",
			args: []string{"-operation", usecases.OperationAddIntoExist, "-src-dir", "images", "-add", "existing.pdf", "-src-file", "input.pdf", "-output-dir", "output"},
		},
		{
			name: "ExtractImagesでSrcFile未指定",
			args: []string{"-operation", usecases.OperationExtractImages, "-output-dir", "output"},
		},
		{
			name: "ExtractImagesでAdd指定",
			args: []string{"-operation", usecases.OperationExtractImages, "-src-file", "input.pdf", "-add", "existing.pdf", "-output-dir", "output"},
		},
		{
			name: "Extract廃止フラグ指定",
			args: []string{"-operation", usecases.OperationExtractImages, "-extract", "input.pdf", "-output-dir", "output"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer

			_, err := parseFlags(tt.args, &stderr)
			if err == nil {
				t.Fatal("operationバリデーションエラーが期待されます")
			}
		})
	}
}

func TestParseFlags_OperationSpecificOptions(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedOperation string
		expectedAdd       string
		expectedExtract   string
	}{
		{
			name:              "既存PDFに画像を追加_Normal",
			args:              []string{"-operation", usecases.OperationAddIntoExist, "-src-dir", "images", "-add", "existing.pdf", "-output-dir", "output"},
			expectedOperation: usecases.OperationAddIntoExist,
			expectedAdd:       "existing.pdf",
		},
		{
			name:              "PDFから画像を抽出_Normal",
			args:              []string{"-operation", usecases.OperationExtractImages, "-src-file", "input.pdf", "-output-dir", "output"},
			expectedOperation: usecases.OperationExtractImages,
			expectedExtract:   "input.pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer

			opts, err := parseFlags(tt.args, &stderr)
			if err != nil {
				t.Fatalf("parseFlags()でエラーが発生: %v", err)
			}
			if opts.Operation != tt.expectedOperation {
				t.Errorf("Operation: 期待 %s, 実際 %s", tt.expectedOperation, opts.Operation)
			}
			if opts.Add != tt.expectedAdd {
				t.Errorf("Add: 期待 %s, 実際 %s", tt.expectedAdd, opts.Add)
			}
			if opts.Extract != tt.expectedExtract {
				t.Errorf("Extract: 期待 %s, 実際 %s", tt.expectedExtract, opts.Extract)
			}
		})
	}
}
