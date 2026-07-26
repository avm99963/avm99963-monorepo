package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"avm99963-monorepo/services/wizzair_automation/internal/caldav"
	"avm99963-monorepo/services/wizzair_automation/internal/email"
	"avm99963-monorepo/services/wizzair_automation/internal/gcal"
	"avm99963-monorepo/services/wizzair_automation/internal/schedule"
)

func main() {
	log.Println("Starting Wizz Air Flight & Check-in Automation service...")

	pollIntervalStr := getEnv("POLL_INTERVAL", "12h")
	pollInterval, err := time.ParseDuration(pollIntervalStr)
	if err != nil {
		log.Printf("Invalid POLL_INTERVAL %q: %v", pollIntervalStr, err)
		return
	}

	fallbackTZ := getEnv("TIMEZONE", "Europe/Warsaw")
	dryRun := getEnvBool("DRY_RUN")

	imapPort, _ := strconv.Atoi(getEnv("IMAP_PORT", "993"))
	emailConfig := email.Config{
		Server:         getEnv("IMAP_SERVER", ""),
		Port:           imapPort,
		Username:       getEnv("IMAP_USERNAME", ""),
		Password:       getEnv("IMAP_PASSWORD", ""),
		LabelProcessed: getEnv("LABEL_PROCESSED", "Processed"),
		LabelIgnored:   getEnv("LABEL_IGNORED", "Processed (ignored)"),
		DryRun:         dryRun,
	}

	emailClient := email.NewClient(emailConfig)

	if err := emailClient.EnsureLabelsExist(emailConfig.LabelProcessed, emailConfig.LabelIgnored); err != nil {
		log.Printf("Warning: Error ensuring IMAP labels exist at startup: %v", err)
	}

	gcalClient := gcal.NewClient(
		getEnv("GCAL_CALENDAR_ID", "primary"),
		getEnv("GCAL_SERVICE_ACCOUNT_JSON_PATH", ""),
		dryRun,
	)

	caldavClient := caldav.NewClient(
		getEnv("CALDAV_BASE_URL", ""),
		getEnv("CALDAV_USERNAME", ""),
		getEnv("CALDAV_PASSWORD", ""),
		dryRun,
	)

	processCycle := func() {
		log.Println("Running email polling and processing cycle...")

		messages, err := emailClient.FetchUnreadMessages()
		if err != nil {
			log.Printf("Error fetching unread messages: %v", err)
			return
		}

		for _, msg := range messages {
			log.Printf("Processing message %d...", msg.SeqNum)

			targetLabel, segments, err := email.ProcessMessageBody(msg.HTMLBody, fallbackTZ, emailConfig)
			if err != nil {
				log.Printf("Error processing message %d: %v", msg.SeqNum, err)
				continue
			}

			// Sync Google Calendar Events
			if err := gcalClient.SyncFlightEvents(segments); err != nil {
				log.Printf("Error syncing Google Calendar for message %d: %v", msg.SeqNum, err)
				continue // Skip moving label so we retry next cycle
			}

			// Sync CalDAV Check-in Tasks
			syncFailed := false
			for _, seg := range segments {
				task := schedule.CalculateCheckInTask(seg.DepartureTime)
				if err := caldavClient.PushCheckInTask(seg, task); err != nil {
					log.Printf("Error pushing CalDAV check-in task for flight %s: %v", seg.FlightNumber, err)
					syncFailed = true
					break
				}
			}

			if syncFailed {
				continue // Skip moving label so we retry next cycle
			}

			// Only move label after both Google Calendar and CalDAV sync succeed
			if err := emailClient.MoveMessage(msg.SeqNum, targetLabel); err != nil {
				log.Printf("Error moving message %d to label %s: %v", msg.SeqNum, targetLabel, err)
			}

			log.Printf("Finished processing message %d", msg.SeqNum)
		}

		log.Println("Finished email polling and processing cycle.")
	}

	processCycle()

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		processCycle()
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvBool(key string) bool {
	val := strings.ToLower(os.Getenv(key))
	return val == "true"
}
