package getconstants

import (
	"fmt"
	"math"
	"strings"
)

func GetConstants() map[string]float64 {
	return map[string]float64{
		"pi":  math.Pi,
		"e":   math.E,
		"tau": 2 * math.Pi,
	}
}

func HandleToGetConstants() (string, error) {
	constants := GetConstants()

	var result strings.Builder
	result.WriteString("利用可能な数学定数:\n")
	for name, value := range constants {
		result.WriteString(fmt.Sprintf("%s = %f\n", name, value))
	}

	return result.String(), nil
}

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Execute() (string, error) {
	return HandleToGetConstants()
}
