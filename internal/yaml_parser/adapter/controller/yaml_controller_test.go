package controller_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/landmaster135/devbox/internal/yaml_parser/adapter/controller"
	"github.com/landmaster135/devbox/internal/yaml_parser/domain/entity"
	"github.com/landmaster135/devbox/internal/yaml_parser/usecase"
)

// YAMLUseCaseInterface はYAMLUseCaseのインターフェースを定義します
type YAMLUseCaseInterface interface {
	ParseYAML(yamlContent string) (*entity.YAMLData, error)
	MarshalToYAML(data interface{}) (string, error)
	ParseYAMLToStruct(yamlContent string, out interface{}) error
	ReadYAMLFile(filePath string) (string, error)
	ParseYAMLFile(filePath string) (*entity.YAMLData, error)
	ParseYAMLFileToStruct(filePath string, out interface{}) error
}

// モックYAMLユースケース
type MockYAMLUseCase struct {
	parseYAMLFn             func(yamlContent string) (*entity.YAMLData, error)
	marshalToYAMLFn         func(data interface{}) (string, error)
	parseYAMLToStructFn     func(yamlContent string, out interface{}) error
	readYAMLFileFn          func(filePath string) (string, error)
	parseYAMLFileFn         func(filePath string) (*entity.YAMLData, error)
	parseYAMLFileToStructFn func(filePath string, out interface{}) error
}

func (m *MockYAMLUseCase) ParseYAML(yamlContent string) (*entity.YAMLData, error) {
	return m.parseYAMLFn(yamlContent)
}

func (m *MockYAMLUseCase) MarshalToYAML(data interface{}) (string, error) {
	return m.marshalToYAMLFn(data)
}

func (m *MockYAMLUseCase) ParseYAMLToStruct(yamlContent string, out interface{}) error {
	return m.parseYAMLToStructFn(yamlContent, out)
}

func (m *MockYAMLUseCase) ReadYAMLFile(filePath string) (string, error) {
	if m.readYAMLFileFn != nil {
		return m.readYAMLFileFn(filePath)
	}
	return "", nil
}

func (m *MockYAMLUseCase) ParseYAMLFile(filePath string) (*entity.YAMLData, error) {
	if m.parseYAMLFileFn != nil {
		return m.parseYAMLFileFn(filePath)
	}
	return nil, nil
}

func (m *MockYAMLUseCase) ParseYAMLFileToStruct(filePath string, out interface{}) error {
	if m.parseYAMLFileToStructFn != nil {
		return m.parseYAMLFileToStructFn(filePath, out)
	}
	return nil
}

// テスト用のコントローラー作成関数
func createYAMLController(yamlUseCase YAMLUseCaseInterface) *controller.YAMLController {
	// モックを使用してコントローラーを作成
	return &controller.YAMLController{
		YAMLUseCase: yamlUseCase,
	}
}

func TestNewYAMLController(t *testing.T) {
	// YAMLUseCaseの作成
	yamlUseCase := &usecase.YAMLUseCase{}

	// コントローラーの作成
	c := controller.NewYAMLController(yamlUseCase)

	// アサーション
	if c == nil {
		t.Error("NewYAMLController() returned nil")
	}

	// YAMLUseCaseが正しく設定されているか確認
	if c.YAMLUseCase != yamlUseCase {
		t.Errorf("NewYAMLController() YAMLUseCase = %v, want %v", c.YAMLUseCase, yamlUseCase)
	}
}

func TestYAMLController_ParseYAML(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		parseFn     func(yamlContent string) (*entity.YAMLData, error)
		want        interface{}
		wantErr     bool
	}{
		{
			name:        "正常なパース",
			yamlContent: "name: test",
			parseFn: func(yamlContent string) (*entity.YAMLData, error) {
				return entity.NewYAMLData(map[string]interface{}{"name": "test"}), nil
			},
			want:    map[string]interface{}{"name": "test"},
			wantErr: false,
		},
		{
			name:        "パースエラー",
			yamlContent: "invalid yaml",
			parseFn: func(yamlContent string) (*entity.YAMLData, error) {
				return nil, errors.New("parse error")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockYAMLUseCase{
				parseYAMLFn: tt.parseFn,
			}

			c := createYAMLController(mockUseCase)

			got, err := c.ParseYAML(tt.yamlContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYAMLController_MarshalToYAML(t *testing.T) {
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
			mockUseCase := &MockYAMLUseCase{
				marshalToYAMLFn: tt.marshalFn,
			}

			c := createYAMLController(mockUseCase)

			got, err := c.MarshalToYAML(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalToYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MarshalToYAML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestYAMLController_ParseYAMLToStruct(t *testing.T) {
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
			mockUseCase := &MockYAMLUseCase{
				parseYAMLToStructFn: tt.parseStructFn,
			}

			c := createYAMLController(mockUseCase)

			err := c.ParseYAMLToStruct(tt.yamlContent, tt.out)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseYAMLToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				ts := tt.out.(*TestStruct)
				if ts.Name != "test" {
					t.Errorf("ParseYAMLToStruct() did not set the struct correctly")
				}
			}
		})
	}
}
