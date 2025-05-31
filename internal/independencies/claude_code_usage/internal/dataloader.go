package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// DataLoader handles loading and parsing JSONL files
type DataLoader struct {
	config *Config
}

// NewDataLoader creates a new DataLoader instance
func NewDataLoader(config *Config) *DataLoader {
	return &DataLoader{config: config}
}

// LoadUsageData loads and parses all JSONL files from Claude data directory
func (dl *DataLoader) LoadUsageData() ([]UsageData, error) {
	var allData []UsageData

	projectsPath := filepath.Join(dl.config.ClaudePath, "projects")
	
	err := filepath.WalkDir(projectsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip directories we can't access
			return nil
		}

		if !strings.HasSuffix(path, ".jsonl") {
			return nil
		}

		data, err := dl.parseJSONLFile(path)
		if err != nil {
			// Log error but continue processing other files
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", path, err)
			return nil
		}

		allData = append(allData, data...)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return dl.filterByDateRange(allData), nil
}

// parseJSONLFile parses a single JSONL file
func (dl *DataLoader) parseJSONLFile(filePath string) ([]UsageData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var data []UsageData
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var rawData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &rawData); err != nil {
			// Skip malformed lines
			continue
		}

		usageData, err := dl.parseUsageEntry(rawData)
		if err != nil {
			// Skip entries we can't parse
			continue
		}

		data = append(data, usageData)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return data, nil
}

// parseUsageEntry converts raw JSON data to UsageData struct
func (dl *DataLoader) parseUsageEntry(rawData map[string]interface{}) (UsageData, error) {
	var usage UsageData

	// Parse timestamp
	if ts, ok := rawData["timestamp"].(string); ok {
		t, err := time.Parse(time.RFC3339, ts)
		if err != nil {
			return usage, fmt.Errorf("invalid timestamp: %w", err)
		}
		usage.Timestamp = t
	} else {
		return usage, fmt.Errorf("missing timestamp")
	}

	// Parse token counts
	usage.InputTokens = dl.getIntField(rawData, "input_tokens")
	usage.OutputTokens = dl.getIntField(rawData, "output_tokens")
	usage.CacheCreationTokens = dl.getIntField(rawData, "cache_creation_tokens")
	usage.CacheReadTokens = dl.getIntField(rawData, "cache_read_tokens")

	// Parse cost
	if cost, ok := rawData["cost"].(float64); ok {
		usage.Cost = cost
	}

	return usage, nil
}

// getIntField safely extracts integer field from raw data
func (dl *DataLoader) getIntField(data map[string]interface{}, field string) int {
	if val, ok := data[field]; ok {
		switch v := val.(type) {
		case float64:
			return int(v)
		case int:
			return v
		}
	}
	return 0
}

// filterByDateRange filters usage data by date range if specified
func (dl *DataLoader) filterByDateRange(data []UsageData) []UsageData {
	if dl.config.Since == nil && dl.config.Until == nil {
		return data
	}

	var filtered []UsageData
	for _, entry := range data {
		if dl.config.Since != nil && entry.Timestamp.Before(*dl.config.Since) {
			continue
		}
		if dl.config.Until != nil && entry.Timestamp.After(*dl.config.Until) {
			continue
		}
		filtered = append(filtered, entry)
	}

	return filtered
}

// GetSessionFromPath extracts project and session names from file path
func GetSessionFromPath(filePath, basePath string) (project, session string) {
	rel, err := filepath.Rel(basePath, filePath)
	if err != nil {
		return "", ""
	}

	parts := strings.Split(rel, string(filepath.Separator))
	if len(parts) >= 3 && parts[0] == "projects" {
		return parts[1], parts[2]
	}

	return "", ""
}
