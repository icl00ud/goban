package bot

import (
	"fmt"
	"strings"

	"github.com/icl00ud/goban/pkg/client"
)

func fmtBoards(boards []client.BoardResponse) string {
	if len(boards) == 0 {
		return "No boards found."
	}
	var sb strings.Builder
	sb.WriteString("<b>Your Boards</b>\n\n")
	for _, b := range boards {
		sb.WriteString(fmt.Sprintf("• <b>#%d %s</b>\n", b.ID, escapeHTML(b.Name)))
		if b.Description != "" {
			sb.WriteString(fmt.Sprintf("  <i>%s</i>\n", escapeHTML(b.Description)))
		}
	}
	return sb.String()
}

func fmtBoard(b *client.BoardResponse) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>%s</b> (#%d)\n\n", escapeHTML(b.Name), b.ID))
	for _, col := range b.Columns {
		sb.WriteString(fmt.Sprintf("📋 <b>%s</b> (#%d)\n", escapeHTML(col.Title), col.ID))
		for _, card := range col.Cards {
			priority := priorityEmoji(card.Priority)
			sb.WriteString(fmt.Sprintf("  %s <code>#%d</code> %s\n", priority, card.ID, escapeHTML(card.Title)))
		}
	}
	return sb.String()
}

func fmtCard(c *client.CardResponse) string {
	priority := priorityEmoji(c.Priority)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s <b>%s</b> (#%d)\n", priority, escapeHTML(c.Title), c.ID))
	sb.WriteString(fmt.Sprintf("Column: #%d | Priority: %s\n", c.ColumnID, c.Priority))
	if c.Description != "" {
		sb.WriteString(fmt.Sprintf("\n<i>%s</i>", escapeHTML(c.Description)))
	}
	return sb.String()
}

func priorityEmoji(p string) string {
	switch p {
	case "high":
		return "🔴"
	case "medium":
		return "🟡"
	default:
		return "🟢"
	}
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
