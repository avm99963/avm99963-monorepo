package schedule_test

import (
	"testing"
	"time"

	"avm99963-monorepo/services/wizzair_automation/internal/schedule"
)

func TestCalculateCheckInTask(t *testing.T) {
	loc, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		t.Fatalf("failed to load Europe/Warsaw location: %v", err)
	}

	tests := []struct {
		name             string
		departureTimeStr string
		expectedStartStr string
		expectedDueStr   string
	}{
		{
			name:             "daytime_checkin_opens_at_14_00_taskStart_is_14_00_taskDue_is_16_00",
			departureTimeStr: "2026-09-05 14:00",
			expectedStartStr: "2026-09-04 14:00",
			expectedDueStr:   "2026-09-04 16:00",
		},
		{
			name:             "overnight_checkin_opens_at_03_00_in_rest_period_adjusts_taskStart_to_08_30_and_taskDue_to_10_30",
			departureTimeStr: "2026-09-05 03:00",
			expectedStartStr: "2026-09-04 08:30",
			expectedDueStr:   "2026-09-04 10:30",
		},
		{
			name:             "taskStart_at_23_00_adjusts_taskDue_to_midnight",
			departureTimeStr: "2026-09-05 23:00",
			expectedStartStr: "2026-09-04 23:00",
			expectedDueStr:   "2026-09-05 00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			departureTime := parseInLoc(t, tt.departureTimeStr, loc)
			expectedStart := parseInLoc(t, tt.expectedStartStr, loc)
			expectedDue := parseInLoc(t, tt.expectedDueStr, loc)

			task := schedule.CalculateCheckInTask(departureTime)

			if !task.Start.Equal(expectedStart) {
				t.Errorf("task.Start = %v; want %v", task.Start, expectedStart)
			}
			if !task.Due.Equal(expectedDue) {
				t.Errorf("task.Due = %v; want %v", task.Due, expectedDue)
			}
		})
	}
}

func parseInLoc(t *testing.T, dateStr string, loc *time.Location) time.Time {
	t.Helper()
	parsed, err := time.ParseInLocation("2006-01-02 15:04", dateStr, loc)
	if err != nil {
		t.Fatalf("failed to parse date string %q: %v", dateStr, err)
	}
	return parsed
}
