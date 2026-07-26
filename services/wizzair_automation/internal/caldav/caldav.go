package caldav

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"avm99963-monorepo/services/wizzair_automation/internal/parser"
	"avm99963-monorepo/services/wizzair_automation/internal/schedule"
)

// Client handles pushing VTODO tasks to Nextcloud CalDAV endpoints.
type Client struct {
	BaseURL  string
	Username string
	Password string
	DryRun   bool
}

// NewClient constructs a new Nextcloud CalDAV client.
func NewClient(baseURL, username, password string, dryRun bool) *Client {
	return &Client{
		BaseURL:  strings.TrimRight(baseURL, "/"),
		Username: username,
		Password: password,
		DryRun:   dryRun,
	}
}

// PushCheckInTask constructs the iCalendar payload and issues an HTTP PUT request to CalDAV.
func (c *Client) PushCheckInTask(seg parser.FlightSegment, task schedule.CheckInTask) error {
	icsContent := GenerateVTODO(seg, task)
	filename := fmt.Sprintf("checkin-%s-%s.ics", seg.PNR, strings.ReplaceAll(seg.FlightNumber, " ", "_"))
	targetURL := fmt.Sprintf("%s/%s", c.BaseURL, filename)

	if c.DryRun {
		log.Printf("[DRY-RUN] Nextcloud CalDAV PUT to %s:\n\n---\n%s\n---\n\n", targetURL, icsContent)
		return nil
	}

	req, err := http.NewRequest("PUT", targetURL, strings.NewReader(icsContent))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "text/calendar; charset=utf-8")
	req.SetBasicAuth(c.Username, c.Password)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to issue HTTP PUT to CalDAV: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("CalDAV HTTP PUT returned error status %d", resp.StatusCode)
	}

	log.Printf("Successfully pushed VTODO check-in task for flight %s to Nextcloud CalDAV", seg.FlightNumber)
	return nil
}

// GenerateVTODO constructs an RFC 5545 iCalendar string containing a VTODO with start & due VALARM triggers.
func GenerateVTODO(seg parser.FlightSegment, task schedule.CheckInTask) string {
	summary := fmt.Sprintf("Check-in for flight %s (%s → %s)", seg.FlightNumber, seg.Origin, seg.Destination)
	uid := fmt.Sprintf("checkin-%s-%s", seg.PNR, strings.ReplaceAll(seg.FlightNumber, " ", "_"))
	dtStart := task.Start.UTC().Format("20060102T150405Z")
	dtDue := task.Due.UTC().Format("20060102T150405Z")

	return fmt.Sprintf(`BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//WizzAir Automation//EN
BEGIN:VTODO
UID:%s
SUMMARY:%s
DTSTART:%s
DUE:%s
STATUS:NEEDS-ACTION
BEGIN:VALARM
TRIGGER;RELATED=START:PT0M
ACTION:DISPLAY
DESCRIPTION:Check-in window is open for flight %s (%s → %s)
END:VALARM
BEGIN:VALARM
TRIGGER;RELATED=END:PT0M
ACTION:DISPLAY
DESCRIPTION:Check-in task is DUE NOW for flight %s (%s → %s)
END:VALARM
END:VTODO
END:VCALENDAR`, uid, summary, dtStart, dtDue, seg.FlightNumber, seg.Origin, seg.Destination, seg.FlightNumber, seg.Origin, seg.Destination)
}
