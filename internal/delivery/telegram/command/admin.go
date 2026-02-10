package command

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	"github.com/jus1d/kypidbot/internal/usecase"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) Promote(c tele.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send(messages.M.Admin.Promote.Usage)
	}

	username := strings.TrimPrefix(args[0], "@")

	err := h.Admin.Promote(context.Background(), username)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserNotFound):
			return c.Send(messages.Format(messages.M.Error.UserNotFound, map[string]string{"username": username}))
		case errors.Is(err, usecase.ErrAlreadyAdmin):
			return c.Send(messages.Format(messages.M.Error.AlreadyAdmin, map[string]string{"username": username}))
		default:
			slog.Error("promote", sl.Err(err))
			return nil
		}
	}

	return c.Send(messages.Format(messages.M.Admin.Promote.Success, map[string]string{"username": username}))
}

func (h *Handler) Demote(c tele.Context) error {
	args := c.Args()
	if len(args) == 0 {
		return c.Send(messages.M.Admin.Demote.Usage)
	}

	username := strings.TrimPrefix(args[0], "@")

	if c.Sender().Username == username {
		return c.Send(messages.M.Error.CannotDemoteYourself)
	}

	err := h.Admin.Demote(context.Background(), username)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrUserNotFound):
			return c.Send(messages.Format(messages.M.Error.UserNotFound, map[string]string{"username": username}))
		case errors.Is(err, usecase.ErrNotAdmin):
			return c.Send(messages.Format(messages.M.Error.NotAdmin, map[string]string{"username": username}))
		default:
			slog.Error("demote", sl.Err(err))
			return nil
		}
	}

	return c.Send(messages.Format(messages.M.Admin.Demote.Success, map[string]string{"username": username}))
}

func (h *Handler) AdminPanel(c tele.Context) error {
	s, err := h.Admin.GetStatistics(context.Background())
	if err != nil {
		slog.Error("get statistics", sl.Err(err))
		return c.Send("Ошибка при получении статистики")
	}

	active := s.RegisteredUsers - s.OptedOutUsers

	content := messages.Format(messages.M.Command.Admin, map[string]string{
		"total_users":       fmt.Sprintf("%d", s.TotalUsers),
		"registered_users":  fmt.Sprintf("%d", s.RegisteredUsers),
		"opted_out_users":   fmt.Sprintf("%d", s.OptedOutUsers),
		"active_users":      fmt.Sprintf("%d", active),
		"registered_daily":  fmt.Sprintf("+%d", s.RegisteredDaily),
		"registered_weekly": fmt.Sprintf("+%d", s.RegisteredWeekly),
		"male_count":        fmt.Sprintf("%d", s.MaleCount),
		"female_count":      fmt.Sprintf("%d", s.FemaleCount),
		"meetings_total":    fmt.Sprintf("%d", s.Meetings.Total),
		"meetings_confirmed": fmt.Sprintf("%d", s.Meetings.Confirmed),
		"meetings_cancelled": fmt.Sprintf("%d", s.Meetings.Cancelled),
		"meetings_pending":  fmt.Sprintf("%d", s.Meetings.Pending),
	})

	return c.Send(content)
}
