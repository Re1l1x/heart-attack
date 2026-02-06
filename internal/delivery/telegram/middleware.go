package telegram

import (
	"context"
	"log/slog"

	tele "gopkg.in/telebot.v3"
)

func (b *Bot) AdminOnly(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		isAdmin, err := b.users.IsAdmin(context.Background(), c.Sender().ID)
		if err != nil || !isAdmin {
			return nil
		}
		return next(c)
	}
}

func LogUpdates(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		sender := c.Sender()
		attrs := []any{
			"telegram_id", sender.ID,
			"username", sender.Username,
		}

		if cb := c.Callback(); cb != nil {
			slog.Debug("callback", append(attrs, "unique", cb.Unique, "data", cb.Data)...)
		} else if msg := c.Message(); msg != nil {
			if msg.Text != "" && msg.Text[0] == '/' {
				slog.Debug("command", append(attrs, "text", msg.Text)...)
			} else if msg.Sticker != nil {
				slog.Debug("sticker", attrs...)
			} else {
				slog.Debug("message", append(attrs, "text", msg.Text)...)
			}
		}

		return next(c)
	}
}
