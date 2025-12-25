package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type FlexibleString string

func (fs *FlexibleString) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*fs = ""
		return nil
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return err
		}
		*fs = FlexibleString(s)
		return nil
	}
	var num json.Number
	if err := json.Unmarshal(trimmed, &num); err == nil {
		*fs = FlexibleString(num.String())
		return nil
	}
	return fmt.Errorf("値を文字列に変換できません: %s", string(trimmed))
}

func (fs FlexibleString) String() string {
	return string(fs)
}

type PersonaExtractionRequest struct {
	Context    string    `json:"context,omitempty"`
	Characters []Persona `json:"characters"`
}

type Persona struct {
	Name          string         `json:"name"`
	Age           FlexibleString `json:"age"`
	Gender        string         `json:"gender"`
	Species       string         `json:"species"`
	Job           string         `json:"job"`
	Hobbies       []string       `json:"hobbies"`
	Personality   string         `json:"personality"`
	Motivations   []string       `json:"motivations,omitempty"`
	Goals         []string       `json:"goals,omitempty"`
	Relationships []Relationship `json:"relationships,omitempty"`
	Evidence      []Evidence     `json:"evidence,omitempty"`
	Notes         string         `json:"notes,omitempty"`
}

type Relationship struct {
	With  string `json:"with"`
	Type  string `json:"type,omitempty"`
	Notes string `json:"notes,omitempty"`
}

type Evidence struct {
	Quote  string `json:"quote"`
	Source string `json:"source,omitempty"`
}

func NormalizeExtractionRequest(req *PersonaExtractionRequest) {
	if req == nil {
		return
	}
	req.Context = strings.TrimSpace(req.Context)
	for i := range req.Characters {
		normalizePersona(&req.Characters[i])
	}
}

func normalizePersona(p *Persona) {
	p.Name = strings.TrimSpace(p.Name)
	p.Gender = strings.TrimSpace(p.Gender)
	p.Species = strings.TrimSpace(p.Species)
	p.Job = strings.TrimSpace(p.Job)
	p.Personality = strings.TrimSpace(p.Personality)
	p.Notes = strings.TrimSpace(p.Notes)
	p.Age = FlexibleString(strings.TrimSpace(p.Age.String()))
	p.Hobbies = trimList(p.Hobbies)
	p.Motivations = trimList(p.Motivations)
	p.Goals = trimList(p.Goals)
	p.Relationships = normalizeRelationships(p.Relationships)
	p.Evidence = normalizeEvidence(p.Evidence)
}

