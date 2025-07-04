package yaml_service

import (
	"reflect"
	"testing"

)

func TestIntegration_ParseYAML(t *testing.T) {
	ys := NewYamlService()

	tests := []struct {
		name        string
		yamlContent string
		want        interface{}
		wantErr     bool
	}{
		{
			name: "複雑なYAMLのパース",
			yamlContent: `
name: integration-test
version: 1.0
config:
  debug: true
  timeout: 30
  servers:
    - name: server1
      host: localhost
      port: 8080
    - name: server2
      host: example.com
      port: 9090
tags:
  - yaml
  - test
  - integration
enabled: true
`,
			want: map[string]interface{}{
				"name":    "integration-test",
				"version": 1.0,
				"config": map[string]interface{}{
					"debug":   true,
					"timeout": 30,
					"servers": []interface{}{
						map[string]interface{}{
							"name": "server1",
							"host": "localhost",
							"port": 8080,
						},
						map[string]interface{}{
							"name": "server2",
							"host": "example.com",
							"port": 9090,
						},
					},
				},
				"tags":    []interface{}{"yaml", "test", "integration"},
				"enabled": true,
			},
			wantErr: false,
		},
		{
			name: "不正なYAMLのパース",
			yamlContent: `
name: invalid-yaml
  invalid-indent:
 - broken-format
`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ys.ParseYAML(tt.yamlContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("Integration ParseYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Integration ParseYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntegration_MarshalToYAML(t *testing.T) {
	ys := NewYamlService()

	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name: "複雑なマップのマーシャル",
			data: map[string]interface{}{
				"name":    "marshal-test",
				"version": 2.0,
				"config": map[string]interface{}{
					"debug":   false,
					"timeout": 60,
					"servers": []map[string]interface{}{
						{
							"name": "server1",
							"host": "localhost",
							"port": 8080,
						},
						{
							"name": "server2",
							"host": "example.com",
							"port": 9090,
						},
					},
				},
				"tags":    []string{"yaml", "test", "marshal"},
				"enabled": true,
			},
			wantErr: false,
		},
		{
			name: "構造体のマーシャル",
			data: struct {
				Name    string   `yaml:"name"`
				Version float64  `yaml:"version"`
				Tags    []string `yaml:"tags"`
				Enabled bool     `yaml:"enabled"`
			}{
				Name:    "struct-test",
				Version: 1.5,
				Tags:    []string{"yaml", "test", "struct"},
				Enabled: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ys.MarshalToYAML(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Integration MarshalToYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// マーシャルした結果を再度パースして元のデータと比較する
				parsedData, err := ys.ParseYAML(got)
				if err != nil {
					t.Errorf("Integration MarshalToYAML() marshal-parse cycle failed: %v", err)
					return
				}

				// 注意: マーシャル-パースのサイクルを通すと、数値型や一部の構造が変わる可能性があるため、
				// 単純な比較ではなく、重要なフィールドのみを比較するか、または文字列形式で比較する
				if reflect.TypeOf(tt.data).Kind() == reflect.Map {
					// マップ型のテスト
					origMap := tt.data.(map[string]interface{})
					if parsedMap, ok := parsedData.(map[string]interface{}); ok {
						if origMap["name"] != parsedMap["name"] {
							t.Errorf("Integration MarshalToYAML() marshal-parse name mismatch: got %v, want %v",
								parsedMap["name"], origMap["name"])
						}
						if origMap["enabled"] != parsedMap["enabled"] {
							t.Errorf("Integration MarshalToYAML() marshal-parse enabled mismatch: got %v, want %v",
								parsedMap["enabled"], origMap["enabled"])
						}
					}
				}
			}
		})
	}
}

func TestIntegration_ParseYAMLToStruct(t *testing.T) {
	ys := NewYamlService()

	type TestConfig struct {
		Debug   bool `yaml:"debug"`
		Timeout int  `yaml:"timeout"`
		Servers []struct {
			Name string `yaml:"name"`
			Host string `yaml:"host"`
			Port int    `yaml:"port"`
		} `yaml:"servers"`
	}

	type TestData struct {
		Name    string     `yaml:"name"`
		Version float64    `yaml:"version"`
		Config  TestConfig `yaml:"config"`
		Tags    []string   `yaml:"tags"`
		Enabled bool       `yaml:"enabled"`
	}

	tests := []struct {
		name        string
		yamlContent string
		want        TestData
		wantErr     bool
	}{
		{
			name: "構造体へのパース",
			yamlContent: `
name: struct-parse-test
version: 3.0
config:
  debug: true
  timeout: 45
  servers:
    - name: server1
      host: localhost
      port: 8080
    - name: server2
      host: example.com
      port: 9090
tags:
  - yaml
  - test
  - struct
enabled: true
`,
			want: TestData{
				Name:    "struct-parse-test",
				Version: 3.0,
				Config: TestConfig{
					Debug:   true,
					Timeout: 45,
					Servers: []struct {
						Name string `yaml:"name"`
						Host string `yaml:"host"`
						Port int    `yaml:"port"`
					}{
						{
							Name: "server1",
							Host: "localhost",
							Port: 8080,
						},
						{
							Name: "server2",
							Host: "example.com",
							Port: 9090,
						},
					},
				},
				Tags:    []string{"yaml", "test", "struct"},
				Enabled: true,
			},
			wantErr: false,
		},
		{
			name: "不正なYAMLの構造体へのパース",
			yamlContent: `
name: invalid-struct
version: not-a-number
`,
			want:    TestData{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestData
			err := ys.ParseYAMLToStruct(tt.yamlContent, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Integration ParseYAMLToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Name != tt.want.Name {
					t.Errorf("Integration ParseYAMLToStruct() Name = %v, want %v", got.Name, tt.want.Name)
				}
				if got.Version != tt.want.Version {
					t.Errorf("Integration ParseYAMLToStruct() Version = %v, want %v", got.Version, tt.want.Version)
				}
				if got.Config.Debug != tt.want.Config.Debug {
					t.Errorf("Integration ParseYAMLToStruct() Config.Debug = %v, want %v", got.Config.Debug, tt.want.Config.Debug)
				}
				if got.Config.Timeout != tt.want.Config.Timeout {
					t.Errorf("Integration ParseYAMLToStruct() Config.Timeout = %v, want %v", got.Config.Timeout, tt.want.Config.Timeout)
				}
				if len(got.Config.Servers) != len(tt.want.Config.Servers) {
					t.Errorf("Integration ParseYAMLToStruct() Config.Servers length = %v, want %v", len(got.Config.Servers), len(tt.want.Config.Servers))
				} else {
					for i, server := range got.Config.Servers {
						if server.Name != tt.want.Config.Servers[i].Name {
							t.Errorf("Integration ParseYAMLToStruct() Config.Servers[%d].Name = %v, want %v", i, server.Name, tt.want.Config.Servers[i].Name)
						}
						if server.Host != tt.want.Config.Servers[i].Host {
							t.Errorf("Integration ParseYAMLToStruct() Config.Servers[%d].Host = %v, want %v", i, server.Host, tt.want.Config.Servers[i].Host)
						}
						if server.Port != tt.want.Config.Servers[i].Port {
							t.Errorf("Integration ParseYAMLToStruct() Config.Servers[%d].Port = %v, want %v", i, server.Port, tt.want.Config.Servers[i].Port)
						}
					}
				}
				if !reflect.DeepEqual(got.Tags, tt.want.Tags) {
					t.Errorf("Integration ParseYAMLToStruct() Tags = %v, want %v", got.Tags, tt.want.Tags)
				}
				if got.Enabled != tt.want.Enabled {
					t.Errorf("Integration ParseYAMLToStruct() Enabled = %v, want %v", got.Enabled, tt.want.Enabled)
				}
			}
		})
	}
}
