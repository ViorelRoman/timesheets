package main

import (
	"fmt"
	"time"
)

type Line struct {
	start       time.Time
	end         time.Time
	ticket      string
	description string
}

func (l *Line) ValidateLine(lines []Line) bool {
	for _, line := range lines {
		if (l.start.After(line.start) || l.start.Equal(line.start)) && l.start.Before(line.end) {
			println("Starting is inside other time sequence")
			return false
		}
		if l.end.After(line.start) && l.end.Before(line.end) {
			println("Ending is inside other time sequence")
			return false
		}
		if (line.start.After(l.start) || line.start.Equal(l.start)) && line.start.Before(l.end) {
			println("Sequence is inside other time sequence")
			return false
		}
		if line.end.After(l.start) && line.end.Before(l.end) {
			println("Other time sequence is inside this sequence")
			return false
		}
	}
	return true
}

func (l *Line) Export() string {
	location, _ := time.LoadLocation("Europe/Moscow")
	ret := fmt.Sprintf("%s,%s,%s,%s,%s", l.start.In(location).Format("2006-01-02"), l.start.In(location).Format("15:04"), l.end.In(location).Format("15:04"), l.ticket, l.description)
	return ret
}

func CreateLine(logline string) (Line, error) {
	line := Line{}
	timestamp, commit, err := SplitLogLine(logline)
	if err != nil {
		return line, err
	}
	ticket, err := GetTicketNumber(commit)
	if err != nil {
		return line, err
	}
	end_ts, err := TimeFromHGTimestamp(timestamp)
	if err != nil {
		return line, err
	}
	line.end = end_ts
	line.start = line.end
	line.ticket = ticket
	line.description = commit[len(ticket)+1:]
	return line, nil
}

func GroupLines(lines []Line) []Line {
	var ret []Line
	prev := lines[0]
	for i, line := range lines {
		if line.ticket == prev.ticket && line.description == prev.description {
			prev.end = line.end
		} else {
			ret = append(ret, prev)
			prev = line
		}
		if i == len(lines)-1 {
			ret = append(ret, prev)
		}
	}
	return ret
}

func GetLines(user string, repo Repo, out chan []Line) {
	var lines []Line
	log, err := repo.GetLog(user)
	if err != nil {
		out <- lines
		return
	}
	for _, logLine := range SplitLog(log) {
		line, err := CreateLine(logLine)
		if err == nil {
			lines = append(lines, line)
		}
	}
	out <- lines
}
