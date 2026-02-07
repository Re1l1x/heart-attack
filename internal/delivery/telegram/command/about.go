package command

import (
	"github.com/jus1d/kypidbot/internal/config/messages"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) About(c tele.Context) error {
	return c.Send(messages.M.Command.About)
}
