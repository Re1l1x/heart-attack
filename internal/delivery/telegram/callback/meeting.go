package callback

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/delivery/telegram/view"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) ConfirmMeeting(c tele.Context) error {
	data := c.Callback().Data
	meetingID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		slog.Error("parse meeting id", sl.Err(err), "data", data)
		return c.Respond()
	}

	telegramID := c.Sender().ID

	ok, err := h.Meeting.ConfirmMeeting(context.Background(), meetingID, telegramID)
	if err != nil {
		slog.Error("confirm meeting", sl.Err(err))
		return c.Respond()
	}
	if !ok {
		return c.Respond()
	}

	_ = c.Delete()

	kb := view.CancelKeyboard(fmt.Sprintf("%d", meetingID))
	originalText := c.Message().Text
	newText := originalText + "\n\n" + messages.M.Meeting.Status.Confirmed
	if _, err := h.Bot.Send(c.Chat(), newText, kb); err != nil {
		slog.Error("send confirmed message", sl.Err(err))
	}

	partnerID, err := h.Meeting.GetPartnerTelegramID(context.Background(), meetingID, telegramID)
	if err != nil {
		slog.Error("get partner telegram id", sl.Err(err))
		return nil
	}

	if partnerID != 0 {
		_, err := h.Bot.Send(&tele.User{ID: partnerID}, messages.M.Meeting.Status.PartnerConfirmed)
		if err != nil {
			slog.Error("send partner confirmed", sl.Err(err), "partner_id", partnerID)
		}
	}

	// both, meeting, err := h.Meeting.BothConfirmed(context.Background(), meetingID)
	// if err != nil {
	// 	slog.Error("check both confirmed", sl.Err(err))
	// 	return nil
	// }

	// if both && meeting != nil && meeting.PlaceID != nil && meeting.Time != nil {
	// 	placeDesc, _ := h.Meeting.GetPlaceDescription(context.Background(), *meeting.PlaceID)

	// 	finalMessage := messages.Format(messages.M.Meeting.Status.BothConfirmed, map[string]string{
	// 		"place": placeDesc,
	// 		"time":  *meeting.Time,
	// 	})

	// 	cancelKb := view.CancelKeyboard(fmt.Sprintf("%d", meetingID))

	// 	_, err := h.Bot.Send(&tele.User{ID: telegramID}, finalMessage, cancelKb)
	// 	if err != nil {
	// 		slog.Error("send both confirmed to user", sl.Err(err))
	// 	}

	// 	if partnerID != 0 {
	// 		_, err := h.Bot.Send(&tele.User{ID: partnerID}, finalMessage, cancelKb)
	// 		if err != nil {
	// 			slog.Error("send both confirmed to partner", sl.Err(err))
	// 		}
	// 	}
	// }

	return nil
}

func (h *Handler) CancelMeeting(c tele.Context) error {
	data := c.Callback().Data
	meetingID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		slog.Error("parse meeting id", sl.Err(err), "data", data)
		return c.Respond()
	}

	telegramID := c.Sender().ID

	ok, err := h.Meeting.CancelMeeting(context.Background(), meetingID, telegramID)
	if err != nil {
		slog.Error("cancel meeting", sl.Err(err))
		return c.Respond()
	}
	if !ok {
		return c.Respond()
	}

	partnerUsername, _ := h.Meeting.GetPartnerUsername(context.Background(), meetingID, telegramID)
	if partnerUsername == "" {
		partnerUsername = "unknown"
	}

	if err := h.DeleteAndSend(c, messages.Format(messages.M.Meeting.Status.Cancelled, map[string]string{
		"partner_username": partnerUsername,
	})); err != nil {
		slog.Error("send cancelled message", sl.Err(err))
	}

	userUsername, _ := h.Users.GetUserUsername(context.Background(), telegramID)
	if userUsername == "" {
		userUsername = "unknown"
	}

	partnerID, _ := h.Meeting.GetPartnerTelegramID(context.Background(), meetingID, telegramID)
	if partnerID != 0 {
		_, err := h.Bot.Send(&tele.User{ID: partnerID}, messages.Format(messages.M.Meeting.Status.PartnerCancelled, map[string]string{
			"partner_username": userUsername,
		}))
		if err != nil {
			slog.Error("send partner cancelled", sl.Err(err), "partner_id", partnerID)
		}
	}

	return nil
}
