package callback

import (
	"context"
	"log/slog"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/domain"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) CancelSupport(c tele.Context) error {
	err := h.Registration.SetState(context.Background(), c.Sender().ID, domain.UserStateCompleted)
	if err != nil {
		slog.Error("set state", sl.Err(err))
		return c.Respond()
	}

	err = h.DeleteAndSend(c, messages.M.Command.Support.Cancelled)
	if err != nil {
		slog.Error("send cancelled message", sl.Err(err))
	}
	return c.Respond()
}
