package cron_workflow

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	templ "github.com/a-h/templ"

	workflowPkg "github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow"
	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
	button "github.com/landmaster135/devbox/internal/templ_components/button"
	code "github.com/landmaster135/devbox/internal/templ_components/code"
	heading "github.com/landmaster135/devbox/internal/templ_components/heading"
	hiddenInput "github.com/landmaster135/devbox/internal/templ_components/hidden_input"
	paragraph "github.com/landmaster135/devbox/internal/templ_components/paragraph"
)

const (
	BaseEndpoint            = "/cron-workflow"
	ManualRunEndpoint       = BaseEndpoint + "/manual-run"
	manualWorkflowFieldName = "workflow"
)

const (
	categoryWeatherAutomationID = "weather-automation"
	categoryDailyHeadingID      = "daily-heading"
)

var workflowCategoryDefinitions = []workflowCategoryDefinition{
	{
		ID:          categoryWeatherAutomationID,
		Title:       "Weather & Climate",
		Description: "Workflows that source OpenWeatherMap data and fan out Discord notifications.",
	},
	{
		ID:          categoryDailyHeadingID,
		Title:       "Daily Briefing",
		Description: "Text-generation workflows that prepare Discord-ready summaries.",
	},
}

var (
	workflowDefinitionsByDescription = initWorkflowDefinitions()
	workflowDefinitionsByKey         = buildWorkflowDefinitionsByKey(workflowDefinitionsByDescription)
)

var (
	errWorkflowNotRegistered = errors.New("workflow not registered")
	errWorkflowUnavailable   = errors.New("workflow unavailable")
)

// Handler exposes cron-workflow metadata as a templ-rendered HTML page.
type Handler struct {
	listWorkflows func() ([]usecases.Workflow, error)
}

// NewHandler creates a handler backed by workflow.List.
func NewHandler() *Handler {
	return &Handler{listWorkflows: workflowPkg.List}
}

// NewHandlerWithLister allows supplying a custom workflow lister (primarily for tests).
func NewHandlerWithLister(listFn func() ([]usecases.Workflow, error)) *Handler {
	if listFn == nil {
		listFn = workflowPkg.List
	}
	return &Handler{listWorkflows: listFn}
}

// HandleCronWorkflowPage renders the cron-workflow dashboard page.
func (h *Handler) HandleCronWorkflowPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	workflows, err := h.listWorkflows()
	if err != nil {
		log.Printf("cron-workflow: failed to load workflows: %v", err)
		http.Error(w, "failed to load workflow definitions", http.StatusInternalServerError)
		return
	}

	data, err := buildWorkflowPageData(workflows)
	if err != nil {
		log.Printf("cron-workflow: failed to build workflow page data: %v", err)
		http.Error(w, "failed to render workflow page", http.StatusInternalServerError)
		return
	}

	templ.Handler(CronWorkflowPage(data)).ServeHTTP(w, r)
}

// HandleManualRun triggers a workflow immediately without navigating away from the page.
func (h *Handler) HandleManualRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		respondJSONError(w, http.StatusBadRequest, "invalid form payload")
		return
	}

	workflowKey := strings.TrimSpace(r.FormValue(manualWorkflowFieldName))
	if workflowKey == "" {
		respondJSONError(w, http.StatusBadRequest, "workflow key is required")
		return
	}

	workflows, err := h.listWorkflows()
	if err != nil {
		log.Printf("cron-workflow: failed to load workflows for manual run: %v", err)
		respondJSONError(w, http.StatusInternalServerError, "failed to load workflows")
		return
	}

	process, def, err := workflowProcessByKey(workflows, workflowKey)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, errWorkflowNotRegistered):
			status = http.StatusBadRequest
		case errors.Is(err, errWorkflowUnavailable):
			status = http.StatusServiceUnavailable
		}
		respondJSONError(w, status, err.Error())
		return
	}

	if err := process(r.Context()); err != nil {
		log.Printf("cron-workflow: manual run %s failed: %v", def.Key, err)
		respondJSONError(w, http.StatusInternalServerError, fmt.Sprintf("%s failed: %v", def.DisplayName, err))
		return
	}

	respondJSON(w, http.StatusOK, manualRunResponse{
		Status:   "ok",
		Workflow: def.Key,
		Message:  fmt.Sprintf("%s を実行しました", def.DisplayName),
	})
}

