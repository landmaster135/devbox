package config

import (
    "os"
    "path/filepath"
    "testing"
)

func TestParseHexColor(t *testing.T) {
    rgb, err := ParseHexColor("#1a2b3c")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if rgb[0] <= 0.1 || rgb[0] >= 0.11 {
        t.Fatalf("unexpected red value: %f", rgb[0])
    }
    if _, err := ParseHexColor("zzz"); err == nil {
        t.Fatalf("expected error for invalid hex")
    }
}

func TestConfigNormalise(t *testing.T) {
    dir := t.TempDir()
    input := filepath.Join(dir, "sample.png")
    if err := os.WriteFile(input, []byte("dummy"), 0o644); err != nil {
        t.Fatalf("failed to setup input: %v", err)
    }

    cfg := Config{
        InputPath: input,
        Mode:      FilterModeGrayscale,
        Strength:  0.5,
    }

    if err := cfg.Validate(); err != nil {
        t.Fatalf("unexpected validation error: %v", err)
    }

    cfg.Normalise()
    if cfg.OutputPath == "" {
        t.Fatalf("expected output path to be set")
    }
    if cfg.TintHex == "" {
        t.Fatalf("expected tint hex to be set")
    }
}
