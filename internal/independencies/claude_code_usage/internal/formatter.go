package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Formatter handles output formatting for different display modes
type Formatter struct{}

// NewFormatter creates a new Formatter instance
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatDailyReport formats daily usage data for display
func (f *Formatter) FormatDailyReport(data []DailyUsage, config *Config) error {
	if config.JSON {
		return f.outputJSON(data)
	}
	return f.outputDailyTable(data)
}

// FormatSessionReport formats session usage data for display
func (f *Formatter) FormatSessionReport(data []SessionUsage, config *Config) error {
	if config.JSON {
		return f.outputJSON(data)
	}
	return f.outputSessionTable(data)
}

// outputJSON outputs data in JSON format
func (f *Formatter) outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// outputDailyTable outputs daily usage in table format
func (f *Formatter) outputDailyTable(data []DailyUsage) error {
	f.printHeader("Claude Code Token Usage Report - Daily")
	
	// Calculate column widths
	dateWidth := 14
	inputWidth := 8
	outputWidth := 9
	cacheCreateWidth := 14
	cacheReadWidth := 12
	totalWidth := 14
	costWidth := 12

	// Print table header
	f.printTableBorder(dateWidth, inputWidth, outputWidth, cacheCreateWidth, cacheReadWidth, totalWidth, costWidth)
	fmt.Printf("│ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │\n",
		dateWidth-2, "Date",
		inputWidth-2, "Input",
		outputWidth-2, "Output",
		cacheCreateWidth-2, "Cache Create",
		cacheReadWidth-2, "Cache Read",
		totalWidth-2, "Total Tokens",
		costWidth-2, "Cost (USD)")
	f.printTableSeparator(dateWidth, inputWidth, outputWidth, cacheCreateWidth, cacheReadWidth, totalWidth, costWidth)

	// Print data rows
	for _, daily := range data {
		fmt.Printf("│ %-*s │ %*s │ %*s │ %*s │ %*s │ %*s │ %*s │\n",
			dateWidth-2, daily.Date,
			inputWidth-2, f.formatNumber(daily.InputTokens),
			outputWidth-2, f.formatNumber(daily.OutputTokens),
			cacheCreateWidth-2, f.formatNumber(daily.CacheCreationTokens),
			cacheReadWidth-2, f.formatNumber(daily.CacheReadTokens),
			totalWidth-2, f.formatNumber(daily.TotalTokens),
			costWidth-2, f.formatCost(daily.Cost))
	}

	f.printTableBorder(dateWidth, inputWidth, outputWidth, cacheCreateWidth, cacheReadWidth, totalWidth, costWidth)
	return nil
}

// outputSessionTable outputs session usage in table format
func (f *Formatter) outputSessionTable(data []SessionUsage) error {
	f.printHeader("Claude Code Token Usage Report - By Session")
	
	// Calculate column widths
	projectWidth := 13
	sessionWidth := 12
	inputWidth := 8
	outputWidth := 9
	cacheCreateWidth := 14
	cacheReadWidth := 12
	totalWidth := 14
	costWidth := 12
	activityWidth := 15

	// Print table header
	f.printSessionTableBorder(projectWidth, sessionWidth, inputWidth, outputWidth, cacheCreateWidth, cacheReadWidth, totalWidth, costWidth, activityWidth)
	fmt.Printf("│ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │ %-*s │\n",
		projectWidth-2, "Project",
		sessionWidth-2, "Session",
		inputWidth-2, "Input",
		outputWidth-2, "Output",
		cacheCreateWidth-2, "Cache Create",
		cacheReadWidth-2, "Cache Read",
		totalWidth-2, "Total Tokens",
		costWidth-2, "Cost (USD)",
		activityWidth-2, "Last Activity")
	f.printSessionTableSeparator(projectWidth, sessionWidth, inputWidth, outputWidth, cacheCreateWidth, cacheReadWidth, totalWidth, costWidth, activityWidth)

	// Print data rows
	for _, session := range data {
		activityStr := ""
		if !session.LastActivity.IsZero() {
			activityStr = session.LastActivity.Format("2006-01-02")
		}
		
		fmt.Printf("│ %-*s │ %-*s │ %*s │ %*s │ %*s │ %*s │ %*s │ %*s │ %-*s │\n",
			projectWidth-2, session.Project,
			sessionWidth-2, session.Session,
			inputWidth-2, f.formatNumber(session.InputTokens),
			outputWidth-2, f.formatNumber(session.OutputTokens),
			cacheCreateWidth-2, f.formatNumber(session.CacheCreationTokens),
			cacheReadWidth-2, f.formatNumber(session.CacheReadTokens),
			totalWidth-2, f.formatNumber(session.TotalTokens),
			costWidth-2, f.formatCost(session.Cost),
			activityWidth-2, activityStr)
	}

	f.printSessionTableBorder(projectWidth, sessionWidth, inputWidth, outputWidth, cacheCreateWidth, cacheReadWidth, totalWidth, costWidth, activityWidth)
	return nil
}

// printHeader prints a formatted header
func (f *Formatter) printHeader(title string) {
	width := len(title) + 4
	fmt.Printf("╭%s╮\n", strings.Repeat("─", width))
	fmt.Printf("│%*s│\n", width, "")
	fmt.Printf("│  %s  │\n", title)
	fmt.Printf("│%*s│\n", width, "")
	fmt.Printf("╰%s╯\n\n", strings.Repeat("─", width))
}

// printTableBorder prints table border for daily report
func (f *Formatter) printTableBorder(widths ...int) {
	fmt.Print("┌")
	for i, w := range widths {
		fmt.Print(strings.Repeat("─", w))
		if i < len(widths)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Println("┐")
}

// printTableSeparator prints table separator for daily report
func (f *Formatter) printTableSeparator(widths ...int) {
	fmt.Print("├")
	for i, w := range widths {
		fmt.Print(strings.Repeat("─", w))
		if i < len(widths)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤")
}

// printSessionTableBorder prints table border for session report
func (f *Formatter) printSessionTableBorder(widths ...int) {
	fmt.Print("┌")
	for i, w := range widths {
		fmt.Print(strings.Repeat("─", w))
		if i < len(widths)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Println("┐")
}

// printSessionTableSeparator prints table separator for session report
func (f *Formatter) printSessionTableSeparator(widths ...int) {
	fmt.Print("├")
	for i, w := range widths {
		fmt.Print(strings.Repeat("─", w))
		if i < len(widths)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤")
}

// formatNumber formats numbers with thousand separators
func (f *Formatter) formatNumber(n int) string {
	s := strconv.Itoa(n)
	if len(s) <= 3 {
		return s
	}
	
	var result strings.Builder
	for i, char := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(char)
	}
	return result.String()
}

// formatCost formats cost as currency
func (f *Formatter) formatCost(cost float64) string {
	return fmt.Sprintf("$%.2f", cost)
}