type workflowDefinition struct {
	Key             string
	DisplayName     string
	Summary         string
	CategoryID      string
	ProcessDisplay  string
	ListDescription string
}

type workflowCategoryDefinition struct {
	ID          string
	Title       string
	Description string
}

func initWorkflowDefinitions() map[string]workflowDefinition {
	defs := map[string]workflowDefinition{
		"Daily Tokyo weather notification": {
			Key:            "daily-tokyo-weather",
			DisplayName:    "Daily Tokyo weather notification",
			Summary:        "Fetches a three-day forecast for Tokyo and shares it with the weather Discord channel.",
			CategoryID:     categoryWeatherAutomationID,
			ProcessDisplay: "WorkflowHandler.NotifyWeather",
		},
		"Daily heading Discord notification": {
			Key:            "daily-heading",
			DisplayName:    "Daily heading Discord notification",
			Summary:        "Generates the day's heading template and posts it to the daily heading Discord webhook.",
			CategoryID:     categoryDailyHeadingID,
			ProcessDisplay: "WorkflowHandler.NotifyDailyHeading",
		},
	}
	for desc, def := range defs {
		def.ListDescription = desc
		defs[desc] = def
	}
	return defs
}

func buildWorkflowDefinitionsByKey(defs map[string]workflowDefinition) map[string]workflowDefinition {
	result := make(map[string]workflowDefinition, len(defs))
	for _, def := range defs {
		result[def.Key] = def
	}
	return result
}

type workflowPageData struct {
	Title       string
	Description string
	Categories  []workflowCategoryView
}

type workflowCategoryView struct {
	ID          string
	Title       string
	Description string
	Workflows   []workflowCardView
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

func buildWorkflowPageData(workflows []usecases.Workflow) (workflowPageData, error) {
	categories := make(map[string]*workflowCategoryView, len(workflowCategoryDefinitions))
	for _, catDef := range workflowCategoryDefinitions {
		categories[catDef.ID] = &workflowCategoryView{
			ID:          catDef.ID,
			Title:       catDef.Title,
			Description: catDef.Description,
		}
	}

	for _, wf := range workflows {
		descKey := strings.TrimSpace(wf.Description)
		def, ok := workflowDefinitionsByDescription[descKey]
		if !ok {
			continue
		}

		cronDef, _, err := wf.GetCronDefinition()
		if err != nil {
			return workflowPageData{}, fmt.Errorf("get cron definition for %s: %w", descKey, err)
		}

		card := workflowCardView{
			Key:                 def.Key,
			Name:                def.DisplayName,
			Summary:             def.Summary,
			CronDefinition:      cronDef,
			ManualAction:        ManualRunEndpoint,
			ManualMethod:        http.MethodPost,
			ManualWorkflowField: manualWorkflowFieldName,
			ManualWorkflowValue: def.Key,
			ProcessName:         def.ProcessDisplay,
		}
		categoryView := categories[def.CategoryID]
		categoryView.Workflows = append(categoryView.Workflows, card)
	}

	data := workflowPageData{
		Title:       "Cron Workflow Orchestrator",
		Description: "Workflows registered in workflow.List() are grouped by category for quick inspection.",
	}

	for _, catDef := range workflowCategoryDefinitions {
		catView := categories[catDef.ID]
		if len(catView.Workflows) == 0 {
			continue
		}
		sort.Slice(catView.Workflows, func(i, j int) bool {
			return strings.Compare(catView.Workflows[i].Name, catView.Workflows[j].Name) < 0
		})
		data.Categories = append(data.Categories, *catView)
	}

	return data, nil
}

// CronWorkflowPage produces the templ component for the dashboard.
func CronWorkflowPage(data workflowPageData) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := writeString(w, "<!DOCTYPE html><html lang=\"ja\">"); err != nil {
			return err
		}
		if err := writeString(w, "<head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><title>"); err != nil {
			return err
		}
		if err := writeEscapedString(w, data.Title); err != nil {
			return err
		}
		if err := writeString(w, "</title><style>"+cronWorkflowPageStyles+"</style></head><body><main class=\"page\">"); err != nil {
			return err
		}
		if err := renderHeroSection(ctx, w, data); err != nil {
			return err
		}
		if len(data.Categories) == 0 {
			if err := renderEmptyState(ctx, w); err != nil {
				return err
			}
		} else {
			if err := writeString(w, "<div class=\"page-sections\">"); err != nil {
				return err
			}
			for _, cat := range data.Categories {
				if err := renderCategorySection(ctx, w, cat); err != nil {
					return err
				}
			}
			if err := writeString(w, "</div>"); err != nil {
				return err
			}
		}
		if err := writeString(w, "</main><script>"+cronWorkflowPageScript+"</script></body></html>"); err != nil {
			return err
		}
		return nil
	})
}

