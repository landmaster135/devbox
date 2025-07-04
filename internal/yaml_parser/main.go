package main

import (
	"fmt"
	"log"
	"os"

	yaml_service "github.com/landmaster135/devbox/internal/yaml_parser/service"
)

func main() {
	ys := yaml_service.NewYamlService()

	// YAMLデータの例
	yamlContent := `
name: example
version: 1.0
details:
  description: This is an example
  tags:
    - yaml
    - parser
    - example
`

	// YAMLのパース
	parsedData, err := ys.ParseYAML(yamlContent)
	if err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	// パースされたデータの表示
	fmt.Printf("Parsed data: %+v\n", parsedData)

	// 特定の構造体へのパース例
	type Config struct {
		Name    string  `yaml:"name"`
		Version float64 `yaml:"version"`
		Details struct {
			Description string   `yaml:"description"`
			Tags        []string `yaml:"tags"`
		} `yaml:"details"`
	}

	var config Config
	err = ys.ParseYAMLToStruct(yamlContent, &config)
	if err != nil {
		log.Fatalf("Failed to parse YAML to struct: %v", err)
	}

	fmt.Printf("Parsed to struct: %+v\n", config)
	fmt.Printf("Name: %s\n", config.Name)
	fmt.Printf("Tags: %v\n", config.Details.Tags)

	// YAMLへの変換例
	newData := map[string]interface{}{
		"name":     "new example",
		"values":   []int{1, 2, 3},
		"isActive": true,
	}

	yamlString, err := ys.MarshalToYAML(newData)
	if err != nil {
		log.Fatalf("Failed to marshal to YAML: %v", err)
	}

	fmt.Printf("Marshaled YAML:\n%s\n", yamlString)

	// ファイルからYAMLを読み込む例
	// 一時ファイルの作成
	tempFile := "temp_config.yml"
	err = os.WriteFile(tempFile, []byte(yamlContent), 0644)
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile) // プログラム終了時に一時ファイルを削除

	// ファイルからYAMLを読み込む
	fileContent, err := ys.ReadYAMLFile(tempFile)
	if err != nil {
		log.Fatalf("Failed to read YAML file: %v", err)
	}
	fmt.Printf("File content:\n%s\n", fileContent)

	// ファイルからYAMLをパースする
	parsedFileData, err := ys.ParseYAMLFile(tempFile)
	if err != nil {
		log.Fatalf("Failed to parse YAML file: %v", err)
	}
	fmt.Printf("Parsed file data: %+v\n", parsedFileData)

	// ファイルからYAMLを構造体にパースする
	var fileConfig Config
	err = ys.ParseYAMLFileToStruct(tempFile, &fileConfig)
	if err != nil {
		log.Fatalf("Failed to parse YAML file to struct: %v", err)
	}
	fmt.Printf("Parsed file to struct: %+v\n", fileConfig)
	fmt.Printf("File Name: %s\n", fileConfig.Name)
	fmt.Printf("File Tags: %v\n", fileConfig.Details.Tags)

	fmt.Println("YAML parser completed successfully!")
}
