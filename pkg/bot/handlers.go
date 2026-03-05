package bot

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/icl00ud/goban/pkg/client"
	"github.com/icl00ud/goban/pkg/config"
)

type botContext struct {
	api      *tgbotapi.BotAPI
	sessions *SessionStore
	cfg      *config.Config
}

func (b *botContext) handle(update tgbotapi.Update) {
	msg := update.Message
	if msg == nil || !msg.IsCommand() {
		return
	}

	chatID := msg.Chat.ID

	// Check user allowlist (if configured)
	if len(b.cfg.TelegramAllowedUsers) > 0 && !b.isAllowed(msg.From.ID) {
		b.reply(chatID, "⛔ You are not authorized to use this bot.")
		return
	}

	args := strings.Fields(msg.CommandArguments())
	switch msg.Command() {
	case "start":
		b.reply(chatID, "👋 Welcome to <b>Goban</b>!\n\nAuthenticate with: <code>/auth &lt;token&gt;</code>\nGet token via: <code>goban login --print-token</code>\n\nSend /help for all commands.")
	case "help":
		b.reply(chatID, helpText())
	case "auth":
		if len(args) == 0 {
			b.reply(chatID, "Usage: /auth &lt;token&gt;")
			return
		}
		token := args[0]
		c := client.New(b.cfg.ServerURL, token)
		user, err := c.Me()
		if err != nil {
			b.reply(chatID, fmt.Sprintf("❌ Invalid token: %v", err))
			return
		}
		b.sessions.Set(chatID, token)
		b.reply(chatID, fmt.Sprintf("✅ Authenticated as <b>%s</b> (%s)", escapeHTML(user.Name), escapeHTML(user.Email)))
	case "logout":
		b.sessions.Delete(chatID)
		b.reply(chatID, "Logged out.")
	case "boards":
		c := b.clientFor(chatID)
		if c == nil {
			return
		}
		boards, err := c.ListBoards()
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		b.reply(chatID, fmtBoards(boards))
	case "board":
		if len(args) == 0 {
			b.reply(chatID, "Usage: /board &lt;id&gt;")
			return
		}
		id, _ := strconv.ParseUint(args[0], 10, 64)
		c := b.clientFor(chatID)
		if c == nil {
			return
		}
		board, err := c.GetBoard(uint(id))
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		b.reply(chatID, fmtBoard(board))
	case "cards":
		if len(args) == 0 {
			b.reply(chatID, "Usage: /cards &lt;board_id&gt;")
			return
		}
		boardID, _ := strconv.ParseUint(args[0], 10, 64)
		c := b.clientFor(chatID)
		if c == nil {
			return
		}
		board, err := c.GetBoard(uint(boardID))
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		b.reply(chatID, fmtBoard(board))
	case "card":
		if len(args) == 0 {
			b.reply(chatID, "Usage: /card &lt;id&gt;")
			return
		}
		id, _ := strconv.ParseUint(args[0], 10, 64)
		c := b.clientFor(chatID)
		if c == nil {
			return
		}
		card, err := c.GetCard(uint(id))
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		b.reply(chatID, fmtCard(card))
	case "add":
		// /add <column_id> <title...>
		if len(args) < 2 {
			b.reply(chatID, "Usage: /add &lt;column_id&gt; &lt;title...&gt;")
			return
		}
		colID, _ := strconv.ParseUint(args[0], 10, 64)
		title := strings.Join(args[1:], " ")
		c := b.clientFor(chatID)
		if c == nil {
			return
		}
		card, err := c.CreateCard(uint(colID), title, "", "medium")
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		b.reply(chatID, fmt.Sprintf("✅ Created card <code>#%d</code>: %s", card.ID, escapeHTML(card.Title)))
	case "done":
		// /done <card_id>  — moves card to last column of its board
		if len(args) == 0 {
			b.reply(chatID, "Usage: /done &lt;card_id&gt;")
			return
		}
		cardID, _ := strconv.ParseUint(args[0], 10, 64)
		c := b.clientFor(chatID)
		if c == nil {
			return
		}
		// Get card to find its board
		card, err := c.GetCard(uint(cardID))
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		// Find board via column — we need to list boards and find the right one
		boards, err := c.ListBoards()
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		var lastColID uint
		for _, board := range boards {
			full, err := c.GetBoard(board.ID)
			if err != nil {
				continue
			}
			for _, col := range full.Columns {
				if col.ID == card.ColumnID {
					// found the board; last column is Done
					if len(full.Columns) > 0 {
						lastColID = full.Columns[len(full.Columns)-1].ID
					}
					break
				}
			}
			if lastColID != 0 {
				break
			}
		}
		if lastColID == 0 {
			b.reply(chatID, "❌ Could not find target column")
			return
		}
		if _, err := c.MoveCard(uint(cardID), lastColID, 0); err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		b.reply(chatID, fmt.Sprintf("✅ Card #%d moved to Done column", cardID))
	case "move":
		// /move <card_id> <target_column_id>
		if len(args) < 2 {
			b.reply(chatID, "Usage: /move &lt;card_id&gt; &lt;target_column_id&gt;")
			return
		}
		cardID, _ := strconv.ParseUint(args[0], 10, 64)
		colID, _ := strconv.ParseUint(args[1], 10, 64)
		c := b.clientFor(chatID)
		if c == nil {
			return
		}
		card, err := c.MoveCard(uint(cardID), uint(colID), 0)
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		b.reply(chatID, fmt.Sprintf("✅ Card #%d moved to column #%d", card.ID, card.ColumnID))
	case "newboard":
		if len(args) == 0 {
			b.reply(chatID, "Usage: /newboard &lt;name&gt;")
			return
		}
		name := strings.Join(args, " ")
		c := b.clientFor(chatID)
		if c == nil {
			return
		}
		board, err := c.CreateBoard(name, "", "#3b82f6")
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		b.reply(chatID, fmt.Sprintf("✅ Created board <b>#%d %s</b>", board.ID, escapeHTML(board.Name)))
	case "newcol":
		// /newcol <board_id> <title...>
		if len(args) < 2 {
			b.reply(chatID, "Usage: /newcol &lt;board_id&gt; &lt;title...&gt;")
			return
		}
		boardID, _ := strconv.ParseUint(args[0], 10, 64)
		title := strings.Join(args[1:], " ")
		c := b.clientFor(chatID)
		if c == nil {
			return
		}
		col, err := c.CreateColumn(uint(boardID), title)
		if err != nil {
			b.reply(chatID, "❌ "+err.Error())
			return
		}
		b.reply(chatID, fmt.Sprintf("✅ Created column <b>#%d %s</b>", col.ID, escapeHTML(col.Title)))
	default:
		b.reply(chatID, "Unknown command. Send /help for all commands.")
	}
}

