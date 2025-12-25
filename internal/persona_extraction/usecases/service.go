package usecases

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/landmaster135/devbox/internal/persona_extraction/domain"
)

type PersonaService struct{}

func NewPersonaService() *PersonaService {
	return &PersonaService{}
}

func (s *PersonaService) HandleExtraction(args map[string]any) (string, error) {
	payload, err := decodePersonaArguments(args)
	if err != nil {
		return "", fmt.Errorf("引数の解析に失敗しました: %w", err)
	}
	domain.NormalizeExtractionRequest(&payload)
	summary, err := s.Extract(payload)
	if err != nil {
		return "", fmt.Errorf("ペルソナの処理に失敗しました: %w", err)
	}
	return summary, nil
}

func (s *PersonaService) Extract(req domain.PersonaExtractionRequest) (string, error) {
	if err := domain.ValidateExtractionRequest(req); err != nil {
		return "", err
	}
	return formatPersonaSummary(req), nil
}

func decodePersonaArguments(args map[string]any) (domain.PersonaExtractionRequest, error) {
	var payload domain.PersonaExtractionRequest
	if len(args) == 0 {
		return payload, fmt.Errorf("characters引数が指定されていません")
	}
	normalizeSpeciesAlias(args)
	raw, err := json.Marshal(args)
	if err != nil {
		return payload, fmt.Errorf("引数のシリアライズに失敗しました: %w", err)
	}
	if len(raw) == 0 || string(raw) == "null" {
		return payload, fmt.Errorf("characters引数が指定されていません")
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return payload, fmt.Errorf("引数のデシリアライズに失敗しました: %w", err)
	}
	return payload, nil
}

func normalizeSpeciesAlias(args map[string]any) {
	rawCharacters, ok := args["characters"]
	if !ok {
		return
	}
	switch list := rawCharacters.(type) {
	case []map[string]any:
		for _, entry := range list {
			ensureSpeciesField(entry)
		}
	case []any:
		for _, item := range list {
			if entry, ok := item.(map[string]any); ok {
				ensureSpeciesField(entry)
			}
		}
	default:
		// nothing to normalize
	}
}

func ensureSpeciesField(entry map[string]any) {
	if entry == nil {
		return
	}
	if _, ok := entry["species"]; ok {
		return
	}
	if alias, ok := entry["spieces"]; ok {
		entry["species"] = alias
	}
}

func formatPersonaSummary(req domain.PersonaExtractionRequest) string {
	var builder strings.Builder

	if req.Context != "" {
		builder.WriteString("Context:\n")
		builder.WriteString(req.Context)
		builder.WriteString("\n\n")
	}

	builder.WriteString("Personas:\n")
	for idx, persona := range req.Characters {
		builder.WriteString(fmt.Sprintf("%d. %s | Age: %s | Gender: %s | Species: %s | Job: %s\n",
			idx+1,
			persona.Name,
			persona.Age.String(),
			persona.Gender,
			persona.Species,
			persona.Job,
		))
		builder.WriteString(fmt.Sprintf("   Personality: %s\n", persona.Personality))
		builder.WriteString(fmt.Sprintf("   Hobbies: %s\n", strings.Join(persona.Hobbies, ", ")))
		if len(persona.Motivations) > 0 {
			builder.WriteString(fmt.Sprintf("   Motivations: %s\n", strings.Join(persona.Motivations, ", ")))
		}
		if len(persona.Goals) > 0 {
			builder.WriteString(fmt.Sprintf("   Goals: %s\n", strings.Join(persona.Goals, ", ")))
		}
		if len(persona.Relationships) > 0 {
			builder.WriteString("   Relationships:\n")
			for _, rel := range persona.Relationships {
				builder.WriteString("     - ")
				builder.WriteString(rel.With)
				if rel.Type != "" && rel.Notes != "" {
					builder.WriteString(fmt.Sprintf(" (%s): %s\n", rel.Type, rel.Notes))
				} else if rel.Type != "" {
					builder.WriteString(fmt.Sprintf(" (%s)\n", rel.Type))
				} else if rel.Notes != "" {
					builder.WriteString(fmt.Sprintf(": %s\n", rel.Notes))
				} else {
					builder.WriteString("\n")
				}
			}
		}
		if len(persona.Evidence) > 0 {
			builder.WriteString("   Evidence:\n")
			for _, ev := range persona.Evidence {
				if ev.Source != "" {
					builder.WriteString(fmt.Sprintf("     - \"%s\" (%s)\n", ev.Quote, ev.Source))
				} else {
					builder.WriteString(fmt.Sprintf("     - \"%s\"\n", ev.Quote))
				}
			}
		}
		if persona.Notes != "" {
			builder.WriteString(fmt.Sprintf("   Notes: %s\n", persona.Notes))
		}
		if idx < len(req.Characters)-1 {
			builder.WriteString("\n")
		}
	}

	return strings.TrimRight(builder.String(), "\n")
}
