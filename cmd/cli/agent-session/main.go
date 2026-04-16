package main

import (
	"fmt"
	"os"

	config "github.com/landmaster135/devbox/internal/agent_session/config"
	usecases "github.com/landmaster135/devbox/internal/agent_session/usecases"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	switch cfg.Operation {
	case "retrieve-session":
		handleRetrieveSessionOperation(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応のoperationです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

func handleRetrieveSessionOperation(cfg *config.Config) {
	service := usecases.NewAgentSessionService()

	result, err := service.RetrieveSessions(usecases.RetrieveSessionsInput{
		AgentType:    cfg.AgentType,
		Limit:        cfg.Limit,
		StartDate:    cfg.StartDateValue,
		EndDate:      cfg.EndDateValue,
		AgentHomeDir: cfg.AgentHomeDir,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}
