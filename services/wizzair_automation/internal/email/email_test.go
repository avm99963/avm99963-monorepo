package email_test

import (
	"os"
	"path/filepath"
	"testing"

	"avm99963-monorepo/services/wizzair_automation/internal/email"
)

func TestProcessMessageBody(t *testing.T) {
	config := email.Config{
		LabelProcessed: "LABEL_PROCESSED",
		LabelIgnored:   "LABEL_IGNORED",
	}
	fallbackTZ := "Asia/Tokyo"

	t.Run("returns_LabelProcessed_and_segments_for_valid_itinerary", func(t *testing.T) {
		htmlBody := readSampleFile(t, "itinerary_email_one_flight_english.html")
		label, segments, err := email.ProcessMessageBody(htmlBody, fallbackTZ, config)
		if err != nil {
			t.Fatalf("expected no error processing valid email body, got: %v", err)
		}
		if label != "LABEL_PROCESSED" {
			t.Errorf("label = %q; want %q", label, "LABEL_PROCESSED")
		}
		if len(segments) != 1 {
			t.Errorf("len(segments) = %d; want 1", len(segments))
		}
	})

	t.Run("returns_LabelProcessed_for_multipart_mime_email", func(t *testing.T) {
		htmlPart := readSampleFile(t, "itinerary_email_one_flight_english.html")
		mimeRaw := "Content-Type: multipart/alternative; boundary=\"boundary123\"\r\n" +
			"\r\n" +
			"--boundary123\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			"Plain text snippet\r\n" +
			"--boundary123\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			htmlPart + "\r\n" +
			"--boundary123--\r\n"

		label, segments, err := email.ProcessMessageBody(mimeRaw, fallbackTZ, config)
		if err != nil {
			t.Fatalf("expected no error processing multipart email, got: %v", err)
		}
		if label != "LABEL_PROCESSED" {
			t.Errorf("label = %q; want %q", label, "LABEL_PROCESSED")
		}
		if len(segments) != 1 {
			t.Errorf("len(segments) = %d; want 1", len(segments))
		}
	})

	t.Run("returns_LabelIgnored_and_nil_segments_for_non_itinerary_email", func(t *testing.T) {
		htmlBody := "<html><body><h1>Promotional Sale!</h1></body></html>"
		label, segments, err := email.ProcessMessageBody(htmlBody, fallbackTZ, config)
		if err != nil {
			t.Fatalf("expected no error for non-itinerary email, got: %v", err)
		}
		if label != "LABEL_IGNORED" {
			t.Errorf("label = %q; want %q", label, "LABEL_IGNORED")
		}
		if len(segments) != 0 {
			t.Errorf("len(segments) = %d; want 0", len(segments))
		}
	})
}

func TestEnsureLabelsExist_DryRun(t *testing.T) {
	config := email.Config{
		DryRun: true,
	}
	client := email.NewClient(config)

	t.Run("dry_run_mode_skips_network_calls", func(t *testing.T) {
		err := client.EnsureLabelsExist("Processed", "Processed (ignored)")
		if err != nil {
			t.Errorf("expected no error in dry run mode, got: %v", err)
		}
	})
}

func TestMoveMessage_DryRun(t *testing.T) {
	config := email.Config{
		DryRun: true,
	}
	client := email.NewClient(config)

	t.Run("dry_run_mode_skips_imap_network_calls", func(t *testing.T) {
		err := client.MoveMessage(123, "LABEL_PROCESSED")
		if err != nil {
			t.Errorf("expected no error in dry run mode, got: %v", err)
		}
	})
}

func readSampleFile(t *testing.T, filename string) string {
	t.Helper()
	path := filepath.Join("../../fixtures/sample_emails", filename)
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read sample email file %q (path %q): %v", filename, path, err)
	}
	return string(content)
}
