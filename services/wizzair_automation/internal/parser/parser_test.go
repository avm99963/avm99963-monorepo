package parser_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"avm99963-monorepo/services/wizzair_automation/internal/parser"
	"avm99963-monorepo/services/wizzair_automation/internal/timezone"
)

func TestParse_ValidItineraries(t *testing.T) {
	fallbackTZ := "Asia/Tokyo"
	wawLoc := timezone.GetLocation("WAW", fallbackTZ)
	bcnLoc := timezone.GetLocation("BCN", fallbackTZ)
	agpLoc := timezone.GetLocation("AGP", fallbackTZ)

	tests := []struct {
		name         string
		fixtureFile  string
		expectedSegs []parser.FlightSegment
	}{
		{
			name:        "itinerary_email_two_flights_catalan_round_trip_segments",
			fixtureFile: "itinerary_email_two_flights_catalan.html",
			expectedSegs: []parser.FlightSegment{
				{
					PNR:           "RESCOD",
					FlightNumber:  "W6 1475",
					Origin:        "WAW",
					Destination:   "BCN",
					DepartureTime: time.Date(2027, 10, 15, 5, 40, 0, 0, wawLoc),
					ArrivalTime:   time.Date(2027, 10, 15, 8, 45, 0, 0, bcnLoc),
				},
				{
					PNR:           "RESCOD",
					FlightNumber:  "W6 1478",
					Origin:        "BCN",
					Destination:   "WAW",
					DepartureTime: time.Date(2027, 10, 20, 19, 45, 0, 0, bcnLoc),
					ArrivalTime:   time.Date(2027, 10, 20, 22, 50, 0, 0, wawLoc),
				},
			},
		},
		{
			name:        "itinerary_email_one_flight_english_single_segment_AGP_to_WAW",
			fixtureFile: "itinerary_email_one_flight_english.html",
			expectedSegs: []parser.FlightSegment{
				{
					PNR:           "RESCOD",
					FlightNumber:  "W6 1360",
					Origin:        "AGP",
					Destination:   "WAW",
					DepartureTime: time.Date(2027, 11, 12, 18, 30, 0, 0, agpLoc),
					ArrivalTime:   time.Date(2027, 11, 12, 22, 30, 0, 0, wawLoc),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlBody := readSampleFile(t, tt.fixtureFile)
			segments, err := parser.Parse(htmlBody, fallbackTZ)
			if err != nil {
				t.Fatalf("expected no error parsing %s, got: %v", tt.fixtureFile, err)
			}

			if len(segments) != len(tt.expectedSegs) {
				t.Fatalf("len(segments) = %d; want %d for %s", len(segments), len(tt.expectedSegs), tt.fixtureFile)
			}

			for i, want := range tt.expectedSegs {
				got := segments[i]
				if got.PNR != want.PNR {
					t.Errorf("seg[%d].PNR = %q; want %q", i, got.PNR, want.PNR)
				}
				if got.FlightNumber != want.FlightNumber {
					t.Errorf("seg[%d].FlightNumber = %q; want %q", i, got.FlightNumber, want.FlightNumber)
				}
				if got.Origin != want.Origin {
					t.Errorf("seg[%d].Origin = %q; want %q", i, got.Origin, want.Origin)
				}
				if got.Destination != want.Destination {
					t.Errorf("seg[%d].Destination = %q; want %q", i, got.Destination, want.Destination)
				}
				if !got.DepartureTime.Equal(want.DepartureTime) {
					t.Errorf("seg[%d].DepartureTime = %v; want %v", i, got.DepartureTime, want.DepartureTime)
				}
				if got.DepartureTime.Location().String() != want.DepartureTime.Location().String() {
					t.Errorf("seg[%d].DepartureTime.Location() = %s; want %s", i, got.DepartureTime.Location(), want.DepartureTime.Location())
				}
				if !got.ArrivalTime.Equal(want.ArrivalTime) {
					t.Errorf("seg[%d].ArrivalTime = %v; want %v", i, got.ArrivalTime, want.ArrivalTime)
				}
				if got.ArrivalTime.Location().String() != want.ArrivalTime.Location().String() {
					t.Errorf("seg[%d].ArrivalTime.Location() = %s; want %s", i, got.ArrivalTime.Location(), want.ArrivalTime.Location())
				}
			}
		})
	}
}

func TestParse_NonItineraryEmail(t *testing.T) {
	fallbackTZ := "Asia/Tokyo"

	t.Run("returns_ErrNotAnItinerary_for_plain_text_promotional_email", func(t *testing.T) {
		promoBody := "<html><body><h1>Wizz Air Sale 20% off flights!</h1></body></html>"
		_, err := parser.Parse(promoBody, fallbackTZ)
		if !errors.Is(err, parser.ErrNotAnItinerary) {
			t.Errorf("expected ErrNotAnItinerary, got: %v", err)
		}
	})

	t.Run("returns_ErrNotAnItinerary_for_transactional_email_lacking_itinerary_utm_campaign", func(t *testing.T) {
		seatEmailBody := `<html><body>
			Confirmation code: RESCOD
			Flight Number: W6 1475
			<a href="https://wizzair.com/en-gb?utm_campaign=system_email_seat_selection">Select Seat</a>
		</body></html>`

		_, err := parser.Parse(seatEmailBody, fallbackTZ)
		if !errors.Is(err, parser.ErrNotAnItinerary) {
			t.Errorf("expected ErrNotAnItinerary for seat selection email, got: %v", err)
		}
	})

	t.Run("returns_ErrNotAnItinerary_for_empty_email_body", func(t *testing.T) {
		_, err := parser.Parse("", fallbackTZ)
		if !errors.Is(err, parser.ErrNotAnItinerary) {
			t.Errorf("expected ErrNotAnItinerary, got: %v", err)
		}
	})
}

func readSampleFile(t *testing.T, filename string) string {
	t.Helper()
	path := filepath.Join("../../fixtures/sample_emails", filename)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read sample itinerary file %q (path %q): %v", filename, path, err)
	}
	return string(content)
}
