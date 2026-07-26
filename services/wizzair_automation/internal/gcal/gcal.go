package gcal

import (
	"context"
	"fmt"
	"log"
	"time"

	"avm99963-monorepo/services/wizzair_automation/internal/parser"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

// Client handles Google Calendar sync operations.
type Client struct {
	CalendarID      string
	CredentialsPath string
	DryRun          bool
}

// NewClient constructs a new Google Calendar sync client.
func NewClient(calendarID, credentialsPath string, dryRun bool) *Client {
	if calendarID == "" {
		calendarID = "primary"
	}
	return &Client{
		CalendarID:      calendarID,
		CredentialsPath: credentialsPath,
		DryRun:          dryRun,
	}
}

// SyncFlightEvents inserts/upserts flight events to Google Calendar.
func (c *Client) SyncFlightEvents(segments []parser.FlightSegment) error {
	for _, seg := range segments {
		event := BuildFlightEvent(seg)

		if c.DryRun {
			log.Printf("[DRY-RUN] Google Calendar Event: %s (Location: %s, Start: %s, End: %s, ICalUID: %s)",
				event.Summary, event.Location, event.Start.DateTime, event.End.DateTime, event.ICalUID)
			continue
		}

		ctx := context.Background()
		srv, err := calendar.NewService(ctx, option.WithCredentialsFile(c.CredentialsPath), option.WithScopes(calendar.CalendarEventsScope))
		if err != nil {
			return fmt.Errorf("failed to create calendar service: %w", err)
		}

		_, err = srv.Events.Insert(c.CalendarID, event).Do()
		if err != nil {
			log.Printf("WARNING: Failed to insert event %s (may already exist): %v", event.ICalUID, err)
		} else {
			log.Printf("Successfully created Google Calendar event: %s (%s)", event.Summary, event.ICalUID)
		}
	}
	return nil
}

// BuildFlightEvent creates a Google Calendar Event struct for a given flight segment.
func BuildFlightEvent(seg parser.FlightSegment) *calendar.Event {
	summary := fmt.Sprintf("🛫 %s → %s", seg.Origin, seg.Destination)
	location := fmt.Sprintf("%s Airport", seg.Origin)
	iCalUID := fmt.Sprintf("wizz-flight-%s-%s", seg.PNR, seg.FlightNumber)
	description := fmt.Sprintf("Flight: %s\nPNR: %s\nRoute: %s → %s", seg.FlightNumber, seg.PNR, seg.Origin, seg.Destination)

	dep := seg.DepartureTime
	arr := seg.ArrivalTime

	return &calendar.Event{
		Summary:     summary,
		Description: description,
		Location:    location,
		ICalUID:     iCalUID,
		Start: &calendar.EventDateTime{
			DateTime: dep.Format(time.RFC3339),
			TimeZone: dep.Location().String(),
		},
		End: &calendar.EventDateTime{
			DateTime: arr.Format(time.RFC3339),
			TimeZone: arr.Location().String(),
		},
	}
}
