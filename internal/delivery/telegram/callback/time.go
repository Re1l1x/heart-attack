package callback

import (
	"context"
	"log/slog"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/delivery/telegram/view"
	"github.com/jus1d/kypidbot/internal/domain"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) ConfirmTime(c tele.Context) error {
	sender := c.Sender()

	if err := h.Registration.SetState(context.Background(), sender.ID, domain.UserStateCompleted); err != nil {
		slog.Error("set state", sl.Err(err))
		return c.Respond()
	}

	binaryStr, err := h.Registration.GetTimeRanges(context.Background(), sender.ID)
	if err != nil {
		slog.Error("get time ranges", sl.Err(err))
		return c.Respond()
	}

	selected := domain.BinaryToSet(binaryStr)
	merged := domain.MergeSelectedRanges(selected)
	summary := messages.M.UI.Chosen
	for _, tr := range merged {
		summary += "\n- " + tr
	}

	if _, err := h.Bot.Edit(c.Message(), c.Message().Text+"\n\n"+summary); err != nil {
		slog.Error("edit time message", sl.Err(err))
	}

	return c.Send(messages.M.Registration.Completed, view.ResubmitKeyboard())
}

func (h *Handler) Time(c tele.Context) error {
	sender := c.Sender()
	timeRange := c.Callback().Data

	binaryStr, err := h.Registration.GetTimeRanges(context.Background(), sender.ID)
	if err != nil {
		slog.Error("get time ranges", sl.Err(err))
		return c.Respond()
	}

	selected := domain.BinaryToSet(binaryStr)

	if selected[timeRange] {
		delete(selected, timeRange)
	} else {
		selected[timeRange] = true
	}

	newBinary := domain.SetToBinary(selected)
	if err := h.Registration.SaveTimeRanges(context.Background(), sender.ID, newBinary); err != nil {
		slog.Error("save time ranges", sl.Err(err))
		return c.Respond()
	}

	return c.Edit(messages.M.Profile.Schedule.Request, view.TimeKeyboard(selected))
}
