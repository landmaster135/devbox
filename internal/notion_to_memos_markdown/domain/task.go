package domain

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Task struct {
	ConID       string       `json:"con_id"`
	PageTitle   string       `json:"page_title"`
	StatusID    *string      `json:"status_id"`
	Priority    TaskPriority `json:"priority"`
	DoneAtStart string       `json:"done_at_start"`
	UpdatedAt   string       `json:"updated_at"`
	Tags        []TaskTag    `json:"tags"`
}

type TaskTag struct {
	PageTitle string `json:"page_title"`
}

type TaskPriority struct {
	IntValue *int
	Text     string
}

func (p *TaskPriority) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		p.IntValue = nil
		p.Text = ""
		return nil
	}

	var number float64
	if err := json.Unmarshal(data, &number); err == nil {
		v := int(number)
		p.IntValue = &v
		p.Text = strconv.Itoa(v)
		return nil
	}

	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		p.IntValue = nil
		p.Text = strings.TrimSpace(text)
		return nil
	}

	var object struct {
		PageTitle string   `json:"page_title"`
		Name      string   `json:"name"`
		Value     *float64 `json:"value"`
		Number    *float64 `json:"number"`
	}
	if err := json.Unmarshal(data, &object); err == nil {
		if object.Value != nil {
			v := int(*object.Value)
			p.IntValue = &v
			p.Text = strconv.Itoa(v)
			return nil
		}
		if object.Number != nil {
			v := int(*object.Number)
			p.IntValue = &v
			p.Text = strconv.Itoa(v)
			return nil
		}
		priorityText := strings.TrimSpace(object.PageTitle)
		if priorityText == "" {
			priorityText = strings.TrimSpace(object.Name)
		}
		p.IntValue = nil
		p.Text = priorityText
		return nil
	}

	return fmt.Errorf("priority の解析に失敗しました: %s", trimmed)
}
