package gateway_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/landmaster135/devbox/internal/independencies/yaml_parser/adapter/gateway"
	"github.com/landmaster135/devbox/internal/independencies/yaml_parser/domain/entity"
)

// モックYAMLパーサー
type MockYAMLParser struct {
	parseYAMLFn             func(yamlContent string) (interface{}, error)
	marshalToYAMLFn         func(data interface{}) (string, error)
	parseYAMLToStructFn     func(yamlContent string, out interface{}) error
	readYAMLFileFn          func(filePath string) (string, error)
	parseYAMLFileFn         func(filePath string) (interface{}, error)
	parseYAMLFileToStructFn func(filePath string, out interface{}) error
}

func (m *MockYAMLParser) ParseYAML(yamlContent string) (interface{}, error) {
	return m.parseYAMLFn(yamlContent)
}

func (m *MockYAMLParser) MarshalToYAML(data interface{}) (string, error) {
	return m.marshalToYAMLFn(data)
}

func (m *MockYAMLParser) ParseYAMLToStruct(yamlContent string, out interface{}) error {
	return m.parseYAMLToStructFn(yamlContent, out)
}

func (m *MockYAMLParser) ReadYAMLFile(filePath string) (string, error) {
	if m.readYAMLFileFn != nil {
		return m.readYAMLFileFn(filePath)
	}
	return "", nil
}

func (m *MockYAMLParser) ParseYAMLFile(filePath string) (interface{}, error) {
	if m.parseYAMLFileFn != nil {
		return m.parseYAMLFileFn(filePath)
	}
	return nil, nil
}

func (m *MockYAMLParser) ParseYAMLFileToStruct(filePath string, out interface{}) error {
	if m.parseYAMLFileToStructFn != nil {
		return m.parseYAMLFileToStructFn(filePath, out)
	}
	return nil
}

func TestNewYAMLGateway(t *testing.T) {
	// モックパーサーの作成
	mockParser := &MockYAMLParser{}

	// ゲートウェイの作成
	g := gateway.NewYAMLGateway(mockParser)

	// アサーション
	if g == nil {
		t.Error("NewYAMLGateway() returned nil")
	}
}

func TestYAMLGateway_Parse(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		parseFn     func(yamlContent string) (interface{}, error)
		want        *entity.YAMLData
		wantErr     bool
	}{
		{
			name:        "正常なパース",
			yamlContent: "name: test",
			parseFn: func(yamlContent string) (interface{}, error) {
				return map[string]interface{}{"name": "test"}, nil
			},
			want:    entity.NewYAMLData(map[string]interface{}{"name": "test"}),
			wantErr: false,
		},
		{
			name:        "パースエラー",
			yamlContent: "invalid yaml",
			parseFn: func(yamlContent string) (interface{}, error) {
				return nil, errors.New("parse error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := &MockYAMLParser{
				parseYAMLFn: tt.parseFn,
			}

			g := gateway.NewYAMLGateway(mockParser)

			got, err := g.Parse(tt.yamlContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil || tt.want == nil {
					if got != tt.want {
						t.Errorf("Parse() = %v, want %v", got, tt.want)
					}
				} else if !reflect.DeepEqual(got.GetData(), tt.want.GetData()) {
					t.Errorf("Parse().GetData() = %v, want %v", got.GetData(), tt.want.GetData())
				}
			}
		})
	}
}

func TestYAMLGateway_Marshal(t *testing.T) {
	tests := []struct {
		name      string
		data      interface{}
		marshalFn func(data interface{}) (string, error)
		want      string
		wantErr   bool
	}{
		{
			name: "正常なマーシャル",
			data: map[string]interface{}{"name": "test"},
			marshalFn: func(data interface{}) (string, error) {
				return "name: test\n", nil
			},
			want:    "name: test\n",
			wantErr: false,
		},
		{
			name: "マーシャルエラー",
			data: make(chan int), // マーシャルできない型
			marshalFn: func(data interface{}) (string, error) {
				return "", errors.New("marshal error")
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := &MockYAMLParser{
				marshalToYAMLFn: tt.marshalFn,
			}

			g := gateway.NewYAMLGateway(mockParser)

			got, err := g.Marshal(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYAMLGateway_ParseToStruct(t *testing.T) {
	type TestStruct struct {
		Name string `yaml:"name"`
	}

	tests := []struct {
		name          string
		yamlContent   string
		out           interface{}
		parseStructFn func(yamlContent string, out interface{}) error
		wantErr       bool
	}{
		{
			name:        "正常なパース",
			yamlContent: "name: test",
			out:         &TestStruct{},
			parseStructFn: func(yamlContent string, out interface{}) error {
				ts := out.(*TestStruct)
				ts.Name = "test"
				return nil
			},
			wantErr: false,
		},
		{
			name:        "パースエラー",
			yamlContent: "invalid yaml",
			out:         &TestStruct{},
			parseStructFn: func(yamlContent string, out interface{}) error {
				return errors.New("parse error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockParser := &MockYAMLParser{
				parseYAMLToStructFn: tt.parseStructFn,
			}

			g := gateway.NewYAMLGateway(mockParser)

			err := g.ParseToStruct(tt.yamlContent, tt.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				ts := tt.out.(*TestStruct)
				if ts.Name != "test" {
					t.Errorf("ParseToStruct() did not set the struct correctly")
				}
			}
		})
	}
}
