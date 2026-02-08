package callback

import (
	"github.com/jus1d/kypidbot/internal/config/messages"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) HowItWorks(c tele.Context) error {
	_ = c.Respond()
	text := messages.M.Start.Welcome + "\n\n" + messages.M.Command.About
	return c.Edit(text)
}
