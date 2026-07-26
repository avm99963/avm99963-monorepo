package caldav_test

import (
	"strings"
	"testing"
	"time"

	"avm99963-monorepo/services/wizzair_automation/internal/caldav"
	"avm99963-monorepo/services/wizzair_automation/internal/parser"
	"avm99963-monorepo/services/wizzair_automation/internal/schedule"
)

func TestGenerateVTODO(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		t.Fatalf("failed to load location: %v", err)
	}

	seg := parser.FlightSegment{
		PNR:          "RESCOD",
		FlightNumber: "W6 1475",
		Origin:       "WAW",
		Destination:  "BCN",
	}

	taskStart := time.Date(2027, 1, 4, 8, 30, 0, 0, loc)
	taskDue := time.Date(2027, 1, 4, 10, 30, 0, 0, loc)
	task := schedule.CheckInTask{
		Start: taskStart,
		Due:   taskDue,
	}

	t.Run("formats_rfc5545_vtodo_with_summary_dtstart_due_and_both_start_and_due_valarms", func(t *testing.T) {
		ics := caldav.GenerateVTODO(seg, task)

		if !strings.Contains(ics, "BEGIN:VCALENDAR") {
			t.Error("expected VCALENDAR header")
		}
		if !strings.Contains(ics, "BEGIN:VTODO") {
			t.Error("expected VTODO header")
		}
		if !strings.Contains(ics, "SUMMARY:Check-in for flight W6 1475 (WAW → BCN)") {
			t.Errorf("unexpected SUMMARY in ics output:\n%s", ics)
		}
		expectedDTStart := "DTSTART:" + taskStart.UTC().Format("20060102T150405Z")
		if !strings.Contains(ics, expectedDTStart) {
			t.Errorf("expected %q in ics output, got:\n%s", expectedDTStart, ics)
		}
		expectedDTDue := "DUE:" + taskDue.UTC().Format("20060102T150405Z")
		if !strings.Contains(ics, expectedDTDue) {
			t.Errorf("expected %q in ics output, got:\n%s", expectedDTDue, ics)
		}
		if !strings.Contains(ics, "TRIGGER;RELATED=START:PT0M") {
			t.Error("expected start-time VALARM trigger")
		}
		if !strings.Contains(ics, "TRIGGER;RELATED=END:PT0M") {
			t.Error("expected due-time VALARM trigger")
		}
		if !strings.Contains(ics, "UID:checkin-RESCOD-W6_1475") {
			t.Errorf("unexpected UID in ics output:\n%s", ics)
		}
	})
}

func TestPushCheckInTask_DryRun(t *testing.T) {
	seg := parser.FlightSegment{
		PNR:          "RESCOD",
		FlightNumber: "W6 1475",
		Origin:       "WAW",
		Destination:  "BCN",
	}
	task := schedule.CheckInTask{
		Start: time.Now(),
		Due:   time.Now().Add(2 * time.Hour),
	}

	t.Run("dry_run_mode_does_not_issue_http_put_requests", func(t *testing.T) {
		client := caldav.NewClient("https://nextcloud.example.com/dav/calendars/user/personal/", "user", "pass", true)
		err := client.PushCheckInTask(seg, task)
		if err != nil {
			t.Errorf("expected no error in dry run mode, got: %v", err)
		}
	})
}
