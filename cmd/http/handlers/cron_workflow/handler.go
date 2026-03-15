package cron_workflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	workflowCLIPkg "github.com/landmaster135/devbox/cmd/cli/cron-workflow/workflow"
	page "github.com/landmaster135/devbox/cmd/http/handlers/cron_workflow/page"
	workflow "github.com/landmaster135/devbox/cmd/http/handlers/cron_workflow/workflow"
	usecases "github.com/landmaster135/devbox/internal/cron_workflow/usecases"
	logging "github.com/landmaster135/devbox/internal/logging"
)

const (
	BaseEndpoint            = "/cron-workflow"
	ManualRunEndpoint       = BaseEndpoint + "/manual-run"
	manualWorkflowFieldName = "workflow"
)

// Handler exposes cron-workflow metadata as a templ-rendered HTML page.
type Handler struct {
	listWorkflows           func(logger *logging.StructuredLogger) ([]usecases.Workflow, error)
	logger                  *logging.StructuredLogger
	manualWorkflowFieldName string
	manualRunEndpoint       string
}

// NewHandler creates a handler backed by workflow.List.
func NewHandler(logger *logging.StructuredLogger) *Handler {
	return &Handler{
		listWorkflows:           workflowCLIPkg.List,
		logger:                  logging.Ensure(logger),
		manualWorkflowFieldName: manualWorkflowFieldName,
		manualRunEndpoint:       ManualRunEndpoint,
	}
}

// NewHandlerWithLister allows supplying a custom workflow lister (primarily for tests).
func NewHandlerWithLister(listFn func(logger *logging.StructuredLogger) ([]usecases.Workflow, error), logger *logging.StructuredLogger) *Handler {
	if listFn == nil {
		listFn = workflowCLIPkg.List
	}
	return &Handler{
		listWorkflows:           listFn,
		logger:                  logging.Ensure(logger),
		manualWorkflowFieldName: manualWorkflowFieldName,
		manualRunEndpoint:       ManualRunEndpoint,
	}
}

// HandleCronWorkflowPage renders the cron-workflow dashboard page.
func (h *Handler) HandleCronWorkflowPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	pageLogger := h.logger.WithTags("page")
	workflows, err := h.listWorkflows(h.logger)
	if err != nil {
		pageLogger.WithTags("list").Errorf("failed to load workflows: %v", err)
		http.Error(w, "failed to load workflow definitions", http.StatusInternalServerError)
		return
	}

	page.Serve(w, r, pageLogger, workflows, h.manualRunEndpoint, h.manualWorkflowFieldName)
}

// HandleManualRun triggers a workflow immediately without navigating away from the page.
func (h *Handler) HandleManualRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "only POST is allowed", http.StatusMethodNotAllowed)
		return
	}

	manualLogger := h.logger.WithTags("manual-run")
	if err := r.ParseForm(); err != nil {
		respondJSONError(manualLogger, w, http.StatusBadRequest, "invalid form payload")
		return
	}

	workflowKey := strings.TrimSpace(r.FormValue(manualWorkflowFieldName))
	if workflowKey == "" {
		respondJSONError(manualLogger, w, http.StatusBadRequest, "workflow key is required")
		return
	}

	workflows, err := h.listWorkflows(h.logger)
	if err != nil {
		manualLogger.WithTags("list").Errorf("failed to load workflows for manual run: %v", err)
		respondJSONError(manualLogger, w, http.StatusInternalServerError, "failed to load workflows")
		return
	}

	process, def, err := workflow.WorkflowProcessByKey(workflows, workflowKey)
	if err != nil {
		status := http.StatusInternalServerError
		switch {
		case errors.Is(err, workflow.ErrWorkflowNotRegistered):
			status = http.StatusBadRequest
		case errors.Is(err, workflow.ErrWorkflowUnavailable):
			status = http.StatusServiceUnavailable
		}
		respondJSONError(manualLogger, w, status, err.Error())
		return
	}

	if err := process(r.Context()); err != nil {
		manualLogger.WithTags(def.Key).Errorf("manual run failed: %v", err)
		respondJSONError(manualLogger, w, http.StatusInternalServerError, fmt.Sprintf("%s failed: %v", def.DisplayName, err))
		return
	}

	respondJSON(manualLogger, w, http.StatusOK, manualRunResponse{
		Status:   "ok",
		Workflow: def.Key,
		Message:  fmt.Sprintf("%s を実行しました", def.DisplayName),
	})
}

type manualRunResponse struct {
	Status   string `json:"status"`
	Workflow string `json:"workflow"`
	Message  string `json:"message"`
}

func respondJSON(logger *logging.StructuredLogger, w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		logging.Ensure(logger).WithTags("response").Errorf("failed to encode JSON response: %v", err)
	}
}

func respondJSONError(logger *logging.StructuredLogger, w http.ResponseWriter, status int, message string) {
	respondJSON(logger, w, status, map[string]string{
		"status":  "error",
		"message": message,
		"error":   message,
	})
}
