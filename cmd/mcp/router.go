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
	gdrive "github.com/landmaster135/devbox/cmd/mcp/gdrive"
	git_commit_history_retriever "github.com/landmaster135/devbox/cmd/mcp/git_commit_history_retriever"
	git_diff_recorder "github.com/landmaster135/devbox/cmd/mcp/git_diff_recorder"
	github "github.com/landmaster135/devbox/cmd/mcp/github"
	http_request "github.com/landmaster135/devbox/cmd/mcp/http_request"
	notion_sync "github.com/landmaster135/devbox/cmd/mcp/notion_sync"
	open_weather_map "github.com/landmaster135/devbox/cmd/mcp/open_weather_map"
	ops_for_golang "github.com/landmaster135/devbox/cmd/mcp/ops_for_golang"
	persona_extraction "github.com/landmaster135/devbox/cmd/mcp/persona_extraction"
	plan "github.com/landmaster135/devbox/cmd/mcp/plan"
	postgresql "github.com/landmaster135/devbox/cmd/mcp/postgresql"
	sequentialthinking "github.com/landmaster135/devbox/cmd/mcp/sequentialthinking"
	service_implementing_viewer "github.com/landmaster135/devbox/cmd/mcp/service_implementing_viewer"
	web_clipper "github.com/landmaster135/devbox/cmd/mcp/web_clipper"

	shell "github.com/landmaster135/devbox/cmd/mcp/shell"
	timezone "github.com/landmaster135/devbox/cmd/mcp/timezone"
	util "github.com/landmaster135/devbox/cmd/mcp/util"
	weather_notificator "github.com/landmaster135/devbox/cmd/mcp/weather_notificator"
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
	case "plan":
		plan.BuildPlanServer()
	case "persona_extraction":
		persona_extraction.BuildPersonaExtractionServer()
	case "notion_sync":
		notion_sync.BuildNotionSyncServer()
	case "service_implementing_viewer":
		service_implementing_viewer.BuildServiceImplementingViewerServer()
	case "open_weather_map":
		open_weather_map.BuildOpenWeatherMapServer()
	case "weather_notificator":
		weather_notificator.BuildWeatherNotificatorServer()
	case "ops_for_golang":
		ops_for_golang.BuildGolangOpsServer()
	case "shell":
		shell.BuildShellServer()
	case "web_clipper":
		web_clipper.BuildWebClipperServer()
	default:
		fmt.Fprintln(os.Stderr, "argument is invalid")
		os.Exit(1)
	}
	util.OutLog("main: built mcp server!")
}
