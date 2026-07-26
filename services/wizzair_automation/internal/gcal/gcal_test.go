package gcal_test

import (
	"testing"
	"time"

	"avm99963-monorepo/services/wizzair_automation/internal/gcal"
	"avm99963-monorepo/services/wizzair_automation/internal/parser"
)

func TestBuildFlightEvent(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		t.Fatalf("failed to load Europe/Warsaw location: %v", err)
	}

	depTime := time.Date(2027, 1, 4, 5, 40, 0, 0, loc)
	arrTime := time.Date(2027, 1, 4, 8, 45, 0, 0, loc)

	seg := parser.FlightSegment{
		PNR:           "RESCOD",
		FlightNumber:  "W6 1475",
		Origin:        "WAW",
		Destination:   "BCN",
		DepartureTime: depTime,
		ArrivalTime:   arrTime,
	}

	t.Run("constructs_event_with_correct_summary_location_and_iCalUID", func(t *testing.T) {
		event := gcal.BuildFlightEvent(seg)

		expectedSummary := "🛫 WAW → BCN"
		if event.Summary != expectedSummary {
			t.Errorf("event.Summary = %q; want %q", event.Summary, expectedSummary)
		}

		expectedLocation := "WAW Airport"
		if event.Location != expectedLocation {
			t.Errorf("event.Location = %q; want %q", event.Location, expectedLocation)
		}

		expectedICalUID := "wizz-flight-RESCOD-W6 1475"
		if event.ICalUID != expectedICalUID {
			t.Errorf("event.ICalUID = %q; want %q", event.ICalUID, expectedICalUID)
		}
	})
}

func TestSyncFlightEvents_DryRun(t *testing.T) {
	seg := parser.FlightSegment{
		PNR:           "RESCOD",
		FlightNumber:  "W6 1475",
		Origin:        "WAW",
		Destination:   "BCN",
		DepartureTime: time.Now(),
		ArrivalTime:   time.Now().Add(2 * time.Hour),
	}

	t.Run("dry_run_mode_logs_and_skips_network_calls", func(t *testing.T) {
		client := gcal.NewClient("primary", "", true) // dryRun = true
		err := client.SyncFlightEvents([]parser.FlightSegment{seg})
		if err != nil {
			t.Errorf("expected no error in dry run mode, got: %v", err)
		}
	})
}
