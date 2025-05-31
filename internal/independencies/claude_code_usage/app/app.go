package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/landmaster135/devbox/internal/independencies/claude_code_usage/config"
	"github.com/landmaster135/devbox/internal/independencies/claude_code_usage/internal"
)

// App represents the main application
type App struct {
	config    *config.AppConfig
	loader    *internal.DataLoader
	calc      *internal.Calculator
	formatter *internal.Formatter
}

// NewApp creates a new App instance
func NewApp(cfg *config.AppConfig) *App {
	return &App{
		config:    cfg,
		calc:      internal.NewCalculator(),
		formatter: internal.NewFormatter(),
	}
}

// Run executes the application
func (a *App) Run(stdout, stderr io.Writer) int {
	// Set defaults
	a.config.SetDefaults()

	// Handle help and version
	if a.config.ShowHelp {
		a.printHelp(stdout)
		return 0
	}

	if a.config.ShowVersion {
		a.printVersion(stdout)
		return 0
	}

	// Set Claude path if not specified
	if a.config.ClaudePath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(stderr, "Error: failed to get home directory: %v\n", err)
			return 1
		}
		a.config.ClaudePath = filepath.Join(homeDir, ".claude")
	}

	// Create internal config
	internalConfig := &internal.Config{
		ClaudePath: a.config.ClaudePath,
		Since:      a.config.Since,
		Until:      a.config.Until,
		JSON:       a.config.OutputJSON,
	}

	a.loader = internal.NewDataLoader(internalConfig)

	// Execute command
	switch a.config.Command {
	case "daily":
		return a.runDailyCommand(stdout, stderr)
	case "session":
		return a.runSessionCommand(stdout, stderr)
	default:
		fmt.Fprintf(stderr, "Error: unknown command: %s\n", a.config.Command)
		return 1
	}
}

// runDailyCommand executes the daily report command
func (a *App) runDailyCommand(stdout, stderr io.Writer) int {
	data, err := a.loader.LoadUsageData()
	if err != nil {
		fmt.Fprintf(stderr, "Error: failed to load usage data: %v\n", err)
		return 1
	}

	if len(data) == 0 {
		fmt.Fprintln(stdout, "No usage data found.")
		return 0
	}

	dailyUsage := a.calc.AggregateDailyUsage(data)
	if len(dailyUsage) == 0 {
		fmt.Fprintln(stdout, "No usage data found for the specified date range.")
		return 0
	}

	// Add total row
	total := a.calc.CalculateDailyTotal(dailyUsage)
	dailyUsage = append(dailyUsage, total)

	internalConfig := &internal.Config{
		JSON: a.config.OutputJSON,
	}

	if err := a.formatter.FormatDailyReport(dailyUsage, internalConfig); err != nil {
		fmt.Fprintf(stderr, "Error: failed to format report: %v\n", err)
		return 1
	}

	return 0
}

// runSessionCommand executes the session report command
func (a *App) runSessionCommand(stdout, stderr io.Writer) int {
	data, err := a.loader.LoadUsageData()
	if err != nil {
		fmt.Fprintf(stderr, "Error: failed to load usage data: %v\n", err)
		return 1
	}

	if len(data) == 0 {
		fmt.Fprintln(stdout, "No usage data found.")
		return 0
	}

	sessionUsage := a.calc.AggregateSessionUsage(data, a.config.ClaudePath)
	if len(sessionUsage) == 0 {
		fmt.Fprintln(stdout, "No usage data found for the specified date range.")
		return 0
	}

	// Add total row
	total := a.calc.CalculateSessionTotal(sessionUsage)
	sessionUsage = append(sessionUsage, total)

	internalConfig := &internal.Config{
		JSON: a.config.OutputJSON,
	}

	if err := a.formatter.FormatSessionReport(sessionUsage, internalConfig); err != nil {
		fmt.Fprintf(stderr, "Error: failed to format report: %v\n", err)
		return 1
	}

	return 0
}

// printHelp prints usage information
func (a *App) printHelp(w io.Writer) {
	fmt.Fprint(w, `claude-code-usage - Claude Code Usage Analysis Tool

USAGE:
    claude-code-usage [COMMAND] [OPTIONS]

COMMANDS:
    daily    Show daily usage report (default)
    session  Show usage grouped by conversation sessions

OPTIONS:
    -since <date>     Filter from date (YYYYMMDD format)
    -until <date>     Filter until date (YYYYMMDD format)
    -path <path>      Custom path to Claude data directory (default: ~/.claude)
    -json            Output results in JSON format instead of table
    -help            Display this help message
    -version         Display version

EXAMPLES:
    claude-code-usage daily                        # Show daily report for all data
    claude-code-usage session                      # Show session-based report
    claude-code-usage daily -since 20250525        # Show daily report from May 25, 2025
    claude-code-usage session -json                # Show session report in JSON format
    claude-code-usage daily -since 20250525 -until 20250530  # Show data for date range

For more information, visit: https://github.com/landmaster135/devbox
`)
}

// printVersion prints version information
func (a *App) printVersion(w io.Writer) {
	fmt.Fprintln(w, "claude-code-usage v1.0.0")
	fmt.Fprintln(w, "A Claude Code usage analysis tool")
}

// parseDate parses date string in YYYYMMDD format
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("20060102", dateStr)
}
