package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

func GetTicketNumber(commit string) (string, error) {
	ticketNumber := regexp.MustCompile(`^[A-Z_]+-[0-9]+`)
	ret := ticketNumber.Find([]byte(commit))
	if ret != nil {
		return string(ret), nil
	}
	return "", fmt.Errorf("There is no a ticket number in line %s", commit)
}

func SplitLog(log string) []string {
	re := regexp.MustCompile(`\n`)
	split := re.Split(log, -1)
	var set []string
	for i := range split {
		set = append(set, split[i])
	}
	return set
}

func TimeFromHGTimestamp(timestamp string) (time.Time, error) {
	reTs := regexp.MustCompile(`[0-9]+`)
	ts, err := strconv.ParseInt(string(reTs.Find([]byte(timestamp))), 10, 64)
	if err != nil {
		return time.Unix(0, 0), err
	}
	tm := time.Unix(ts, 0)
	return tm, nil
}

func SplitLogLine(logLine string) (string, string, error) {
	re := regexp.MustCompile(`\|`)
	parsed := re.Split(logLine, -1)
	if len(parsed) != 2 {
		return "", "", fmt.Errorf("Wrong log format for entry %s", logLine)
	}
	timestamp := parsed[0]
	commit := parsed[1]
	return timestamp, commit, nil
}
