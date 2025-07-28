package mcp

import (
	"fmt"
	"os"

	arithmetic_calculator "github.com/landmaster135/devbox/cmd/mcp/arithmetic_calculator"
	brave_search "github.com/landmaster135/devbox/cmd/mcp/brave_search"
	context7 "github.com/landmaster135/devbox/cmd/mcp/context7"
	datetime_calculator "github.com/landmaster135/devbox/cmd/mcp/datetime_calculator"
	duckduckgo_search "github.com/landmaster135/devbox/cmd/mcp/duckduckgo_search"
	everart "github.com/landmaster135/devbox/cmd/mcp/everart"
	figma "github.com/landmaster135/devbox/cmd/mcp/figma"
	filesystem "github.com/landmaster135/devbox/cmd/mcp/filesystem"
	git_commit_history_retriever "github.com/landmaster135/devbox/cmd/mcp/git_commit_history_retriever"
	git_diff_recorder "github.com/landmaster135/devbox/cmd/mcp/git_diff_recorder"
	github "github.com/landmaster135/devbox/cmd/mcp/github"
	http_request "github.com/landmaster135/devbox/cmd/mcp/http_request"
	postgresql "github.com/landmaster135/devbox/cmd/mcp/postgresql"
	sequentialthinking "github.com/landmaster135/devbox/cmd/mcp/sequentialthinking"
	service_implementing_viewer "github.com/landmaster135/devbox/cmd/mcp/service_implementing_viewer"
	gdrive "github.com/landmaster135/devbox/cmd/mcp/gdrive"
	// shell "github.com/landmaster135/devbox/cmd/mcp/shell" // TODO: unapplicable for WSL...
	timezone "github.com/landmaster135/devbox/cmd/mcp/timezone"
	util "github.com/landmaster135/devbox/cmd/mcp/util"
	youtube_transcript "github.com/landmaster135/devbox/cmd/mcp/youtube_transcript"
)

func Router() {
	args := os.Args
	// check arguments
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run main.go [arguments]")
		fmt.Fprintln(os.Stderr, "arguments are lack")
		os.Exit(1)
	}
	for i, arg := range args[1:] {
		fmt.Fprintf(os.Stderr, "argument %d: %s\n", i+1, arg)
	}

	util.OutLog("main: building mcp server...")
	a1 := args[1]
	switch a1 {
	case "arith_calc":
		arithmetic_calculator.BuildArithCalculatorServer()
	case "context7":
		context7.BuildContext7Server()
	case "datetime_calc":
		datetime_calculator.BuildTimeCalculatorServer()
	case "http_request":
		http_request.BuildMcpServer()
	case "brave_web_search":
		brave_search.BuildBraveSearchServer()
	case "duckduckgo_search":
		duckduckgo_search.BuildDuckDuckGoSearchServer()
	case "timezone":
		timezone.BuildTimezoneServer()
	case "filesystem":
		filesystem.BuildFileSystemServer()
	case "youtube_transcript":
		youtube_transcript.BuildYouTubeTranscriptServer()
	case "github":
		github.BuildGitHubServer()
	case "postgresql":
		postgresql.BuildPostgreSQLServer()
	case "everart":
		everart.BuildEverArtServer()
	case "sequentialthinking":
		sequentialthinking.BuildSequentialThinkingServer()
	case "figma":
		figma.BuildFigmaServer()
	case "gdrive":
		gdrive.BuildGoogleDriveServer()
	case "git_diff_recorder":
		git_diff_recorder.BuildMcpServer()
	case "git_commit_history_retriever":
		git_commit_history_retriever.BuildMcpServer()
	case "service_implementing_viewer":
		service_implementing_viewer.BuildServiceImplementingViewerServer()
	default:
		fmt.Fprintln(os.Stderr, "argument is invalid")
		os.Exit(1)
	}
	util.OutLog("main: built mcp server!")
}