func renderHeroSection(ctx context.Context, w io.Writer, data workflowPageData) error {
	if err := writeString(w, "<section class=\"hero\">"); err != nil {
		return err
	}
	if err := paragraph.Text("hero-eyebrow", "Cron Workflow").Render(ctx, w); err != nil {
		return err
	}
	if err := heading.Heading(1, data.Title).Render(ctx, w); err != nil {
		return err
	}
	if data.Description != "" {
		if err := paragraph.Text("hero-description", data.Description).Render(ctx, w); err != nil {
			return err
		}
	}
	return writeString(w, "</section>")
}

func renderEmptyState(ctx context.Context, w io.Writer) error {
	if err := writeString(w, "<section class=\"empty-state\">"); err != nil {
		return err
	}
	if err := heading.Heading(2, "No workflows").Render(ctx, w); err != nil {
		return err
	}
	if err := paragraph.Text("", "workflow.List() did not expose weatherWorkflow or dailyHeadingWorkflow.").Render(ctx, w); err != nil {
		return err
	}
	if err := writeString(w, "</section>"); err != nil {
		return err
	}
	return nil
}

func renderCategorySection(ctx context.Context, w io.Writer, cat workflowCategoryView) error {
	if err := writeString(w, "<section class=\"workflow-category\" id=\""); err != nil {
		return err
	}
	if err := writeEscapedString(w, cat.ID); err != nil {
		return err
	}
	if err := writeString(w, "\"><div class=\"workflow-category__header\">"); err != nil {
		return err
	}
	if err := heading.Heading(2, cat.Title).Render(ctx, w); err != nil {
		return err
	}
	if cat.Description != "" {
		if err := paragraph.Text("", cat.Description).Render(ctx, w); err != nil {
			return err
		}
	}
	if err := writeString(w, "</div><div class=\"workflow-grid\">"); err != nil {
		return err
	}
	for _, wf := range cat.Workflows {
		if err := renderWorkflowCard(ctx, w, wf); err != nil {
			return err
		}
	}
	return writeString(w, "</div></section>")
}

func renderWorkflowCard(ctx context.Context, w io.Writer, card workflowCardView) error {
	if err := writeString(w, "<article class=\"workflow-card\" data-workflow-key=\""); err != nil {
		return err
	}
	if err := writeEscapedString(w, card.Key); err != nil {
		return err
	}
	if err := writeString(w, "\"><div class=\"workflow-card__header\">"); err != nil {
		return err
	}
	if err := paragraph.Content("workflow-card__process", code.Text("", card.ProcessName)).Render(ctx, w); err != nil {
		return err
	}
	if err := heading.Heading(3, card.Name).Render(ctx, w); err != nil {
		return err
	}
	if err := writeString(w, "</div>"); err != nil {
		return err
	}
	if card.Summary != "" {
		if err := paragraph.Text("workflow-card__summary", card.Summary).Render(ctx, w); err != nil {
			return err
		}
	}
	if card.CronDefinition != "" {
		if err := writeString(w, "<div class=\"workflow-card__cron\"><span>CRON</span>"); err != nil {
			return err
		}
		if err := code.Text("", card.CronDefinition).Render(ctx, w); err != nil {
			return err
		}
		if err := writeString(w, "</div>"); err != nil {
			return err
		}
	}
	if card.ManualAction != "" {
		if err := writeString(w, "<form class=\"workflow-card__manual\" method=\""); err != nil {
			return err
		}
		method := card.ManualMethod
		if method == "" {
			method = http.MethodPost
		}
		if err := writeEscapedString(w, strings.ToUpper(method)); err != nil {
			return err
		}
		if err := writeString(w, "\" action=\""); err != nil {
			return err
		}
		action := card.ManualAction
		if err := writeEscapedString(w, action); err != nil {
			return err
		}
		if err := writeString(w, "\" data-endpoint=\""); err != nil {
			return err
		}
		if err := writeEscapedString(w, action); err != nil {
			return err
		}
		if err := writeString(w, "\">"); err != nil {
			return err
		}
		if card.ManualWorkflowField != "" && card.ManualWorkflowValue != "" {
			if err := hiddenInput.HiddenField(card.ManualWorkflowField, card.ManualWorkflowValue).Render(ctx, w); err != nil {
				return err
			}
		}
		if err := button.Submit("Manual Run").Render(ctx, w); err != nil {
			return err
		}
		if err := writeString(w, "</form>"); err != nil {
			return err
		}
	}
	if err := paragraph.Status("workflow-card__status").Render(ctx, w); err != nil {
		return err
	}
	return writeString(w, "</article>")
}

