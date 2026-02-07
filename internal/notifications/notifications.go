package notifications

import (
	"context"
	"log/slog"
	"time"

	"github.com/jus1d/kypidbot/internal/config"
	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/domain"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func Go(ctx context.Context, c *config.Notifications, bot *tele.Bot, users domain.UserRepository, places domain.PlaceRepository, meetings domain.MeetingRepository) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		list, err := meetings.GetMeetingsStartingIn(ctx, c.UpcomingIn)
		if err != nil {
			slog.Error("notifications: get meetings", sl.Err(err))
			sleep(ctx, c.PollInterval)
			continue
		}

		for _, m := range list {
			if m.UsersNotified {
				continue
			}

			if m.DillState != domain.StateConfirmed || m.DoeState != domain.StateConfirmed {
				continue
			}

			if m.PlaceID == nil || m.Time == nil {
				continue
			}

			dill, err := users.GetUser(ctx, m.DillID)
			if err != nil {
				slog.Error("notifications: get dill", sl.Err(err))
				continue
			}

			doe, err := users.GetUser(ctx, m.DoeID)
			if err != nil {
				slog.Error("notifications: get doe", sl.Err(err))
				continue
			}

			if dill == nil || doe == nil {
				continue
			}

			place, err := places.GetPlaceDescription(ctx, *m.PlaceID)
			if err != nil {
				slog.Error("notifications: get place description", sl.Err(err))
				continue
			}

			msg := messages.Format(messages.M.Meeting.Reminder, map[string]string{
				"place": place,
				"time":  domain.Timef(*m.Time),
			})

			if _, err := bot.Send(&tele.User{ID: dill.TelegramID}, msg); err != nil {
				slog.Error("notifications: send to dill", sl.Err(err), "telegram_id", dill.TelegramID)
			}

			if _, err := bot.Send(&tele.User{ID: doe.TelegramID}, msg); err != nil {
				slog.Error("notifications: send to doe", sl.Err(err), "telegram_id", doe.TelegramID)
			}

			if err := meetings.MarkNotified(ctx, m.ID); err != nil {
				slog.Error("notifications: mark notified", sl.Err(err), "meeting_id", m.ID)
			}
		}

		sleep(ctx, c.PollInterval)
	}
}

func sleep(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
