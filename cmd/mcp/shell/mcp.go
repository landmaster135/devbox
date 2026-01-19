package shell

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	"github.com/landmaster135/devbox/internal/shell/domain"
	usecases "github.com/landmaster135/devbox/internal/shell/usecases"
)

func handleShellExecute(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command, err := requireStringSlice(request, "command")
	if err != nil {
		return nil, err
	}
	if len(command) == 0 {
		return nil, fmt.Errorf("commandは1要素以上必要です")
	}

	workdir := request.GetString("workdir", "")
	if workdir == "" {
		workdir = request.GetString("cwd", "")
	}

	baseDir := request.GetString("base_dir", "")
	if baseDir == "" {
		baseDir = request.GetString("basedir", "")
	}
	if baseDir == "" {
		baseDir = os.Getenv("SHELL_BASE_DIRECTORY")
	}

	timeout := request.GetInt("timeout_ms", request.GetInt("timeout", 0))
	if timeout < 0 {
		return nil, fmt.Errorf("timeout_msは0以上で指定してください")
	}

	env, err := getEnvMap(request, "env")
	if err != nil {
		return nil, err
	}

	sandboxValue := request.GetString("sandbox_permissions", request.GetString("sandbox", string(domain.SandboxPermissionsUseDefault)))
	if sandboxValue == "" {
		sandboxValue = string(domain.SandboxPermissionsUseDefault)
	}
	sandbox, err := domain.ParseSandboxPermissions(sandboxValue)
	if err != nil {
		return nil, err
	}

	justification := request.GetString("justification", request.GetString("reason", ""))

	input := &usecases.ExecuteCommandInput{
		Command:            command,
		WorkDir:            workdir,
		BaseDir:            baseDir,
		TimeoutMs:          uint64(timeout),
		Env:                env,
		SandboxPermissions: sandbox,
		Justification:      justification,
	}

	service := usecases.NewShellService()
	result, err := service.ExecuteCommand(input)
	if err != nil {
		return nil, err
	}

	payload, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("結果のシリアライズに失敗しました: %w", err)
	}

	return mcp.NewToolResultText(string(payload)), nil
}

func handleListDenied(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	service := usecases.NewShellService()
	commands := service.ListDeniedCommands()

	payload := map[string][]string{
		"commands": commands,
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(data)), nil
}

func requireStringSlice(request mcp.CallToolRequest, key string) ([]string, error) {
	raw, ok := request.GetArguments()[key]
	if !ok {
		return nil, fmt.Errorf("%sは必須です", key)
	}
	items, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%sは文字列の配列で指定してください", key)
	}
	values := make([]string, 0, len(items))
	for i, v := range items {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%s[%d]は文字列である必要があります", key, i)
		}
		values = append(values, str)
	}
	return values, nil
}

func getEnvMap(request mcp.CallToolRequest, key string) (map[string]string, error) {
	raw, ok := request.GetArguments()[key]
	if !ok {
		return nil, nil
	}
	valueMap, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%sはオブジェクトで指定してください", key)
	}
	result := make(map[string]string, len(valueMap))
	for k, v := range valueMap {
		result[k] = fmt.Sprintf("%v", v)
	}
	return result, nil
}

func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("shell_prompt",
		mcp.WithPromptDescription("Guardrails for shell operations"),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for the shell executor.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("Use shell_execute responsibly. Favor read-only commands unless escalation is justified."),
				},
			},
		}, nil
	})
	return s
}

func setShellTools(s *server.MCPServer) *server.MCPServer {
	executeTool := mcp.NewTool("execute",
		mcp.WithDescription("Codex CLIのshellツール互換エクゼキュータ。command配列をそのままexecvpに渡して実行します。"),
		mcp.WithArray("command",
			mcp.Required(),
			mcp.Description("execvpに渡すコマンドの配列"),
		),
		mcp.WithString("workdir",
			mcp.Description("作業ディレクトリ。未指定時はbase_dirもしくはカレントディレクトリ"),
		),
		mcp.WithString("base_dir",
			mcp.Description("ベースディレクトリ。未指定時は環境変数SHELL_BASE_DIRECTORYまたはカレント"),
		),
		mcp.WithNumber("timeout_ms",
			mcp.Description("ミリ秒単位のタイムアウト。0の場合は60秒"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("timeout_msのエイリアス"),
		),
		mcp.WithObject("env",
			mcp.Description("KEY:VALUE形式の環境変数マップ"),
		),
		mcp.WithString("sandbox_permissions",
			mcp.Description("use_defaultまたはrequire_escalated"),
		),
		mcp.WithString("justification",
			mcp.Description("sandbox_permissions: require_escalated を設定した時の理由"),
		),
	)
	s.AddTool(executeTool, handleShellExecute)

	deniedTool := mcp.NewTool("list_denied_commands",
		mcp.WithDescription("defaultDeniedCommandsに登録されている危険コマンド一覧を返します"),
	)
	s.AddTool(deniedTool, handleListDenied)

	return s
}

func BuildShellServer() {
	s := server.NewMCPServer(
		"Shell CLI Server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	s = setShellTools(s)
	s = addPromptIntoServer(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}
