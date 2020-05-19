package main

import (
	"testing"
	"time"
)

func TestTimeSequenceInOtherTimeSequenceValidation(t *testing.T) {
	line1 := Line{
		time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 13, 0, 0, 0, time.UTC),
		"BROADSIGN-2020",
		"first activity",
	}
	line2 := Line{
		time.Date(2020, time.April, 1, 11, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 14, 0, 0, 0, time.UTC),
		"BROADSIGN-2021",
		"second activity",
	}
	lines := []Line{line1}
	got := line2.ValidateLine(lines)
	if got != false {
		t.Error("Line 2 is actually invalid")
	}
}

func TestTimeSequenceInOtherTimeSequenceReverceValidation(t *testing.T) {
	line1 := Line{
		time.Date(2020, time.April, 1, 11, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 14, 0, 0, 0, time.UTC),
		"BROADSIGN-2021",
		"second activity",
	}
	line2 := Line{
		time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 13, 0, 0, 0, time.UTC),
		"BROADSIGN-2020",
		"first activity",
	}
	lines := []Line{line1}
	got := line2.ValidateLine(lines)
	if got != false {
		t.Error("Line 2 is actually invalid")
	}
}

func TestTimeEndIntersectsTimeSequenceValidation(t *testing.T) {
	line1 := Line{
		time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 13, 0, 0, 0, time.UTC),
		"BROADSIGN-2020",
		"first activity",
	}
	line2 := Line{
		time.Date(2020, time.April, 1, 11, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 12, 30, 0, 0, time.UTC),
		"BROADSIGN-2021",
		"second activity",
	}
	lines := []Line{line1}
	got := line2.ValidateLine(lines)
	if got != false {
		t.Error("Line 2 is actually invalid")
	}
}

func TestTimeEndEqualsTimeStartValidation(t *testing.T) {
	line1 := Line{
		time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 13, 0, 0, 0, time.UTC),
		"BROADSIGN-2020",
		"first activity",
	}
	line2 := Line{
		time.Date(2020, time.April, 1, 13, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 14, 30, 0, 0, time.UTC),
		"BROADSIGN-2021",
		"second activity",
	}
	lines := []Line{line1}
	got := line2.ValidateLine(lines)
	if got == false {
		t.Error("Line 2 is actually valid")
	}
}

func TestTimeEndEqualsTimeStartReverseValidation(t *testing.T) {
	line1 := Line{
		time.Date(2020, time.April, 1, 13, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 14, 30, 0, 0, time.UTC),
		"BROADSIGN-2021",
		"second activity",
	}
	line2 := Line{
		time.Date(2020, time.April, 1, 12, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 13, 0, 0, 0, time.UTC),
		"BROADSIGN-2020",
		"first activity",
	}
	lines := []Line{line1}
	got := line2.ValidateLine(lines)
	if got == false {
		t.Error("Line 2 is actually valid")
	}
}

func TestExportLine(t *testing.T) {
	line := Line{
		time.Date(2020, time.April, 1, 13, 0, 0, 0, time.UTC),
		time.Date(2020, time.April, 1, 14, 30, 0, 0, time.UTC),
		"BROADSIGN-2021",
		"second activity",
	}
	got := line.Export()
	if got != "2020-04-01,16:00,17:30,BROADSIGN-2021,second activity" {
		t.Error("Export format is wrong", got)
	}
}

func TestFindTicketNumber(t *testing.T) {
	commit := "BROADSIGN-2020 some work done"
	got, err := GetTicketNumber(commit)
	if err != nil {
		t.Error(err)
	}
	if got != "BROADSIGN-2020" {
		t.Error("Ticket number got wrong", got)
	}
}

func TestFindTicketNumberWithUnderscore(t *testing.T) {
	commit := "GRID_ADC-863 some work done"
	got, err := GetTicketNumber(commit)
	if err != nil {
		t.Error(err)
	}
	if got != "GRID_ADC-863" {
		t.Error("Ticket number got wrong", got)
	}
}

func TestSplitHGLog(t *testing.T) {
	log := `1589393622|BROADSIGN-2548 deploy 3.20.1 on QA
1589378723|BROADSIGN-2548 deploy 3.20.1 on staging`
	got := SplitLog(log)
	if len(got) != 2 {
		t.Error("Split works wrong", got)
	}
}

func TestGetHGlog(t *testing.T) {
	user := "vroman"
	repo := HGRepo{"/home/vroman/iponweb/broadsign-ci-test/project"}
	log, err := repo.GetLog(user)
	if err != nil {
		t.Error(err)
	}
	got := SplitLog(log)
	if len(got) < 2 {
		t.Error("Wrong formated log", got, log)
	}
}

func TestParseTimestamp(t *testing.T) {
	timestamp := "1588090089.0-10800"
	got, err := TimeFromHGTimestamp(timestamp)
	if err != nil {
		t.Error(err)
	}
	location, _ := time.LoadLocation("Europe/Moscow")
	expected := time.Date(2020, 4, 28, 19, 8, 9, 0, location)
	if !(got.Equal(expected)) {
		t.Error("Wrong time parsed", got, expected)
	}
}

