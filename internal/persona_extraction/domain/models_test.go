package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFlexibleStringUnmarshalJSON(t *testing.T) {
	t.Parallel()

	var age FlexibleString
	if err := json.Unmarshal([]byte("27"), &age); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if age.String() != "27" {
		t.Fatalf("expected 27, got %q", age.String())
	}

	if err := json.Unmarshal([]byte("\" 35 \""), &age); err != nil {
		t.Fatalf("unexpected error for quoted string: %v", err)
	}
	if age.String() != " 35 " {
		t.Fatalf("expected raw string to remain, got %q", age.String())
	}

	if err := json.Unmarshal([]byte("null"), &age); err != nil {
		t.Fatalf("null should be allowed and produce empty string: %v", err)
	}
	if age.String() != "" {
		t.Fatalf("expected empty string for null, got %q", age.String())
	}

	if err := json.Unmarshal([]byte("{}"), &age); err == nil {
		t.Fatalf("expected error for invalid type")
	}
}

func TestNormalizeExtractionRequest(t *testing.T) {
	req := PersonaExtractionRequest{
		Context: "  mission ",
		Characters: []Persona{
			{
				Name:        "  Mika ",
				Age:         FlexibleString(" 27 "),
				Gender:      " female ",
				Species:     " human ",
				Job:         " reporter ",
				Hobbies:     []string{" photography ", " ", "urban exploration"},
				Personality: " curious ",
				Motivations: []string{" truth ", ""},
				Goals:       []string{" expose ", " "},
				Relationships: []Relationship{
					{With: " Ken ", Type: " brother ", Notes: " info "},
					{With: " ", Type: "", Notes: ""},
				},
				Evidence: []Evidence{
					{Quote: " I won't miss the truth ", Source: " chapter 1 "},
					{Quote: " ", Source: " "},
				},
				Notes: " cautious ",
			},
		},
	}

	NormalizeExtractionRequest(&req)

	if req.Context != "mission" {
		t.Fatalf("context not trimmed: %q", req.Context)
	}
	persona := req.Characters[0]
	if persona.Name != "Mika" || persona.Gender != "female" || persona.Species != "human" || persona.Job != "reporter" {
		t.Fatalf("basic fields not trimmed: %+v", persona)
	}
	if persona.Personality != "curious" || persona.Notes != "cautious" {
		t.Fatalf("personality/notes not trimmed: %+v", persona)
	}
	if persona.Age.String() != "27" {
		t.Fatalf("age not trimmed: %q", persona.Age.String())
	}
	expectedHobbies := []string{"photography", "urban exploration"}
	if strings.Join(persona.Hobbies, ";") != strings.Join(expectedHobbies, ";") {
		t.Fatalf("unexpected hobbies: %#v", persona.Hobbies)
	}
	expectedMotivations := []string{"truth"}
	if strings.Join(persona.Motivations, ";") != strings.Join(expectedMotivations, ";") {
		t.Fatalf("unexpected motivations: %#v", persona.Motivations)
	}
	if len(persona.Relationships) != 1 || persona.Relationships[0].With != "Ken" || persona.Relationships[0].Type != "brother" || persona.Relationships[0].Notes != "info" {
		t.Fatalf("relationships not normalized: %#v", persona.Relationships)
	}
	if len(persona.Evidence) != 1 || persona.Evidence[0].Quote != "I won't miss the truth" || persona.Evidence[0].Source != "chapter 1" {
		t.Fatalf("evidence not normalized: %#v", persona.Evidence)
	}
}

func TestValidateExtractionRequest(t *testing.T) {
	makeValid := func() PersonaExtractionRequest {
		return PersonaExtractionRequest{
			Characters: []Persona{
				{
					Name:        "Mika",
					Age:         FlexibleString("27"),
					Gender:      "female",
					Species:     "human",
					Job:         "reporter",
					Hobbies:     []string{"photography"},
					Personality: "curious",
				},
			},
		}
	}

	if err := ValidateExtractionRequest(makeValid()); err != nil {
		t.Fatalf("valid request should not error: %v", err)
	}

	cases := []struct {
		name   string
		mutate func(req *PersonaExtractionRequest)
		want   string
	}{
		{
			name: "nil characters",
			mutate: func(req *PersonaExtractionRequest) {
				req.Characters = nil
			},
			want: "charactersフィールドは必須です",
		},
		{
			name: "empty list",
			mutate: func(req *PersonaExtractionRequest) {
				req.Characters = []Persona{}
			},
			want: "charactersは最低1件必要です",
		},
		{
			name: "missing name",
			mutate: func(req *PersonaExtractionRequest) {
				req.Characters[0].Name = ""
			},
			want: "characters[0].name",
		},
		{
			name: "missing hobbies",
			mutate: func(req *PersonaExtractionRequest) {
				req.Characters[0].Hobbies = nil
			},
			want: "hobbiesは1件以上",
		},
		{
			name: "empty hobby",
			mutate: func(req *PersonaExtractionRequest) {
				req.Characters[0].Hobbies = []string{""}
			},
			want: "hobbies[0]",
		},
		{
			name: "missing personality",
			mutate: func(req *PersonaExtractionRequest) {
				req.Characters[0].Personality = ""
			},
			want: "personality",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			req := makeValid()
			tt.mutate(&req)
			if err := ValidateExtractionRequest(req); err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error containing %q, got %v", tt.want, err)
			}
		})
	}
}
