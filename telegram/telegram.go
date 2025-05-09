// telegram/telegram.go
package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"smtp2chat/types"
)

func SendMessage(token, chatID string, email *types.Email) {
	// Clean up body: remove "Subject:" line
	body := removeSubjectLine(email.Body)
	snippet := escapeHTML(body)
	if len(snippet) > 500 {
		snippet = snippet[:500] + "..."
	}

	msg := fmt.Sprintf(
		"<b>📩</b>\n"+
			"<b>Subject:</b> <i>%s</i>\n"+
			"<b>From:</b> <i>%s</i>\n"+
			"<b>To:</b> <i>%s</i>\n"+
			"<b>Message:</b>\n"+
			"<blockquote expandable>%s</blockquote>",
		escapeHTML(email.Subject),
		escapeHTML(email.From),
		escapeHTML(email.To),
		snippet,
	)

	if err := Send(token, chatID, msg); err != nil {
		fmt.Printf("telegram send error: %v\n", err)
	}
}

// removeSubjectLine removes any line starting with "Subject:" (case-insensitive)
func removeSubjectLine(body string) string {
	lines := strings.Split(body, "\n")
	filtered := make([]string, 0, len(lines))

	for _, line := range lines {
		if !strings.HasPrefix(strings.ToLower(line), "subject:") {
			filtered = append(filtered, line)
		}
	}

	return strings.Join(filtered, "\n")
}

// escapeHTML escapes special characters for Telegram HTML mode
func escapeHTML(text string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "<",
		">", ">",
	).Replace(text)
}

func Send(token, chatID, text string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := map[string]interface{}{
		"chat_id":    chatID,
		"text":       text,
		"parse_mode": "HTML",
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body := new(bytes.Buffer)
		_, _ = body.ReadFrom(resp.Body)
		return fmt.Errorf("telegram api error [%d]: %s", resp.StatusCode, body.String())
	}

	return nil
}