func TestParsetTimestampWithoutOffset(t *testing.T) {
	timestamp := "1589393622"
	got, err := TimeFromHGTimestamp(timestamp)
	if err != nil {
		t.Error(err)
	}
	location, _ := time.LoadLocation("Europe/Moscow")
	expected := time.Date(2020, 5, 13, 21, 13, 42, 0, location)
	if !(got.Equal(expected)) {
		t.Error("Wrong time parsed", got, expected)
	}
}

func TestParseLineFromLogEntry(t *testing.T) {
	logline := "1588090089.0-10800|BROADSIGN-2546 absolute import"
	location, _ := time.LoadLocation("Europe/Moscow")
	expected := Line{
		time.Date(2020, time.April, 28, 19, 8, 9, 0, location),
		time.Date(2020, time.April, 28, 19, 8, 9, 0, location),
		"BROADSIGN-2546",
		"absolute import",
	}
	got, err := CreateLine(logline)
	if err != nil {
		t.Error(err)
	}
	if !(got.start.Equal(expected.start)) {
		t.Error("Start date parsed wrong", got.start)
	}
	if !(got.end.Equal(expected.end)) {
		t.Error("End date parsed wrong", got.end)
	}
	if got.ticket != expected.ticket {
		t.Error("Ticket number parsed wrong", got.ticket)
	}
	if got.description != expected.description {
		t.Error("Comment parsed wrong", got.description)
	}
}

func TestNotElegibleCommit(t *testing.T) {
	logline := "1588011842.0-10800|Added tag 3.19.0-auto-approve-21270420 for changeset f9f1f72203a5"
	_, err := CreateLine(logline)
	if err == nil {
		t.Error("Parsed a not timesheetable line")
	}
}

func TestGroupLines(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Moscow")
	lines := []Line{
		{
			time.Date(2020, time.April, 28, 19, 8, 9, 0, location),
			time.Date(2020, time.April, 28, 19, 8, 9, 0, location),
			"BROADSIGN-2546",
			"absolute import",
		},
		{
			time.Date(2020, time.April, 28, 20, 8, 9, 0, location),
			time.Date(2020, time.April, 28, 20, 8, 9, 0, location),
			"BROADSIGN-2546",
			"absolute import",
		},
		{
			time.Date(2020, time.April, 28, 21, 8, 9, 0, location),
			time.Date(2020, time.April, 28, 21, 8, 9, 0, location),
			"BROADSIGN-2546",
			"absolute import 1",
		},
		{
			time.Date(2020, time.April, 28, 22, 8, 9, 0, location),
			time.Date(2020, time.April, 28, 22, 8, 9, 0, location),
			"BROADSIGN-2546",
			"absolute import 1",
		},
		{
			time.Date(2020, time.April, 28, 23, 8, 9, 0, location),
			time.Date(2020, time.April, 28, 23, 8, 9, 0, location),
			"BROADSIGN-2548",
			"absolute import 1",
		},
		{
			time.Date(2020, time.April, 29, 1, 8, 9, 0, location),
			time.Date(2020, time.April, 29, 1, 8, 9, 0, location),
			"BROADSIGN-2547",
			"absolute import",
		},
		{
			time.Date(2020, time.April, 29, 2, 8, 9, 0, location),
			time.Date(2020, time.April, 29, 2, 8, 9, 0, location),
			"BROADSIGN-2547",
			"absolute import",
		},
		{
			time.Date(2020, time.April, 29, 3, 8, 9, 0, location),
			time.Date(2020, time.April, 29, 3, 8, 9, 0, location),
			"BROADSIGN-2549",
			"single commit at the end",
		},
	}
	got := GroupLines(lines)
	if len(got) != 5 {
		t.Error("Lines where not grouped", len(got))
	}
	if !(got[0].start.Equal(lines[0].start)) {
		t.Error("Start time is wrong #1, got", got[0].start)
	}
	if !(got[0].end.Equal(lines[1].end)) {
		t.Error("End time is wrong #1, got", got[0].end)
	}
	if !(got[1].start.Equal(lines[2].start)) {
		t.Error("Start time is wrong #2, got", got[1].start)
	}
	if !(got[1].end.Equal(lines[3].end)) {
		t.Error("End time is wrong #2, got", got[1].end)
	}
	if !(got[2].start.Equal(lines[4].start)) {
		t.Error("Start time is wrong #3, got", got[2].start)
	}
	if !(got[2].end.Equal(lines[4].end)) {
		t.Error("End time is wrong #3, got", got[2].end)
	}
	if !(got[3].start.Equal(lines[5].start)) {
		t.Error("Start time is wrong #4, got", got[3].start)
	}
	if !(got[3].end.Equal(lines[6].end)) {
		t.Error("End time is wrong #4, got", got[3].end)
	}
	if !(got[4].start.Equal(lines[7].start)) {
		t.Error("Start time is wrong #4, got", got[4].start)
	}
	if !(got[4].end.Equal(lines[7].end)) {
		t.Error("End time is wrong #4, got", got[4].end)
	}
}
