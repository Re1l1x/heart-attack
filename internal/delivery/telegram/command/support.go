package command

import (
	"context"
	"log/slog"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/delivery/telegram/view"
	"github.com/jus1d/kypidbot/internal/domain"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) Support(c tele.Context) error {
	if err := h.Registration.SetState(context.Background(), c.Sender().ID, domain.UserStateAwaitingSupport); err != nil {
		slog.Error("set state", sl.Err(err))
		return nil
	}

	return c.Send(messages.M.Command.Support.Request, view.CancelSupportKeyboard())
}
