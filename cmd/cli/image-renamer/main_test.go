package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseFlags_ExtensionsDefault_Normal(t *testing.T) {
	var stderr bytes.Buffer
	config, err := parseFlags([]string{"-prefix", "article01", "-name"}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags() が予期せず失敗しました: %v", err)
	}

	if config.Extensions != nil {
		t.Errorf("extensions未指定時はnilを期待します。実際: %v", config.Extensions)
	}
}

func TestParseFlags_WithoutPrefix_Normal(t *testing.T) {
	var stderr bytes.Buffer
	config, err := parseFlags([]string{"-name"}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags() が予期せず失敗しました: %v", err)
	}

	if config.Prefix != "" {
		t.Errorf("prefix未指定時は空文字を期待します。実際: %q", config.Prefix)
	}
}

func TestParseFlags_ExtensionsCustom_Normal(t *testing.T) {
	var stderr bytes.Buffer
	config, err := parseFlags([]string{
		"-prefix", "article01",
		"-name",
		"-extensions", ".jpg, png ,HEIC",
	}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags() が予期せず失敗しました: %v", err)
	}

	expected := []string{".jpg", "png", "HEIC"}
	if !reflect.DeepEqual(config.Extensions, expected) {
		t.Errorf("extensionsが期待値と異なります。期待: %v, 実際: %v", expected, config.Extensions)
	}
}

func TestParseFlags_ExtensionsEmpty_Normal(t *testing.T) {
	var stderr bytes.Buffer
	config, err := parseFlags([]string{
		"-prefix", "article01",
		"-name",
		"-extensions", " , , ",
	}, &stderr)
	if err != nil {
		t.Fatalf("parseFlags() が予期せず失敗しました: %v", err)
	}

	if config.Extensions != nil {
		t.Errorf("実質空のextensions指定時はnilを期待します。実際: %v", config.Extensions)
	}
}
