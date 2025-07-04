package usecase_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/landmaster135/devbox/internal/yaml_parser/domain/entity"
	"github.com/landmaster135/devbox/internal/yaml_parser/usecase"
)

// モックYAMLリポジトリ
type MockYAMLRepository struct {
	parseFn            func(yamlContent string) (*entity.YAMLData, error)
	marshalFn          func(data interface{}) (string, error)
	parseToStructFn    func(yamlContent string, out interface{}) error
	readFileFn         func(filePath string) (string, error)
	parseFileFn        func(filePath string) (*entity.YAMLData, error)
	parseFileToStructFn func(filePath string, out interface{}) error
}

func (m *MockYAMLRepository) Parse(yamlContent string) (*entity.YAMLData, error) {
	return m.parseFn(yamlContent)
}

func (m *MockYAMLRepository) Marshal(data interface{}) (string, error) {
	return m.marshalFn(data)
}

func (m *MockYAMLRepository) ParseToStruct(yamlContent string, out interface{}) error {
	return m.parseToStructFn(yamlContent, out)
}

func (m *MockYAMLRepository) ReadFile(filePath string) (string, error) {
	if m.readFileFn != nil {
		return m.readFileFn(filePath)
	}
	return "", nil
}

func (m *MockYAMLRepository) ParseFile(filePath string) (*entity.YAMLData, error) {
	if m.parseFileFn != nil {
		return m.parseFileFn(filePath)
	}
	return nil, nil
}

func (m *MockYAMLRepository) ParseFileToStruct(filePath string, out interface{}) error {
	if m.parseFileToStructFn != nil {
		return m.parseFileToStructFn(filePath, out)
	}
	return nil
}

func TestYAMLUseCase_ParseYAML(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		parseFn     func(yamlContent string) (*entity.YAMLData, error)
		want        *entity.YAMLData
		wantErr     bool
	}{
		{
			name:        "正常なパース",
			yamlContent: "name: test",
			parseFn: func(yamlContent string) (*entity.YAMLData, error) {
				return entity.NewYAMLData(map[string]interface{}{"name": "test"}), nil
			},
			want:    entity.NewYAMLData(map[string]interface{}{"name": "test"}),
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
			mockRepo := &MockYAMLRepository{
				parseFn: tt.parseFn,
			}

			uc := usecase.NewYAMLUseCase(mockRepo)

			got, err := uc.ParseYAML(tt.yamlContent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got == nil || tt.want == nil {
					if got != tt.want {
						t.Errorf("ParseYAML() = %v, want %v", got, tt.want)
					}
				} else if !reflect.DeepEqual(got.GetData(), tt.want.GetData()) {
					t.Errorf("ParseYAML().GetData() = %v, want %v", got.GetData(), tt.want.GetData())
				}
			}
		})
	}
}

func TestYAMLUseCase_MarshalToYAML(t *testing.T) {
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
			mockRepo := &MockYAMLRepository{
				marshalFn: tt.marshalFn,
			}

			uc := usecase.NewYAMLUseCase(mockRepo)

			got, err := uc.MarshalToYAML(tt.data)
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

func TestYAMLUseCase_ParseYAMLToStruct(t *testing.T) {
	type TestStruct struct {
		Name string `yaml:"name"`
	}

	tests := []struct {
		name        string
		yamlContent string
		out         interface{}
		parseToStructFn func(yamlContent string, out interface{}) error
		wantErr     bool
	}{
		{
			name:        "正常なパース",
			yamlContent: "name: test",
			out:         &TestStruct{},
			parseToStructFn: func(yamlContent string, out interface{}) error {
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
			parseToStructFn: func(yamlContent string, out interface{}) error {
				return errors.New("parse error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockYAMLRepository{
				parseToStructFn: tt.parseToStructFn,
			}

			uc := usecase.NewYAMLUseCase(mockRepo)

			err := uc.ParseYAMLToStruct(tt.yamlContent, tt.out)
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
