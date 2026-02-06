package command

import (
	"context"
	"fmt"

	"github.com/jus1d/kypidbot/internal/delivery/telegram/view"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

const mmSticker = "CAACAgIAAxkBAANtaYKDDtR5d1478iPkCrZr2xnZOpMAAgIBAAJWnb0KTuJsgctA5P84BA"

func (h *Handler) MM(c tele.Context) error {
	sticker := &tele.Sticker{File: tele.File{FileID: mmSticker}}
	stickerMsg, err := h.Bot.Send(c.Chat(), sticker)
	if err != nil {
		h.Log.Error("send sticker", sl.Err(err))
	}

	result, err := h.Matching.RunMatch(context.Background())
	if err != nil {
		if stickerMsg != nil {
			_ = h.Bot.Delete(stickerMsg)
		}
		h.Log.Error("run match", sl.Err(err))
		return c.Send(view.Msg("mm", "not_enough_users"))
	}

	if stickerMsg != nil {
		_ = h.Bot.Delete(stickerMsg)
	}

	fullInfo := ""
	if result.FullMatchCount > 0 {
		fullInfo = fmt.Sprintf("\n\nполных совпадений (без общего времени): %d", result.FullMatchCount)
	}

	if err := c.Send(view.Msgf(map[string]string{
		"pairs":     fmt.Sprintf("%d", result.PairsCount),
		"users":     fmt.Sprintf("%d", result.UsersCount),
		"full_info": fullInfo,
	}, "mm", "matched")); err != nil {
		h.Log.Error("send match result", sl.Err(err))
	}

	meetResult, err := h.Meeting.CreateMeetings(context.Background())
	if err != nil {
		h.Log.Error("create meetings", sl.Err(err))
		if err.Error() == "no pairs" {
			return c.Send(view.Msg("mm", "no_pairs"))
		}
		if err.Error() == "no places" {
			return c.Send(view.Msg("mm", "no_places"))
		}
		return nil
	}

	count := 0

	for _, m := range meetResult.Meetings {
		message := view.Msgf(map[string]string{
			"place": m.Place,
			"time":  m.Time,
		}, "meet", "notification")

		kb := view.MeetingKeyboard(fmt.Sprintf("%d", m.MeetingID))

		_, err := h.Bot.Send(&tele.User{ID: m.DillID}, message, kb)
		if err != nil {
			h.Log.Error("send meeting to dill", sl.Err(err), "telegram_id", m.DillID)
		}

		_, err = h.Bot.Send(&tele.User{ID: m.DoeID}, message, kb)
		if err != nil {
			h.Log.Error("send meeting to doe", sl.Err(err), "telegram_id", m.DoeID)
		}

		count++
	}

	for _, fm := range meetResult.FullMatches {
		dillMsg := view.Msgf(map[string]string{
			"partner_username": fm.DoeUsername,
		}, "meet", "full_match")

		doeMsg := view.Msgf(map[string]string{
			"partner_username": fm.DillUsername,
		}, "meet", "full_match")

		_, err := h.Bot.Send(&tele.User{ID: fm.DillTelegramID}, dillMsg)
		if err != nil {
			h.Log.Error("send full match to dill", sl.Err(err), "telegram_id", fm.DillTelegramID)
		}

		_, err = h.Bot.Send(&tele.User{ID: fm.DoeTelegramID}, doeMsg)
		if err != nil {
			h.Log.Error("send full match to doe", sl.Err(err), "telegram_id", fm.DoeTelegramID)
		}

		count++
	}

	return c.Send(view.Msgf(map[string]string{
		"count": fmt.Sprintf("%d", count),
	}, "mm", "success"))
}
