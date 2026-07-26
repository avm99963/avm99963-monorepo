package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"avm99963-monorepo/services/wizzair_automation/internal/parser"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
)

// Config holds IMAP connection and label management settings.
type Config struct {
	Server         string
	Port           int
	Username       string
	Password       string
	LabelProcessed string
	LabelIgnored   string
	DryRun         bool
}

// Client manages IMAP operations using go-imap/v2.
type Client struct {
	Config Config
}

// Message holds an unread email's sequence number and HTML body.
type Message struct {
	SeqNum   uint32
	HTMLBody string
}

// NewClient constructs a new IMAP email client.
func NewClient(config Config) *Client {
	return &Client{Config: config}
}

// Connect opens an IMAP TLS connection and logs in to the server.
func (c *Client) Connect() (*imapclient.Client, error) {
	addr := fmt.Sprintf("%s:%d", c.Config.Server, c.Config.Port)

	var client *imapclient.Client
	var err error
	if c.Config.Port == 993 {
		client, err = imapclient.DialTLS(addr, &imapclient.Options{
			TLSConfig: &tls.Config{ServerName: c.Config.Server},
		})
	} else {
		client, err = imapclient.DialStartTLS(addr, &imapclient.Options{
			TLSConfig: &tls.Config{ServerName: c.Config.Server},
		})
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to IMAP server %s: %w", addr, err)
	}

	if err := client.Login(c.Config.Username, c.Config.Password).Wait(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed IMAP login for user %s: %w", c.Config.Username, err)
	}

	return client, nil
}

// EnsureLabelsExist verifies that the specified IMAP folders/labels exist, and creates them if missing.
func (c *Client) EnsureLabelsExist(labels ...string) error {
	if c.Config.DryRun {
		return nil
	}

	client, err := c.Connect()
	if err != nil {
		return err
	}
	defer client.Logout()

	listCmd := client.List("", "*", nil)
	existing := make(map[string]bool)
	for {
		mbox := listCmd.Next()
		if mbox == nil {
			break
		}
		existing[mbox.Mailbox] = true
	}
	if err := listCmd.Close(); err != nil {
		return fmt.Errorf("failed to list IMAP mailboxes: %w", err)
	}

	for _, label := range labels {
		if label == "" || existing[label] {
			continue
		}
		log.Printf("IMAP label %q does not exist; creating mailbox...", label)
		if err := client.Create(label, nil).Wait(); err != nil {
			return fmt.Errorf("failed to create IMAP mailbox %q: %w", label, err)
		}
		log.Printf("Successfully created IMAP mailbox %q", label)
	}

	return nil
}

// FetchUnreadMessages selects INBOX, searches UNSEEN messages, and fetches raw body content.
func (c *Client) FetchUnreadMessages() ([]Message, error) {
	if c.Config.DryRun {
		log.Printf("[DRY-RUN] Fetching unread emails from INBOX at %s:%d", c.Config.Server, c.Config.Port)
		return nil, nil
	}

	client, err := c.Connect()
	if err != nil {
		return nil, err
	}
	defer client.Logout()

	if _, err := client.Select("INBOX", nil).Wait(); err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %w", err)
	}

	searchData, err := client.Search(&imap.SearchCriteria{
		NotFlag: []imap.Flag{imap.FlagSeen},
	}, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to search UNSEEN messages: %w", err)
	}

	seqNums := searchData.AllSeqNums()
	if len(seqNums) == 0 {
		return nil, nil
	}

	var seqSet imap.SeqSet
	seqSet.AddNum(seqNums...)

	bodySection := &imap.FetchItemBodySection{}
	fetchCmd := client.Fetch(seqSet, &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{bodySection},
	})
	defer fetchCmd.Close()

	var messages []Message
	for {
		msg := fetchCmd.Next()
		if msg == nil {
			break
		}

		var bodyBuf strings.Builder
		for {
			item := msg.Next()
			if item == nil {
				break
			}
			if literal, ok := item.(imapclient.FetchItemDataBodySection); ok {
				b, _ := io.ReadAll(literal.Literal)
				bodyBuf.Write(b)
			}
		}

		messages = append(messages, Message{
			SeqNum:   msg.SeqNum,
			HTMLBody: bodyBuf.String(),
		})
	}

	if err := fetchCmd.Close(); err != nil {
		return nil, fmt.Errorf("failed fetching message bodies: %w", err)
	}

	return messages, nil
}

// ProcessMessageBody inspects an email raw body (handling both single-part and MIME multipart messages),
// determines target label, and extracts flight segments if a valid itinerary is present.
func ProcessMessageBody(rawBody string, fallbackTZ string, config Config) (string, []parser.FlightSegment, error) {
	candidates := ExtractCandidateBodies(rawBody)

	var lastErr error
	for _, candidate := range candidates {
		segments, err := parser.Parse(candidate, fallbackTZ)
		if err == nil {
			return config.LabelProcessed, segments, nil
		} else {
			if lastErr != nil && !errors.Is(lastErr, parser.ErrNotAnItinerary) {
				log.Printf("Error parsing email body (if it is multiparted, other parts will be attempted): %w", lastErr)
			}
			lastErr = err
		}
	}

	if lastErr != nil && !errors.Is(lastErr, parser.ErrNotAnItinerary) {
		return "", nil, fmt.Errorf("error parsing email body; error for last part (if multiparted): %w", lastErr)
	}

	return config.LabelIgnored, nil, nil
}

// MoveMessage moves a message by sequence number from INBOX to targetLabel.
func (c *Client) MoveMessage(seqNum uint32, targetLabel string) error {
	if c.Config.DryRun {
		log.Printf("[DRY-RUN] IMAP move message %d to label %q", seqNum, targetLabel)
		return nil
	}

	client, err := c.Connect()
	if err != nil {
		return err
	}
	defer client.Logout()

	if _, err := client.Select("INBOX", nil).Wait(); err != nil {
		return fmt.Errorf("failed to select INBOX: %w", err)
	}

	var seqSet imap.SeqSet
	seqSet.AddNum(seqNum)

	if _, err := client.Move(seqSet, targetLabel).Wait(); err != nil {
		return fmt.Errorf("failed to move message %d to label %s: %w", seqNum, targetLabel, err)
	}

	return nil
}

// ExtractCandidateBodies parses rawBody as a MIME email or returns rawBody directly.
// It recursively extracts all text parts (text/html and text/plain) from multipart messages.
func ExtractCandidateBodies(rawBody string) []string {
	mr, err := mail.CreateReader(strings.NewReader(rawBody))
	if err != nil {
		return []string{rawBody}
	}

	var candidates []string
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, _ := h.ContentType()
			if strings.HasPrefix(contentType, "text/") {
				b, err := io.ReadAll(p.Body)
				if err == nil && len(b) > 0 {
					candidates = append(candidates, string(b))
				}
			}
		case *mail.AttachmentHeader:
			contentType, _, _ := h.ContentType()
			if strings.HasPrefix(contentType, "text/") {
				b, err := io.ReadAll(p.Body)
				if err == nil && len(b) > 0 {
					candidates = append(candidates, string(b))
				}
			}
		}
	}

	if len(candidates) == 0 {
		return []string{rawBody}
	}

	return candidates
}