type manualRunResponse struct {
	Status   string `json:"status"`
	Workflow string `json:"workflow"`
	Message  string `json:"message"`
}

func workflowProcessByKey(workflows []usecases.Workflow, workflowKey string) (usecases.ProcessFunc, workflowDefinition, error) {
	def, ok := workflowDefinitionsByKey[workflowKey]
	if !ok {
		return nil, workflowDefinition{}, fmt.Errorf("%w: %s", errWorkflowNotRegistered, workflowKey)
	}

	desc := strings.TrimSpace(def.ListDescription)
	for _, wf := range workflows {
		if strings.TrimSpace(wf.Description) != desc {
			continue
		}
		if wf.Process == nil {
			return nil, def, fmt.Errorf("%w: %s missing process", errWorkflowUnavailable, workflowKey)
		}
		return wf.Process, def, nil
	}

	return nil, def, fmt.Errorf("%w: %s not scheduled", errWorkflowUnavailable, workflowKey)
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("cron-workflow: failed to encode JSON response: %v", err)
	}
}

func respondJSONError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"status":  "error",
		"message": message,
		"error":   message,
	})
}

func writeString(w io.Writer, value string) error {
	_, err := io.WriteString(w, value)
	return err
}

func writeEscapedString(w io.Writer, value string) error {
	_, err := io.WriteString(w, templ.EscapeString(value))
	return err
}

