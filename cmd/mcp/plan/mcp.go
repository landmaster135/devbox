package plan

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/plan/usecases"
)

const (
	toolVersion = "0.1.0"
)

func handlePlanUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	payload, err := decodePlanArguments(args)
	if err != nil {
		return nil, fmt.Errorf("引数の解析に失敗しました: %w", err)
	}

	service := usecases.NewPlanService()
	summary, err := service.HandleUpdatePlan(payload)
	if err != nil {
		return nil, fmt.Errorf("プランの処理に失敗しました: %w", err)
	}

	return mcp.NewToolResultText(summary), nil
}

func decodePlanArguments(args map[string]interface{}) (usecases.UpdatePlanRequest, error) {
	var payload usecases.UpdatePlanRequest
	raw, err := json.Marshal(args)
	if err != nil {
		return payload, fmt.Errorf("引数のシリアライズに失敗しました: %w", err)
	}
	if len(raw) == 0 || string(raw) == "null" {
		return payload, fmt.Errorf("plan引数が指定されていません")
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return payload, fmt.Errorf("引数のデシリアライズに失敗しました: %w", err)
	}
	return payload, nil
}

func createPlanTool() mcp.Tool {
	return mcp.NewTool("update_plan",
		mcp.WithDescription(`Structured planning helper. Use this to send your current multi-step plan to the client UI so the user can track progress.`),
		mcp.WithString("explanation",
			mcp.Description("Optional narrative that summarizes the approach."),
		),
		mcp.WithArray("plan",
			mcp.Required(),
			mcp.Description("List of plan items (objects with step + status fields). Allowed statuses: pending, in_progress, completed. At most one item may be in_progress."),
			mcp.Items(usecases.PlanItemJSONSchema()),
		),
		mcp.WithReadOnlyHintAnnotation(true),
		mcp.WithIdempotentHintAnnotation(true),
	)
}

func BuildPlanServer() {
	s := server.NewMCPServer(
		"Plan Tool Server",
		toolVersion,
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
		server.WithLogging(),
	)

	s.AddTool(createPlanTool(), handlePlanUpdate)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Plan server error: %v\n", err)
	}
}
