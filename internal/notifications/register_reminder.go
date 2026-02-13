package notifications

import (
	"context"
	"log/slog"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (n *Notificator) RegisterReminder(ctx context.Context) error {
	value, _ := n.settings.Get(ctx, "registration_closed")
	closed := value == "true"
	if closed {
		return nil
	}

	list, err := n.users.GetNotCompleted(ctx, n.config.RegistrationReminderIn)
	if err != nil {
		return err
	}

	for _, u := range list {
		log := slog.With(slog.String("username", u.Username), slog.Int64("telegram_id", u.TelegramID))

		if u.IsAdmin {
			continue
		}

		content := messages.M.Notifications.Registration
		if _, err := n.bot.Send(&tele.User{ID: u.TelegramID}, content); err != nil {
			log.Error("notifications: send registration notification", sl.Err(err))
		}

		if err := n.users.MarkNotified(ctx, u.TelegramID); err != nil {
			log.Error("notifications: mark user as notified", sl.Err(err))
		}
	}

	return nil
}
