package page

import (
	"errors"
	"fmt"
	"strings"

	page "github.com/landmaster135/devbox/cmd/http/handlers/cron_workflow/page"
	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
)

var (
	ErrWorkflowNotRegistered = errors.New("workflow not registered")
	ErrWorkflowUnavailable   = errors.New("workflow unavailable")
)

func initWorkflowDefinitions() map[string]page.WorkflowDefinition {
	defs := make(map[string]page.WorkflowDefinition)
	for _, cat := range page.WorkflowCategoryDefinitions {
		for _, def := range cat.Workflows {
			desc := strings.TrimSpace(def.Key)
			if desc == "" {
				continue
			}
			defs[desc] = def
		}
	}
	return defs
}

func WorkflowProcessByKey(workflows []usecases.Workflow, workflowKey string) (usecases.ProcessFunc, page.WorkflowDefinition, error) {
	workflowDefinitionsByKey := initWorkflowDefinitions()
	def, ok := workflowDefinitionsByKey[workflowKey]
	if !ok {
		return nil, page.WorkflowDefinition{}, fmt.Errorf("%w: %s", ErrWorkflowNotRegistered, workflowKey)
	}

	desc := strings.TrimSpace(def.DisplayName)
	for _, wf := range workflows {
		if strings.TrimSpace(wf.Description) != desc {
			continue
		}
		if wf.Process == nil {
			return nil, def, fmt.Errorf("%w: %s missing process", ErrWorkflowUnavailable, workflowKey)
		}
		return wf.Process, def, nil
	}

	return nil, def, fmt.Errorf("%w: %s not scheduled", ErrWorkflowUnavailable, workflowKey)
}
