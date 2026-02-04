package main

import (
	"encoding/json"
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/taskfile/config"
	usecases "github.com/landmaster135/devbox/internal/taskfile/usecases"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	switch cfg.Operation {
	case config.OperationInspect:
		handleInspect(cfg)
	case config.OperationFill:
		handleFill(cfg)
	case config.OperationNew:
		handleNew(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未サポートのoperationです: %s\n", cfg.Operation)
		os.Exit(1)
	}
}

type baseResponse struct {
	Operation string `json:"operation"`
	Status    string `json:"status"`
	Code      int    `json:"code"`
	Message   string `json:"message,omitempty"`
}

type inspectResponse struct {
	baseResponse
	MissingFields []string `json:"missing_fields,omitempty"`
}

type taskfileResponse struct {
	baseResponse
	TaskfilePath string `json:"taskfile_path"`
}

type errorResponse struct {
	baseResponse
	Error string `json:"error"`
}

const (
	statusSuccess      = "success"
	statusMissingField = "missing_fields"
	statusRegistered   = "registered"
	statusNoop         = "noop"
	statusError        = "error"
)

const (
	codeOK            = 200
	codeBadRequest    = 400
	codeInternalError = 500
)

func handleInspect(cfg *config.Config) {
	service := usecases.NewService()
	result, err := service.Inspect(cfg.TaskType, cfg.TaskfilePath)
	if err != nil {
		outputError(config.OperationInspect, err)
	}

	if result.HasMissingFields() {
		outputJSONAndExit(1, inspectResponse{
			baseResponse: baseResponse{
				Operation: config.OperationInspect,
				Status:    statusMissingField,
				Code:      codeOK,
				Message:   fmt.Sprintf("不足しているフィールドが %d 個見つかりました", len(result.MissingFields)),
			},
			MissingFields: result.MissingFields,
		})
		return
	}

	outputJSON(inspectResponse{
		baseResponse: baseResponse{
			Operation: config.OperationInspect,
			Status:    statusSuccess,
			Code:      codeOK,
			Message:   "Taskfileには参照Taskfileのすべてのフィールドが含まれています。",
		},
	})
}

func handleFill(cfg *config.Config) {
	service := usecases.NewService()
	updated, err := service.Fill(cfg.TaskType, cfg.TaskfilePath)
	if err != nil {
		outputError(config.OperationFill, err)
	}

	if updated {
		outputJSON(taskfileResponse{
			baseResponse: baseResponse{
				Operation: config.OperationFill,
				Status:    statusRegistered,
				Code:      codeOK,
				Message:   "Taskfileの空欄フィールドをテンプレートの値で補完しました。",
			},
			TaskfilePath: cfg.TaskfilePath,
		})
		return
	}

	outputJSON(taskfileResponse{
		baseResponse: baseResponse{
			Operation: config.OperationFill,
			Status:    statusNoop,
			Code:      codeOK,
			Message:   "補完対象の空欄フィールドは見つかりませんでした。",
		},
		TaskfilePath: cfg.TaskfilePath,
	})
}

func handleNew(cfg *config.Config) {
	service := usecases.NewService()
	if err := service.Create(cfg.TaskType, cfg.TaskfilePath); err != nil {
		outputError(config.OperationNew, err)
	}

	outputJSON(taskfileResponse{
		baseResponse: baseResponse{
			Operation: config.OperationNew,
			Status:    statusRegistered,
			Code:      codeOK,
			Message:   fmt.Sprintf("Taskfileを新規作成しました: %s", cfg.TaskfilePath),
		},
		TaskfilePath: cfg.TaskfilePath,
	})
}

func outputJSON(payload interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(payload); err != nil {
		fmt.Fprintf(os.Stderr, "エラー: JSON出力に失敗しました: %v\n", err)
		os.Exit(1)
	}
}

func outputJSONAndExit(code int, payload interface{}) {
	outputJSON(payload)
	if code != 0 {
		os.Exit(code)
	}
}

func outputError(operation string, err error) {
	outputJSONAndExit(1, errorResponse{
		baseResponse: baseResponse{
			Operation: operation,
			Status:    statusError,
			Code:      codeInternalError,
		},
		Error: err.Error(),
	})
}
