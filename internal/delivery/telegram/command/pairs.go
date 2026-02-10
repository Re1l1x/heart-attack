package command

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) Pairs(c tele.Context) error {
	pairs, err := h.Matching.DryMatch(context.Background())
	if err != nil {
		slog.Error("dry match", sl.Err(err))
		return c.Send(messages.M.Command.Pairs.Error)
	}

	if len(pairs) == 0 {
		return c.Send(messages.M.Command.Pairs.NotFound)
	}

	var sb strings.Builder
	for _, p := range pairs {
		sb.WriteString(fmt.Sprintf("%s -- %s\n",
			messages.Mention(p.DillTelegramID, p.DillFirstName),
			messages.Mention(p.DoeTelegramID, p.DoeFirstName),
		))
	}

	return c.Send(sb.String())
}
