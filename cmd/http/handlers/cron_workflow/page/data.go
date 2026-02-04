package page

import (
	"fmt"
	"net/http"
	"strings"

	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
)

type WorkflowDefinition struct {
	Key         string
	DisplayName string
	Summary     string
	TagDisplay  string
}

type workflowCategoryDefinition struct {
	ID          string
	Title       string
	Description string
	Workflows   []WorkflowDefinition
}

type workflowCardView struct {
	Key                 string
	Name                string
	Summary             string
	CronDefinition      string
	ManualAction        string
	ManualMethod        string
	ManualWorkflowField string
	ManualWorkflowValue string
	ProcessName         string
}

type workflowCategoryView struct {
	ID          string
	Title       string
	Description string
	Workflows   []workflowCardView
}

type workflowPageData struct {
	Title       string
	Description string
	Categories  []workflowCategoryView
}

var WorkflowCategoryDefinitions = []workflowCategoryDefinition{
	{
		ID:          "daily-workflow",
		Title:       "Daily Workflow",
		Description: "Daily workflows that prepare Discord-ready summaries.",
		Workflows: []WorkflowDefinition{
			{
				Key:         "daily-tokyo-weather",
				DisplayName: "Daily Tokyo weather notification",
				Summary:     "Fetches a three-day forecast for Tokyo and shares it with the weather Discord channel.",
				TagDisplay:  "weather",
			},
			{
				Key:         "daily-heading",
				DisplayName: "Daily heading Discord notification",
				Summary:     "Generates the day's heading template and posts it to the daily heading Discord webhook.",
				TagDisplay:  "heading",
			},
		},
	},
	{
		ID:          "hourly-workflow",
		Title:       "Hourly Workflow",
		Description: "Hourly workflows that capture telemetry on a short cadence.",
		Workflows: []WorkflowDefinition{
			{
				Key:         "ubuntu-pc-info",
				DisplayName: "Ubuntu PC info snapshot",
				Summary:     "Collects Ubuntu host CPU, memory, and temperature stats and archives the snapshot to disk.",
				TagDisplay:  "pc-info",
			},
		},
	},
}

func buildWorkflowPageData(workflows []usecases.Workflow, manualRunEndpoint string, manualWorkflowFieldName string) (workflowPageData, error) {
	workflowByDescription := make(map[string]usecases.Workflow, len(workflows))
	for _, wf := range workflows {
		desc := strings.TrimSpace(wf.Description)
		if desc == "" {
			continue
		}
		workflowByDescription[desc] = wf
	}

	data := workflowPageData{
		Title:       "Cron Workflow Orchestrator",
		Description: "Workflows registered in workflow.List() are grouped by category for quick inspection.",
	}

	for _, catDef := range WorkflowCategoryDefinitions {
		catView := workflowCategoryView{
			ID:          catDef.ID,
			Title:       catDef.Title,
			Description: catDef.Description,
		}
		for _, wfDef := range catDef.Workflows {
			wf, ok := workflowByDescription[strings.TrimSpace(wfDef.DisplayName)]
			if !ok {
				continue
			}

			cronDef, _, err := wf.GetCronDefinition()
			if err != nil {
				return workflowPageData{}, fmt.Errorf("get cron definition for %s: %w", wfDef.DisplayName, err)
			}

			catView.Workflows = append(catView.Workflows, workflowCardView{
				Key:                 wfDef.Key,
				Name:                wfDef.DisplayName,
				Summary:             wfDef.Summary,
				CronDefinition:      cronDef,
				ManualAction:        manualRunEndpoint,
				ManualMethod:        http.MethodPost,
				ManualWorkflowField: manualWorkflowFieldName,
				ManualWorkflowValue: wfDef.Key,
				ProcessName:         wfDef.TagDisplay,
			})
		}
		if len(catView.Workflows) == 0 {
			continue
		}
		data.Categories = append(data.Categories, catView)
	}

	return data, nil
}