const cronWorkflowPageStyles = `
:root {
	color-scheme: light dark;
	font-family: "Inter", "Segoe UI", system-ui, -apple-system, BlinkMacSystemFont, sans-serif;
}
* {
	box-sizing: border-box;
}
body {
	margin: 0;
	background: #f5f6fb;
	color: #111827;
}
main.page {
	max-width: 1100px;
	margin: 0 auto;
	padding: 48px 20px 80px;
}
.hero {
	background: linear-gradient(120deg, #eef2ff, #e0f2fe);
	border: 1px solid rgba(99, 102, 241, 0.25);
	border-radius: 20px;
	padding: 32px;
	box-shadow: 0 15px 35px rgba(15, 23, 42, 0.08);
}
.hero-eyebrow {
	text-transform: uppercase;
	letter-spacing: 0.12em;
	font-size: 0.75rem;
	color: #6366f1;
	margin: 0 0 8px;
}
.hero h1 {
	margin: 0 0 16px;
	font-size: 2rem;
}
.hero-description {
	margin: 0 0 12px;
	color: #374151;
}
.hero-source,
.hero-endpoint {
	margin: 4px 0;
	font-size: 0.9rem;
	color: #4b5563;
}
.hero code {
	background: rgba(255, 255, 255, 0.8);
	padding: 4px 8px;
	border-radius: 6px;
}
.page-sections {
	margin-top: 32px;
	display: flex;
	flex-direction: column;
	gap: 32px;
}
.workflow-category {
	background: #ffffff;
	border-radius: 20px;
	border: 1px solid rgba(148, 163, 184, 0.4);
	padding: 28px;
	box-shadow: 0 20px 45px rgba(15, 23, 42, 0.06);
}
.workflow-category__header > h2 {
	margin: 0 0 4px;
}
.workflow-category__header > p {
	margin: 0;
	color: #4b5563;
}
.workflow-grid {
	margin-top: 20px;
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
	gap: 20px;
}
.workflow-card {
	background: #f8fafc;
	border-radius: 16px;
	border: 1px solid rgba(203, 213, 225, 0.7);
	padding: 20px;
	display: flex;
	flex-direction: column;
	gap: 12px;
}
.workflow-card__process {
	margin: 0;
	font-size: 0.85rem;
	color: #6b7280;
}
.workflow-card__process code {
	background: rgba(255, 255, 255, 0.75);
	padding: 2px 6px;
	border-radius: 5px;
}
.workflow-card__header h3 {
	margin: 6px 0 0;
}
.workflow-card__summary {
	margin: 0;
	color: #475569;
	line-height: 1.4;
}
.workflow-card__cron {
	display: flex;
	flex-direction: column;
	font-size: 0.9rem;
}
.workflow-card__cron span {
	font-weight: 600;
	color: #0f172a;
}
.workflow-card__cron code {
	margin-top: 4px;
	padding: 4px 6px;
	background: #e2e8f0;
	border-radius: 6px;
	font-family: "JetBrains Mono", "SFMono-Regular", Consolas, monospace;
}
.workflow-card__manual {
	margin-top: auto;
	display: flex;
	justify-content: flex-start;
}
.workflow-card__manual button {
	padding: 10px 16px;
	background: #4f46e5;
	color: #ffffff;
	border: none;
	border-radius: 8px;
	cursor: pointer;
	font-weight: 600;
}
.workflow-card__manual button:hover {
	background: #4338ca;
}
.workflow-card__status {
	font-size: 0.85rem;
	margin: 8px 0 0;
	min-height: 1.2em;
	color: #4b5563;
}
.workflow-card__status[data-status="success"] {
	color: #0f766e;
}
.workflow-card__status[data-status="error"] {
	color: #dc2626;
}
.workflow-card__status[data-status="pending"] {
	color: #92400e;
}
.empty-state {
	margin-top: 32px;
	padding: 40px;
	text-align: center;
	border-radius: 16px;
	border: 1px dashed rgba(99, 102, 241, 0.5);
	color: #4b5563;
}
@media (prefers-color-scheme: dark) {
	body { background: #020617; color: #e2e8f0; }
	.hero { background: linear-gradient(120deg, #312e81, #1d4ed8); border-color: rgba(99,102,241,0.5); }
	.hero-description, .hero-source, .hero-endpoint { color: #cbd5f5; }
	.workflow-category { background: #0f172a; border-color: rgba(59, 130, 246, 0.4); }
	.workflow-card { background: #1e293b; border-color: rgba(59, 130, 246, 0.25); }
	.workflow-card__summary { color: #cbd5f5; }
	.workflow-card__cron span { color: #f8fafc; }
	.workflow-card__cron code { background: #0f172a; }
	.empty-state { border-color: rgba(99, 102, 241, 0.8); color: #cbd5f5; }
}
`

const cronWorkflowPageScript = `(() => {
	const forms = document.querySelectorAll(".workflow-card__manual");
	forms.forEach((form) => {
		form.addEventListener("submit", async (event) => {
			event.preventDefault();
			const statusEl = form.parentElement?.querySelector("[data-manual-run-status]");
			if (statusEl) {
				statusEl.textContent = "手動実行中...";
				statusEl.dataset.status = "pending";
			}
			const formData = new FormData(form);
			const payload = new URLSearchParams();
			formData.forEach((value, key) => {
				payload.append(key, value.toString());
			});
			try {
				const response = await fetch(form.action, {
					method: form.method || "POST",
					headers: {
						"Content-Type": "application/x-www-form-urlencoded;charset=UTF-8",
					},
					body: payload.toString(),
				});
				let data = null;
				try {
					data = await response.json();
				} catch (error) {
					data = null;
				}
				if (statusEl) {
					if (response.ok) {
						statusEl.textContent = (data && (data.message || data.status)) || "完了しました";
						statusEl.dataset.status = "success";
					} else {
						const message = (data && (data.error || data.message)) || ("HTTP " + response.status);
						statusEl.textContent = message;
						statusEl.dataset.status = "error";
					}
				}
			} catch (error) {
				if (statusEl) {
					statusEl.textContent = "リクエストに失敗しました";
					statusEl.dataset.status = "error";
				}
			}
		});
	});
})();`
