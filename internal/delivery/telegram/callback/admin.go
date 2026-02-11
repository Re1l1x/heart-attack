package callback

import (
	"context"
	"log/slog"

	"github.com/jus1d/kypidbot/internal/delivery/telegram/view"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) RefreshAdmin(c tele.Context) error {
	content, err := h.Admin.FormatPanel(context.Background())
	if err != nil {
		slog.Error("refresh admin panel", sl.Err(err))
		return nil
	}

	if _, err = h.Bot.Edit(c.Message(), content, view.RefreshAdminKeyboard()); err != nil {
		if err == tele.ErrSameMessageContent {
			return c.Respond(&tele.CallbackResponse{Text: "Нет изменений"})
		}
		c.Respond()
		slog.Error("edit admin panel", sl.Err(err))
	}

	return c.Respond()
}
