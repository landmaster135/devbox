package common

import (
	"regexp"
	"strconv"
	"time"
)

var webClipFilePattern = regexp.MustCompile(`^web-summary-(\d{8})-(\d{6})-([A-Za-z0-9][A-Za-z0-9_-]*)\.md$`)
var movieClipFilePattern = regexp.MustCompile(`^movie-summary-(\d{8})-(\d{6})-([A-Za-z0-9][A-Za-z0-9_-]*)\.md$`)

var webClipAttachmentFilePattern = regexp.MustCompile(`^(web-summary-\d{8}-\d{6}-[A-Za-z0-9][A-Za-z0-9_-]*)_(\d+)\.([A-Za-z0-9]+)$`)
var movieClipAttachmentFilePattern = regexp.MustCompile(`^(movie-summary-\d{8}-\d{6}-[A-Za-z0-9][A-Za-z0-9_-]*)_(\d+)\.([A-Za-z0-9]+)$`)
var commonMemoFilePattern = regexp.MustCompile(`^(\d{14})_(\d+)\.md$`)
var commonMemoAttachmentFilePattern = regexp.MustCompile(`^(\d{14})_(\d+)_(\d+)\.([A-Za-z0-9]+)$`)

var jstLocation = time.FixedZone("JST", 9*60*60)

func MatchWebClipFile(baseName string) bool {
	return webClipFilePattern.MatchString(baseName)
}

func MatchMovieClipFile(baseName string) bool {
	return movieClipFilePattern.MatchString(baseName)
}

func ParseWebClipDisplayTime(baseName string) (string, bool) {
	matches := webClipFilePattern.FindStringSubmatch(baseName)
	if len(matches) != 4 {
		return "", false
	}
	return matches[1] + matches[2], true
}

func ParseMovieClipDisplayTime(baseName string) (string, bool) {
	matches := movieClipFilePattern.FindStringSubmatch(baseName)
	if len(matches) != 4 {
		return "", false
	}
	return matches[1] + matches[2], true
}

func ParseWebAttachmentContentBaseName(baseName string) (string, bool) {
	matches := webClipAttachmentFilePattern.FindStringSubmatch(baseName)
	if len(matches) != 4 {
		return "", false
	}
	return matches[1], true
}

func ParseMovieAttachmentContentBaseName(baseName string) (string, bool) {
	matches := movieClipAttachmentFilePattern.FindStringSubmatch(baseName)
	if len(matches) != 4 {
		return "", false
	}
	return matches[1], true
}

func MatchCommonMemoFile(baseName string) bool {
	return commonMemoFilePattern.MatchString(baseName)
}

func ParseCommonMemoDisplayTime(baseName string) (string, bool) {
	matches := commonMemoFilePattern.FindStringSubmatch(baseName)
	if len(matches) != 3 {
		return "", false
	}
	return matches[1], true
}

func ParseCommonMemoFileKey(baseName string) (string, int, bool) {
	matches := commonMemoFilePattern.FindStringSubmatch(baseName)
	if len(matches) != 3 {
		return "", 0, false
	}

	number, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, false
	}
	return matches[1], number, true
}

func ParseCommonMemoAttachmentKey(baseName string) (string, int, bool) {
	matches := commonMemoAttachmentFilePattern.FindStringSubmatch(baseName)
	if len(matches) != 5 {
		return "", 0, false
	}

	number, err := strconv.Atoi(matches[2])
	if err != nil {
		return "", 0, false
	}
	return matches[1], number, true
}
