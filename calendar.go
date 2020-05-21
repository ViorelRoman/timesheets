package main

import (
	"github.com/apognu/gocal"
	"net/http"
	"time"
)

func GetCalendarLines(ical string, defaultTicket string, out chan []Line) {
	var lines []Line
	resp, err := http.Get(ical)
	if err != nil {
		out <- lines
		return
	}
	defer resp.Body.Close()
	start, end := time.Now().AddDate(0, -1, 0), time.Now()
	c := gocal.NewParser(resp.Body)
	c.Start, c.End = &start, &end
	c.Parse()
	for _, e := range c.Events {
		line := Line{
			*e.Start,
			*e.End,
			defaultTicket,
			e.Summary,
		}
		lines = append(lines, line)
	}
	out <- lines
}
