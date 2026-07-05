package envloader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoad_Normal(t *testing.T) {
	withWorkingDirectory(t, map[string]string{
		".env": "ENV_LOADER_TEST_VALUE=from_dotenv\n",
	})

	values, err := Load([]string{"ENV_LOADER_TEST_VALUE"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if values["ENV_LOADER_TEST_VALUE"] != "from_dotenv" {
		t.Fatalf("Load() value = %q, want %q", values["ENV_LOADER_TEST_VALUE"], "from_dotenv")
	}
}

func TestLoad_NormalWithExistingEnvironment(t *testing.T) {
	t.Setenv("ENV_LOADER_EXISTING_VALUE", "from_environment")
	withWorkingDirectory(t, map[string]string{
		".env": "ENV_LOADER_EXISTING_VALUE=from_dotenv\n",
	})

	values, err := Load([]string{"ENV_LOADER_EXISTING_VALUE"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if values["ENV_LOADER_EXISTING_VALUE"] != "from_environment" {
		t.Fatalf("Load() value = %q, want %q", values["ENV_LOADER_EXISTING_VALUE"], "from_environment")
	}
}

func TestLoad_NormalWithoutDotenv(t *testing.T) {
	t.Setenv("ENV_LOADER_OS_ONLY_VALUE", "from_environment")
	withWorkingDirectory(t, nil)

	values, err := Load([]string{"ENV_LOADER_OS_ONLY_VALUE"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if values["ENV_LOADER_OS_ONLY_VALUE"] != "from_environment" {
		t.Fatalf("Load() value = %q, want %q", values["ENV_LOADER_OS_ONLY_VALUE"], "from_environment")
	}
}

func TestLoad_NormalWithEmptyKeys(t *testing.T) {
	withWorkingDirectory(t, nil)

	values, err := Load(nil)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(values) != 0 {
		t.Fatalf("Load() values length = %d, want 0", len(values))
	}
}

func TestLoad_ErrorWhenKeyIsMissing(t *testing.T) {
	withWorkingDirectory(t, nil)

	_, err := Load([]string{"ENV_LOADER_MISSING_VALUE"})
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "ENV_LOADER_MISSING_VALUE") {
		t.Fatalf("Load() error = %v, want missing key name", err)
	}
}

func TestLoad_ErrorWhenValueIsEmpty(t *testing.T) {
	withWorkingDirectory(t, map[string]string{
		".env": "ENV_LOADER_EMPTY_VALUE=\n",
	})

	_, err := Load([]string{"ENV_LOADER_EMPTY_VALUE"})
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "ENV_LOADER_EMPTY_VALUE") {
		t.Fatalf("Load() error = %v, want empty key name", err)
	}
}

func TestLoad_ErrorWhenDotenvIsInvalid(t *testing.T) {
	withWorkingDirectory(t, map[string]string{
		".env": "ENV_LOADER_INVALID LINE\n",
	})

	_, err := Load([]string{"ENV_LOADER_INVALID"})
	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if !strings.Contains(err.Error(), ".envの読み込みに失敗しました") {
		t.Fatalf("Load() error = %v, want dotenv load error", err)
	}
}

func withWorkingDirectory(t *testing.T, files map[string]string) {
	t.Helper()

	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("failed to write test file %s: %v", path, err)
		}
	}

	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd() error = %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("os.Chdir() error = %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})
}
