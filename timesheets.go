package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"sort"
)

func main() {
	f, err := os.Open("config.yaml")
	if err != nil {
		fmt.Errorf("Can't open config file #{err}")
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Errorf("Can't decode config %v", err)
	}

	var lines []Line
	repos := cfg.RepoList()
	lc := make(chan []Line, len(repos)+len(cfg.Calendars.Calendar))
	for _, repo := range repos {
		go GetLines(cfg.Repos.User, repo, lc)
	}
	for _, iCal := range cfg.Calendars.Calendar {
		go GetCalendarLines(iCal, cfg.Calendars.DefaultTicket, lc)
	}
	for i := 0; i < len(repos)+len(cfg.Calendars.Calendar); i++ {
		cLines := <-lc
		lines = append(lines, cLines...)
	}
	sort.SliceStable(lines, func(i, j int) bool {
		return lines[i].start.Before(lines[j].start)
	})
	lines = GroupLines(lines)
	for _, l := range lines {
		println(l.Export())
	}
}