func (b *botContext) clientFor(chatID int64) *client.Client {
	token := b.sessions.Get(chatID)
	if token == "" {
		b.reply(chatID, "❌ Not authenticated. Use /auth &lt;token&gt;")
		return nil
	}
	return client.New(b.cfg.ServerURL, token)
}

func (b *botContext) reply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true
	_, _ = b.api.Send(msg)
}

func (b *botContext) isAllowed(userID int64) bool {
	for _, id := range b.cfg.TelegramAllowedUsers {
		if id == userID {
			return true
		}
	}
	return false
}

func helpText() string {
	return `<b>Goban Bot Commands</b>

<b>Auth</b>
/auth &lt;token&gt; — Authenticate
/logout — Clear session

<b>Boards</b>
/boards — List all boards
/board &lt;id&gt; — Board with columns &amp; cards
/newboard &lt;name&gt; — Create board

<b>Columns</b>
/newcol &lt;board_id&gt; &lt;title&gt; — Create column

<b>Cards</b>
/cards &lt;board_id&gt; — List all cards
/card &lt;id&gt; — Card details
/add &lt;column_id&gt; &lt;title&gt; — Create card
/done &lt;card_id&gt; — Move to last column
/move &lt;card_id&gt; &lt;column_id&gt; — Move card`
}
