package callback

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/domain"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) Sex(c tele.Context) error {
	sender := c.Sender()
	cb := c.Callback()

	sex := "female"
	sexLabel := messages.M.UI.Buttons.Sex.Female
	if cb.Unique == "sex_male" {
		sex = "male"
		sexLabel = messages.M.UI.Buttons.Sex.Male
	}

	if err := h.Registration.SetSex(context.Background(), sender.ID, sex); err != nil {
		slog.Error("set sex", sl.Err(err))
		return c.Respond()
	}

	if err := h.Registration.SetState(context.Background(), sender.ID, domain.UserStateAwaitingAbout); err != nil {
		slog.Error("set state", sl.Err(err))
		return c.Respond()
	}

	content := fmt.Sprintf("%s\n\n%s %s", c.Message().Text, messages.M.UI.Chosen, sexLabel)
	if _, err := h.Bot.Edit(c.Message(), content); err != nil {
		slog.Error("edit sex message", sl.Err(err))
	}

	return c.Send(messages.M.Profile.About.Request)
}
