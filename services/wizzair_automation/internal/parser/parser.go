package parser

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"avm99963-monorepo/services/wizzair_automation/internal/timezone"
)

// ErrNotAnItinerary is returned when an email body does not contain a valid Wizz Air itinerary confirmation.
var ErrNotAnItinerary = errors.New("email does not contain a valid Wizz Air itinerary confirmation")

const itineraryUTMCampaign = "utm_campaign=system_email_itinerary_confirmation"

// FlightSegment represents a single flight segment extracted from an itinerary.
type FlightSegment struct {
	PNR           string
	FlightNumber  string
	Origin        string
	Destination   string
	DepartureTime time.Time
	ArrivalTime   time.Time
}

var (
	pnrURLRegex    = regexp.MustCompile(`(?i)wizzair\.com/[a-z]{2}-[a-z]{2}/itinerary/([A-Z0-9]{6})/`)
	flightNumRegex = regexp.MustCompile(`(?i)(?:Flight Number|Número de vol)[:\s]*(W6?[\s\xA0]*\d+)`)
	iataRegex      = regexp.MustCompile(`\(([A-Z]{3})\)`)
	dateTimeRegex  = regexp.MustCompile(`(\d{1,2})[/.-](\d{1,2})[/.-](\d{4})[\s\xA0]+(\d{1,2}):(\d{2})`)
)

// Parse inspects an email HTML body, resolves airport timezones using fallbackTZ,
// and extracts all contained Wizz Air flight segments.
func Parse(htmlBody string, fallbackTZ string) ([]FlightSegment, error) {
	if strings.TrimSpace(htmlBody) == "" {
		return nil, ErrNotAnItinerary
	}

	// If the UTM campaign created by Wizz Air for confirmation emails isn't part
	// of the email, we probably are dealing with another type of email sent by
	// Wizz Air or even an email sent by another company.
	if !strings.Contains(htmlBody, itineraryUTMCampaign) {
		return nil, ErrNotAnItinerary
	}

	pnr := extractPNR(htmlBody)
	if pnr == "" {
		return nil, ErrNotAnItinerary
	}

	segments := extractSegments(htmlBody, pnr, fallbackTZ)
	if len(segments) == 0 {
		return nil, ErrNotAnItinerary
	}

	return segments, nil
}

func extractPNR(htmlBody string) string {
	if m := pnrURLRegex.FindStringSubmatch(htmlBody); len(m) > 1 {
		return strings.ToUpper(m[1])
	}
	return ""
}

func extractSegments(htmlBody string, pnr string, fallbackTZ string) []FlightSegment {
	var segments []FlightSegment

	flightNumMatches := flightNumRegex.FindAllStringSubmatchIndex(htmlBody, -1)
	if len(flightNumMatches) == 0 {
		return nil
	}

	for i, loc := range flightNumMatches {
		var endIdx int
		if i+1 < len(flightNumMatches) {
			endIdx = flightNumMatches[i+1][0]
		} else {
			endIdx = len(htmlBody)
		}
		block := htmlBody[loc[0]:endIdx]

		rawFlightNum := htmlBody[loc[2]:loc[3]]
		flightNum := strings.ReplaceAll(rawFlightNum, "\u00a0", " ")
		flightNum = strings.Join(strings.Fields(flightNum), " ")

		iatas := iataRegex.FindAllStringSubmatch(block, -1)
		if len(iatas) < 2 {
			continue
		}
		origin := iatas[0][1]
		destination := iatas[1][1]

		dateMatches := dateTimeRegex.FindAllStringSubmatch(block, -1)
		if len(dateMatches) < 2 {
			continue
		}

		depLoc := timezone.GetLocation(origin, fallbackTZ)
		arrLoc := timezone.GetLocation(destination, fallbackTZ)

		depTime, errDep := parseDateTime(dateMatches[0], depLoc)
		arrTime, errArr := parseDateTime(dateMatches[1], arrLoc)
		if errDep != nil || errArr != nil {
			continue
		}

		segments = append(segments, FlightSegment{
			PNR:           pnr,
			FlightNumber:  flightNum,
			Origin:        origin,
			Destination:   destination,
			DepartureTime: depTime,
			ArrivalTime:   arrTime,
		})
	}

	return segments
}

func parseDateTime(m []string, loc *time.Location) (time.Time, error) {
	day, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	year, _ := strconv.Atoi(m[3])
	hour, _ := strconv.Atoi(m[4])
	minute, _ := strconv.Atoi(m[5])

	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, loc), nil
}
