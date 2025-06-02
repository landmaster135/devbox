package usecases

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// UnixToISO8601 converts a Unix timestamp to ISO-8601 format
func UnixToISO8601(unixTimestamp string) (string, error) {
	// Parse the Unix timestamp
	timestamp, err := strconv.ParseInt(unixTimestamp, 10, 64)
	if err != nil {
		return "", errors.New("invalid Unix timestamp: " + err.Error())
	}

	// Convert to time.Time
	t := time.Unix(timestamp, 0)

	// Format as ISO-8601
	return t.Format(time.RFC3339), nil
}

// ISO8601ToUnix converts an ISO-8601 formatted time to Unix timestamp
func ISO8601ToUnix(iso8601 string) (string, error) {
	// Parse the ISO-8601 timestamp
	t, err := time.Parse(time.RFC3339, iso8601)
	if err != nil {
		return "", errors.New("invalid ISO-8601 format: " + err.Error())
	}

	// Convert to Unix timestamp
	unixTimestamp := t.Unix()

	// Convert to string
	return strconv.FormatInt(unixTimestamp, 10), nil
}

// DateToUnix converts a date string to Unix timestamp
// If isJST is true, the time is set to 00:00:00 JST
// If isJST is false, the time is set to 00:00:00 UTC
func DateToUnix(dateStr string, isJST bool) (string, error) {
	// Check if the date string already contains time information
	if strings.Contains(dateStr, "T") {
		return "", errors.New("date string already contains time information, use ISO8601ToUnix instead")
	}

	var loc *time.Location
	// var err error

	// Set the location based on isJST
	loc = time.UTC
	if isJST {
		loc = time.FixedZone("UTC+9", 9*60*60)
	}

	// Try different date formats
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"20060102",
	}

	var t time.Time
	var parseErr error

	for _, format := range formats {
		t, parseErr = time.ParseInLocation(format, dateStr, loc)
		if parseErr == nil {
			break
		}
	}

	if parseErr != nil {
		return "", errors.New("invalid date format: " + parseErr.Error())
	}

	// Convert to Unix timestamp
	unixTimestamp := t.Unix()

	// Convert to string
	return strconv.FormatInt(unixTimestamp, 10), nil
}

func NowToUnix() string {
	currentTime := time.Now()
	unixTimestamp := currentTime.Unix()
	return strconv.FormatInt(unixTimestamp, 10)
}

func NowToISO8601InUTC() string {
	currentTime := time.Now()
	t := currentTime.UTC().Format(time.RFC3339)
	return t
}

func NowToISO8601InJST() string {
	currentTime := time.Now()
	t := currentTime.In(time.FixedZone("JST", 9*60*60)).Format(time.RFC3339)
	return t
}
