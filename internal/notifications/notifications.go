package notifications

import (
	"context"
	"log/slog"
	"time"

	"github.com/jus1d/kypidbot/internal/config"
	"github.com/jus1d/kypidbot/internal/domain"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

type NotifyFunc func(ctx context.Context) error

type Notificator struct {
	bot      *tele.Bot
	users    domain.UserRepository
	places   domain.PlaceRepository
	meetings domain.MeetingRepository
	config   *config.Notifications
	settings domain.SettingsRepository
	funcs    []NotifyFunc
}

func New(c *config.Notifications, bot *tele.Bot, users domain.UserRepository, places domain.PlaceRepository, meetings domain.MeetingRepository, settings domain.SettingsRepository) *Notificator {
	return &Notificator{
		bot:      bot,
		users:    users,
		places:   places,
		meetings: meetings,
		config:   c,
		settings: settings,
	}
}

func (n *Notificator) Register(fns ...NotifyFunc) {
	n.funcs = append(n.funcs, fns...)
}

func (n *Notificator) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		for _, fn := range n.funcs {
			if err := fn(ctx); err != nil {
				slog.Error("notifications: notify func failed", sl.Err(err))
			}
		}

		sleep(ctx, n.config.PollInterval)
	}
}

func sleep(ctx context.Context, d time.Duration) {
	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}