func trimList(values []string) []string {
	if len(values) == 0 {
		return values
	}
	cleaned := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func normalizeRelationships(values []Relationship) []Relationship {
	if len(values) == 0 {
		return values
	}
	cleaned := make([]Relationship, 0, len(values))
	for _, rel := range values {
		rel.With = strings.TrimSpace(rel.With)
		rel.Type = strings.TrimSpace(rel.Type)
		rel.Notes = strings.TrimSpace(rel.Notes)
		if rel.With == "" && rel.Type == "" && rel.Notes == "" {
			continue
		}
		cleaned = append(cleaned, rel)
	}
	return cleaned
}

func normalizeEvidence(values []Evidence) []Evidence {
	if len(values) == 0 {
		return values
	}
	cleaned := make([]Evidence, 0, len(values))
	for _, ev := range values {
		ev.Quote = strings.TrimSpace(ev.Quote)
		ev.Source = strings.TrimSpace(ev.Source)
		if ev.Quote == "" && ev.Source == "" {
			continue
		}
		cleaned = append(cleaned, ev)
	}
	return cleaned
}

func ValidateExtractionRequest(req PersonaExtractionRequest) error {
	if req.Characters == nil {
		return fmt.Errorf("charactersフィールドは必須です")
	}
	if len(req.Characters) == 0 {
		return fmt.Errorf("charactersは最低1件必要です")
	}
	for idx, persona := range req.Characters {
		if persona.Name == "" {
			return fmt.Errorf("characters[%d].nameは空にできません", idx)
		}
		if persona.Age.String() == "" {
			return fmt.Errorf("characters[%d].ageは空にできません", idx)
		}
		if persona.Gender == "" {
			return fmt.Errorf("characters[%d].genderは空にできません", idx)
		}
		if persona.Species == "" {
			return fmt.Errorf("characters[%d].speciesは空にできません", idx)
		}
		if persona.Job == "" {
			return fmt.Errorf("characters[%d].jobは空にできません", idx)
		}
		if len(persona.Hobbies) == 0 {
			return fmt.Errorf("characters[%d].hobbiesは1件以上必要です", idx)
		}
		for hobbyIdx, hobby := range persona.Hobbies {
			if hobby == "" {
				return fmt.Errorf("characters[%d].hobbies[%d]は空にできません", idx, hobbyIdx)
			}
		}
		if persona.Personality == "" {
			return fmt.Errorf("characters[%d].personalityは空にできません", idx)
		}
		for relIdx, rel := range persona.Relationships {
			if rel.With == "" {
				return fmt.Errorf("characters[%d].relationships[%d].withは空にできません", idx, relIdx)
			}
		}
		for evIdx, ev := range persona.Evidence {
			if ev.Quote == "" {
				return fmt.Errorf("characters[%d].evidence[%d].quoteは空にできません", idx, evIdx)
			}
		}
	}
	return nil
}

func PersonaJSONSchema() map[string]any {
	relationshipSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"with": map[string]any{
				"type":        "string",
				"description": "関係する人物名",
			},
			"type": map[string]any{
				"type":        "string",
				"description": "関係性（例: sibling, mentor）",
			},
			"notes": map[string]any{
				"type":        "string",
				"description": "関係についての補足",
			},
		},
		"required":             []string{"with"},
		"additionalProperties": false,
	}

	evidenceSchema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"quote": map[string]any{
				"type":        "string",
				"description": "キャラクターに関する引用や根拠",
			},
			"source": map[string]any{
				"type":        "string",
				"description": "引用元（チャプター名や資料など）",
			},
		},
		"required":             []string{"quote"},
		"additionalProperties": false,
	}

	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": map[string]any{
				"type":        "string",
				"description": "キャラクターの名前",
			},
			"age": map[string]any{
				"description": "年齢。文字列または数値のどちらでも指定可能",
				"oneOf": []map[string]any{
					{"type": "string"},
					{"type": "number"},
				},
			},
			"gender": map[string]any{
				"type":        "string",
				"description": "性別やジェンダー表現",
			},
			"species": map[string]any{
				"type":        "string",
				"description": "種族や所属（ヒト、AIなど）",
			},
			"spieces": map[string]any{
				"type":        "string",
				"description": "speciesのスペルミス互換フィールド。指定された場合は自動的にspeciesへ統合されます。",
			},
			"job": map[string]any{
				"type":        "string",
				"description": "職業や役割",
			},
			"hobbies": map[string]any{
				"type":     "array",
				"minItems": 1,
				"items": map[string]any{
					"type": "string",
				},
				"description": "趣味。1件以上の文字列配列",
			},
			"personality": map[string]any{
				"type":        "string",
				"description": "性格や人格の特徴",
			},
			"motivations": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "動機や価値観",
			},
			"goals": map[string]any{
				"type":        "array",
				"items":       map[string]any{"type": "string"},
				"description": "短期・長期の目標",
			},
			"relationships": map[string]any{
				"type":        "array",
				"items":       relationshipSchema,
				"description": "他キャラクターとの関係",
			},
			"evidence": map[string]any{
				"type":        "array",
				"items":       evidenceSchema,
				"description": "設定の根拠となる引用や出来事",
			},
			"notes": map[string]any{
				"type":        "string",
				"description": "補足メモ",
			},
		},
		"required":             []string{"name", "age", "gender", "species", "job", "hobbies", "personality"},
		"additionalProperties": false,
	}
}
