// teams/teams.go
package teams

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"smtp2chat/types"
)

func SendMessage(webhook string, email *types.Email, verbose bool) {
	snippet := cleanBody(email.Body)
	if len(snippet) > 500 {
		snippet = snippet[:500] + "..."
	}

	card := adaptiveCardMessage(email.Subject, email.From, email.To, snippet)

	data, err := json.Marshal(card)
	if err != nil {
		if verbose {
			fmt.Printf("failed to marshal card: %v\n", err)
		}
		return
	}

	if verbose {
		fmt.Println("sending teams message...")
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(data))
	if err != nil {
		if verbose {
			fmt.Printf("teams send error: %v\n", err)
		}
		return
	}
	defer resp.Body.Close()

	if verbose {
		if resp.StatusCode >= 400 {
			fmt.Printf("teams API error: %d\n", resp.StatusCode)
		} else {
			fmt.Println("teams message sent successfully")
		}
	}
}

func cleanBody(body string) string {
	lines := strings.Split(body, "\n")
	filtered := make([]string, 0, len(lines))

	for _, line := range lines {
		if !strings.HasPrefix(strings.ToLower(line), "subject:") {
			filtered = append(filtered, line)
		}
	}

	return strings.Join(filtered, "\n")
}

// adaptiveCardMessage returns an Adaptive Card payload for Teams
func adaptiveCardMessage(subject, from, to, snippet string) map[string]interface{} {
	return map[string]interface{}{
		"type": "message",
		"attachments": []map[string]interface{}{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"contentUrl":  nil,
				"content": map[string]interface{}{
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"type":    "AdaptiveCard",
					"version": "1.2",
					"body": []map[string]interface{}{
						{
							"type":   "TextBlock",
							"text":   "**📩 New Email Received**",
							"size":   "Medium",
							"weight": "Bolder",
						},
						{
							"type": "FactSet",
							"facts": []map[string]string{
								{"title": "Subject", "value": subject},
								{"title": "From", "value": from},
								{"title": "To", "value": to},
							},
						},
						{
							"type":   "TextBlock",
							"text":   "Message Snippet:",
							"weight": "Bolder",
						},
						{
							"type": "TextBlock",
							"text": escapeHTMLForTeams(snippet),
							"wrap": true,
						},
					},
					"actions": []map[string]interface{}{
						{
							"type":  "Action.OpenUrl",
							"title": "Read More",
							"url":   "https://adaptivecards.io",
						},
					},
				},
			},
		},
	}
}

// escapeHTMLForTeams escapes special characters for Teams Adaptive Card rendering
func escapeHTMLForTeams(text string) string {
	return strings.NewReplacer(
		"&", "&amp;",
		"<", "<",
		">", ">",
	).Replace(text)
}
