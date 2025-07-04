package internal

import (
	"sort"
)

// Calculator handles data aggregation and cost calculations
type Calculator struct{}

// NewCalculator creates a new Calculator instance
func NewCalculator() *Calculator {
	return &Calculator{}
}

// AggregateDailyUsage aggregates usage data by date
func (c *Calculator) AggregateDailyUsage(data []UsageData) []DailyUsage {
	dailyMap := make(map[string]*DailyUsage)

	for _, entry := range data {
		dateKey := entry.Timestamp.Format("2006-01-02")

		if daily, exists := dailyMap[dateKey]; exists {
			daily.InputTokens += entry.InputTokens
			daily.OutputTokens += entry.OutputTokens
			daily.CacheCreationTokens += entry.CacheCreationTokens
			daily.CacheReadTokens += entry.CacheReadTokens
			daily.Cost += entry.Cost
		} else {
			dailyMap[dateKey] = &DailyUsage{
				Date:                dateKey,
				InputTokens:         entry.InputTokens,
				OutputTokens:        entry.OutputTokens,
				CacheCreationTokens: entry.CacheCreationTokens,
				CacheReadTokens:     entry.CacheReadTokens,
				Cost:                entry.Cost,
			}
		}
	}

	// Convert map to slice and calculate total tokens
	var result []DailyUsage
	for _, daily := range dailyMap {
		daily.TotalTokens = daily.InputTokens + daily.OutputTokens +
			daily.CacheCreationTokens + daily.CacheReadTokens
		result = append(result, *daily)
	}

	// Sort by date
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result
}

// AggregateSessionUsage aggregates usage data by session
func (c *Calculator) AggregateSessionUsage(data []UsageData, claudePath string) []SessionUsage {
	sessionMap := make(map[string]*SessionUsage)

	for _, entry := range data {
		// For this implementation, we'll use a simple session key
		// In a real implementation, you'd parse the file path to get project/session
		sessionKey := entry.Timestamp.Format("2006-01-02") // Simplified for demo

		if session, exists := sessionMap[sessionKey]; exists {
			session.InputTokens += entry.InputTokens
			session.OutputTokens += entry.OutputTokens
			session.CacheCreationTokens += entry.CacheCreationTokens
			session.CacheReadTokens += entry.CacheReadTokens
			session.Cost += entry.Cost

			// Update last activity if this entry is more recent
			if entry.Timestamp.After(session.LastActivity) {
				session.LastActivity = entry.Timestamp
			}
		} else {
			sessionMap[sessionKey] = &SessionUsage{
				Project:             "project", // Would be extracted from file path
				Session:             sessionKey,
				InputTokens:         entry.InputTokens,
				OutputTokens:        entry.OutputTokens,
				CacheCreationTokens: entry.CacheCreationTokens,
				CacheReadTokens:     entry.CacheReadTokens,
				Cost:                entry.Cost,
				LastActivity:        entry.Timestamp,
			}
		}
	}

	// Convert map to slice and calculate total tokens
	var result []SessionUsage
	for _, session := range sessionMap {
		session.TotalTokens = session.InputTokens + session.OutputTokens +
			session.CacheCreationTokens + session.CacheReadTokens
		result = append(result, *session)
	}

	// Sort by cost (descending)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Cost > result[j].Cost
	})

	return result
}

// CalculateTotalUsage calculates total usage across all data
func (c *Calculator) CalculateTotalUsage(data []UsageData) (inputTokens, outputTokens, cacheCreationTokens, cacheReadTokens int, totalCost float64) {
	for _, entry := range data {
		inputTokens += entry.InputTokens
		outputTokens += entry.OutputTokens
		cacheCreationTokens += entry.CacheCreationTokens
		cacheReadTokens += entry.CacheReadTokens
		totalCost += entry.Cost
	}
	return
}

// CalculateDailyTotal calculates total usage for daily reports
func (c *Calculator) CalculateDailyTotal(dailyUsage []DailyUsage) DailyUsage {
	var total DailyUsage
	total.Date = "Total"

	for _, daily := range dailyUsage {
		total.InputTokens += daily.InputTokens
		total.OutputTokens += daily.OutputTokens
		total.CacheCreationTokens += daily.CacheCreationTokens
		total.CacheReadTokens += daily.CacheReadTokens
		total.TotalTokens += daily.TotalTokens
		total.Cost += daily.Cost
	}

	return total
}

// CalculateSessionTotal calculates total usage for session reports
func (c *Calculator) CalculateSessionTotal(sessionUsage []SessionUsage) SessionUsage {
	var total SessionUsage
	total.Project = "Total"
	total.Session = ""

	for _, session := range sessionUsage {
		total.InputTokens += session.InputTokens
		total.OutputTokens += session.OutputTokens
		total.CacheCreationTokens += session.CacheCreationTokens
		total.CacheReadTokens += session.CacheReadTokens
		total.TotalTokens += session.TotalTokens
		total.Cost += session.Cost
	}

	return total
}
