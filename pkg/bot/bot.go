// Package bot provides the Telegram bot for Goban.
package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/icl00ud/goban/pkg/config"
)

// Start launches the Telegram bot using long-polling.
func Start() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	if cfg.TelegramBotToken == "" {
		return fmt.Errorf("TELEGRAM_BOT_TOKEN is not set; configure it via env or ~/.config/goban/config.json")
	}

	api, err := tgbotapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		return fmt.Errorf("creating bot: %w", err)
	}

	sessions, err := newSessionStore()
	if err != nil {
		return fmt.Errorf("creating session store: %w", err)
	}

	log.Printf("Goban Telegram bot started as @%s", api.Self.UserName)

	bctx := &botContext{
		api:      api,
		sessions: sessions,
		cfg:      cfg,
	}

	updateCfg := tgbotapi.NewUpdate(0)
	updateCfg.Timeout = 60

	updates := api.GetUpdatesChan(updateCfg)
	for update := range updates {
		go bctx.handle(update)
	}
	return nil
}
