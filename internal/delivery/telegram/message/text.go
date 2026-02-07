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

	switch state {
	case "awaiting_about":
		return h.handleAbout(c, sender)
	case "awaiting_support":
		return h.handleSupport(c, sender)
	}

	return nil
}

func (h *Handler) handleAbout(c tele.Context, sender *tele.User) error {
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

func (h *Handler) handleSupport(c tele.Context, sender *tele.User) error {
	if err := h.Registration.SetState(context.Background(), sender.ID, "completed"); err != nil {
		slog.Error("set state", sl.Err(err))
		return nil
	}

	admins, err := h.Users.GetAdmins(context.Background())
	if err != nil {
		slog.Error("get admins", sl.Err(err))
		return nil
	}

	content := messages.Format(messages.M.Command.Support.Ticket, map[string]string{
		"username":    sender.Username,
		"description": c.Text(),
	})

	for _, admin := range admins {
		if _, err := h.Bot.Send(&tele.User{ID: admin.TelegramID}, content); err != nil {
			slog.Error("send support to admin", sl.Err(err), "admin_id", admin.TelegramID)
		}
	}

	return c.Send(messages.M.Command.Support.ProblemSent)
}
