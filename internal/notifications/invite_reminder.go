package notifications

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/domain"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

const (
	referralCodeLen     = 8
	referralCodeCharset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func (n *Notificator) InviteReminder(ctx context.Context) error {
	value, _ := n.settings.Get(ctx, "registration_closed")
	closed := value == "true"
	if closed {
		return nil
	}

	list, err := n.users.GetForInviteReminder(ctx, n.config.InviteReminderIn)
	if err != nil {
		return err
	}

	for _, u := range list {
		log := slog.With(slog.String("username", u.Username), slog.Int64("telegram_id", u.TelegramID))

		code := u.ReferralCode
		if code == "" {
			code, err = domain.GenerateReferralCode()
			if err != nil {
				log.Error("notifications: generate referral code", sl.Err(err))
				continue
			}

			if err := n.users.SetReferralCode(ctx, u.TelegramID, code); err != nil {
				log.Error("notifications: set referral code", sl.Err(err))
				continue
			}
		}

		link := fmt.Sprintf("https://t.me/%s?start=%s", n.bot.Me.Username, code)
		text := messages.Format(messages.M.Notifications.Invite, map[string]string{"link": link})

		if _, err := n.bot.Send(&tele.User{ID: u.TelegramID}, text); err != nil {
			log.Error("notifications: send invite reminder", sl.Err(err))
		}

		if err := n.users.MarkInviteNotified(ctx, u.TelegramID); err != nil {
			log.Error("notifications: mark invite notified", sl.Err(err))
		}
	}

	return nil
}
