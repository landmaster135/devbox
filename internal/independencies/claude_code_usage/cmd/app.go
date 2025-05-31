package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/landmaster135/devbox/internal/independencies/claude_code_usage/internal"
)

// App represents the main application
type App struct {
	config    *internal.Config
	loader    *internal.DataLoader
	calc      *internal.Calculator
	formatter *internal.Formatter
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{
		calc:      internal.NewCalculator(),
		formatter: internal.NewFormatter(),
	}
}

// Run executes the application
func (a *App) Run(args []string) error {
	if err := a.parseFlags(args); err != nil {
		return err
	}

	a.loader = internal.NewDataLoader(a.config)

	// Default to daily command if no command specified
	command := "daily"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		command = args[0]
		args = args[1:]
	}

	switch command {
	case "daily", "":
		return a.runDailyCommand()
	case "session":
		return a.runSessionCommand()
	case "help", "-h", "--help":
		a.printHelp()
		return nil
	case "version", "-v", "--version":
		a.printVersion()
		return nil
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// runDailyCommand executes the daily report command
func (a *App) runDailyCommand() error {
	data, err := a.loader.LoadUsageData()
	if err != nil {
		return fmt.Errorf("failed to load usage data: %w", err)
	}

	if len(data) == 0 {
		fmt.Println("No usage data found.")
		return nil
	}

	dailyUsage := a.calc.AggregateDailyUsage(data)
	if len(dailyUsage) == 0 {
		fmt.Println("No usage data found for the specified date range.")
		return nil
	}

	// Add total row
	total := a.calc.CalculateDailyTotal(dailyUsage)
	dailyUsage = append(dailyUsage, total)

	return a.formatter.FormatDailyReport(dailyUsage, a.config)
}

// runSessionCommand executes the session report command
func (a *App) runSessionCommand() error {
	data, err := a.loader.LoadUsageData()
	if err != nil {
		return fmt.Errorf("failed to load usage data: %w", err)
	}

	if len(data) == 0 {
		fmt.Println("No usage data found.")
		return nil
	}

	sessionUsage := a.calc.AggregateSessionUsage(data, a.config.ClaudePath)
	if len(sessionUsage) == 0 {
		fmt.Println("No usage data found for the specified date range.")
		return nil
	}

	// Add total row
	total := a.calc.CalculateSessionTotal(sessionUsage)
	sessionUsage = append(sessionUsage, total)

	return a.formatter.FormatSessionReport(sessionUsage, a.config)
}

// parseFlags parses command line flags
func (a *App) parseFlags(args []string) error {
	// Create flag set
	fs := flag.NewFlagSet("ccusage", flag.ContinueOnError)
	fs.Usage = func() {
		a.printHelp()
	}

	var (
		sinceStr   = fs.String("since", "", "Filter from date (YYYYMMDD format)")
		untilStr   = fs.String("until", "", "Filter until date (YYYYMMDD format)")
		claudePath = fs.String("path", "", "Custom path to Claude data directory")
		jsonOutput = fs.Bool("json", false, "Output results in JSON format")
		help       = fs.Bool("help", false, "Display help message")
		version    = fs.Bool("version", false, "Display version")
	)

	// Allow short flags
	fs.StringVar(sinceStr, "s", "", "Filter from date (YYYYMMDD format)")
	fs.StringVar(untilStr, "u", "", "Filter until date (YYYYMMDD format)")
	fs.StringVar(claudePath, "p", "", "Custom path to Claude data directory")
	fs.BoolVar(jsonOutput, "j", false, "Output results in JSON format")
	fs.BoolVar(help, "h", false, "Display help message")
	fs.BoolVar(version, "v", false, "Display version")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *help {
		a.printHelp()
		os.Exit(0)
	}

	if *version {
		a.printVersion()
		os.Exit(0)
	}

	// Initialize config
	a.config = &internal.Config{
		JSON: *jsonOutput,
	}

	// Set Claude path
	if *claudePath != "" {
		a.config.ClaudePath = *claudePath
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		a.config.ClaudePath = filepath.Join(homeDir, ".claude")
	}

	// Parse date filters
	if *sinceStr != "" {
		since, err := parseDate(*sinceStr)
		if err != nil {
			return fmt.Errorf("invalid since date: %w", err)
		}
		a.config.Since = &since
	}

	if *untilStr != "" {
		until, err := parseDate(*untilStr)
		if err != nil {
			return fmt.Errorf("invalid until date: %w", err)
		}
		a.config.Until = &until
	}

	return nil
}

// parseDate parses date string in YYYYMMDD format
func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("20060102", dateStr)
}

// printHelp prints usage information
func (a *App) printHelp() {
	fmt.Print(`ccusage - Claude Code Usage Analysis Tool

USAGE:
    ccusage [COMMAND] [OPTIONS]

COMMANDS:
    daily    Show daily usage report (default)
    session  Show usage grouped by conversation sessions
    help     Show this help message
    version  Show version information

OPTIONS:
    -s, --since <date>    Filter from date (YYYYMMDD format)
    -u, --until <date>    Filter until date (YYYYMMDD format)
    -p, --path <path>     Custom path to Claude data directory (default: ~/.claude)
    -j, --json           Output results in JSON format instead of table
    -h, --help           Display this help message
    -v, --version        Display version

EXAMPLES:
    ccusage                                    # Show daily report for all data
    ccusage daily                             # Same as above
    ccusage session                           # Show session-based report
    ccusage daily --since 20250525           # Show daily report from May 25, 2025
    ccusage session --json                   # Show session report in JSON format
    ccusage daily --since 20250525 --until 20250530  # Show data for date range

For more information, visit: https://github.com/your-username/ccusage-go
`)
}

// printVersion prints version information
func (a *App) printVersion() {
	fmt.Println("ccusage-go v1.0.0")
	fmt.Println("A Claude Code usage analysis tool written in Go")
}
