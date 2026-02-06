package message

import (
	"context"
	"log/slog"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/delivery/telegram/view"
	"github.com/jus1d/kypidbot/internal/domain"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) Text(c tele.Context) error {
	sender := c.Sender()

	state, err := h.Registration.GetState(context.Background(), sender.ID)
	if err != nil {
		slog.Error("get state", sl.Err(err))
		return nil
	}

	if state != "awaiting_about" {
		return nil
	}

	if err := h.Registration.SetAbout(context.Background(), sender.ID, c.Text()); err != nil {
		slog.Error("set about", sl.Err(err))
		return nil
	}

	if err := h.Registration.SetState(context.Background(), sender.ID, "awaiting_time"); err != nil {
		slog.Error("set state", sl.Err(err))
		return nil
	}

	binaryStr, err := h.Registration.GetTimeRanges(context.Background(), sender.ID)
	if err != nil {
		slog.Error("get time ranges", sl.Err(err))
		return nil
	}

	selected := domain.BinaryToSet(binaryStr)

	return c.Send(messages.M.Profile.Schedule.Request, view.TimeKeyboard(selected))
}
