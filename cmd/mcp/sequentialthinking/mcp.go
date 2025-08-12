package sequentialthinking

import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	server "github.com/mark3labs/mcp-go/server"

	usecases "github.com/landmaster135/devbox/internal/sequential_thinking/usecases"
)

const (
	version = "0.2.0"
)

// #==============================================================#
// ##          Handlers                                          ##
// #==============================================================#
func handleSequentialThinking(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	// SequentialThinkingServiceを初期化
	service := usecases.NewSequentialThinkingService()
	result, err := service.HandleSequentialThinking(args)
	if err != nil {
		return nil, fmt.Errorf("シーケンシャルシンキング処理に失敗しました: %v", err)
	}

	return mcp.NewToolResultText(result), nil
}

// #==============================================================#
// ##          Servers                                           ##
// #==============================================================#
func addPromptIntoServer(s *server.MCPServer) *server.MCPServer {
	prompt := mcp.NewPrompt("system_prompt_01",
		mcp.WithPromptDescription("This is a prompt for the sequential thinking."),
	)
	s.AddPrompt(prompt, func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "System prompt for the sequential thinking.",
			Messages: []mcp.PromptMessage{
				{
					Role:    mcp.RoleAssistant,
					Content: mcp.NewTextContent("You use this great sequential thinking well."),
				},
			},
		}, nil
	})
	return s
}

func createSequentialThinkingServer() *server.MCPServer {
	// MCPサーバーを作成
	s := server.NewMCPServer(
		"Sequential Thinking Server",
		version,
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	// ツール: シーケンシャルシンキング
	sequentialThinkingTool := mcp.NewTool("sequentialthinking",
		mcp.WithDescription(`A detailed tool for dynamic and reflective problem-solving through thoughts.
This tool helps analyze problems through a flexible thinking process that can adapt and evolve.
Each thought can build on, question, or revise previous insights as understanding deepens.

When to use this tool:
- Breaking down complex problems into steps
- Planning and design with room for revision
- Analysis that might need course correction
- Problems where the full scope might not be clear initially
- Problems that require a multi-step solution
- Tasks that need to maintain context over multiple steps
- Situations where irrelevant information needs to be filtered out

Key features:
- You can adjust total_thoughts up or down as you progress
- You can question or revise previous thoughts
- You can add more thoughts even after reaching what seemed like the end
- You can express uncertainty and explore alternative approaches
- Not every thought needs to build linearly - you can branch or backtrack
- Generates a solution hypothesis
- Verifies the hypothesis based on the Chain of Thought steps
- Repeats the process until satisfied
- Provides a correct answer

Parameters explained:
- thought: Your current thinking step, which can include:
* Regular analytical steps
* Revisions of previous thoughts
* Questions about previous decisions
* Realizations about needing more analysis
* Changes in approach
* Hypothesis generation
* Hypothesis verification
- next_thought_needed: True if you need more thinking, even if at what seemed like the end
- thought_number: Current number in sequence (can go beyond initial total if needed)
- total_thoughts: Current estimate of thoughts needed (can be adjusted up/down)
- is_revision: A boolean indicating if this thought revises previous thinking
- revises_thought: If is_revision is true, which thought number is being reconsidered
- branch_from_thought: If branching, which thought number is the branching point
- branch_id: Identifier for the current branch (if any)
- needs_more_thoughts: If reaching end but realizing more thoughts needed

You should:
1. Start with an initial estimate of needed thoughts, but be ready to adjust
2. Feel free to question or revise previous thoughts
3. Don't hesitate to add more thoughts if needed, even at the "end"
4. Express uncertainty when present
5. Mark thoughts that revise previous thinking or branch into new paths
6. Ignore information that is irrelevant to the current step
7. Generate a solution hypothesis when appropriate
8. Verify the hypothesis based on the Chain of Thought steps
9. Repeat the process until satisfied with the solution
10. Provide a single, ideally correct answer as the final output
11. Only set next_thought_needed to false when truly done and a satisfactory answer is reached`),
		mcp.WithString("thought",
			mcp.Required(),
			mcp.Description("Your current thinking step"),
		),
		mcp.WithBoolean("nextThoughtNeeded",
			mcp.Required(),
			mcp.Description("Whether another thought step is needed"),
		),
		mcp.WithNumber("thoughtNumber",
			mcp.Required(),
			mcp.Description("Current thought number"),
			mcp.Min(1),
		),
		mcp.WithNumber("totalThoughts",
			mcp.Required(),
			mcp.Description("Estimated total thoughts needed"),
			mcp.Min(1),
		),
		mcp.WithBoolean("isRevision",
			mcp.Description("Whether this revises previous thinking"),
		),
		mcp.WithNumber("revisesThought",
			mcp.Description("Which thought is being reconsidered"),
			mcp.Min(1),
		),
		mcp.WithNumber("branchFromThought",
			mcp.Description("Branching point thought number"),
			mcp.Min(1),
		),
		mcp.WithString("branchId",
			mcp.Description("Branch identifier"),
		),
		mcp.WithBoolean("needsMoreThoughts",
			mcp.Description("If more thoughts are needed"),
		),
	)

	s.AddTool(sequentialThinkingTool, handleSequentialThinking)

	// プロンプト
	s = addPromptIntoServer(s)

	return s
}

func BuildSequentialThinkingServer() {
	s := createSequentialThinkingServer()

	// サーバーを起動
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
