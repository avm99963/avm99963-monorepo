package schedule

import (
	"time"
)

// CheckInTask holds the start date and due date for a flight check-in reminder task.
type CheckInTask struct {
	Start time.Time
	Due   time.Time
}

// CalculateCheckInTask calculates the start and due timestamps for a check-in task
// given a flight departure time (which includes its local airport timezone).
func CalculateCheckInTask(departureTime time.Time) CheckInTask {
	// Check-in opens 24h before departure
	checkInOpens := departureTime.Add(-24 * time.Hour)

	start := adjustForRestPeriodStart(checkInOpens)

	// Due time is start + 2 hours
	prelimDue := start.Add(2 * time.Hour)
	due := adjustForRestPeriodDue(prelimDue)

	return CheckInTask{
		Start: start,
		Due:   due,
	}
}

// adjustForRestPeriodStart shifts a start timestamp to 08:30 if it falls in 00:00-08:30 rest period.
func adjustForRestPeriodStart(t time.Time) time.Time {
	loc := t.Location()
	hour, min, sec := t.Clock()
	secondsOfDay := hour*3600 + min*60 + sec

	// Rest period: 00:00:00 to 08:30:00 (30600 seconds)
	if secondsOfDay < 8*3600+30*60 {
		return time.Date(t.Year(), t.Month(), t.Day(), 8, 30, 0, 0, loc)
	}
	return t
}

// adjustForRestPeriodDue shifts a due timestamp to 00:00 (start of rest period) if it falls in 00:00-08:30.
func adjustForRestPeriodDue(t time.Time) time.Time {
	loc := t.Location()
	hour, min, sec := t.Clock()
	secondsOfDay := hour*3600 + min*60 + sec

	if secondsOfDay < 8*3600+30*60 {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	}
	return t
}
