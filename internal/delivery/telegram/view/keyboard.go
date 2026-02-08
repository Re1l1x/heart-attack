package view

import (
	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/domain"
	tele "gopkg.in/telebot.v3"
)

func SexKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	male := menu.Data(messages.M.UI.Buttons.Sex.Male, "sex_male")
	female := menu.Data(messages.M.UI.Buttons.Sex.Female, "sex_female")

	menu.Inline(
		menu.Row(male, female),
	)
	return menu
}

func TimeKeyboard(selected map[string]bool) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	var rows []tele.Row
	for i := 0; i < len(domain.TimeRanges); i += 2 {
		var row tele.Row
		for _, tr := range domain.TimeRanges[i:min(i+2, len(domain.TimeRanges))] {
			text := tr
			if selected[tr] {
				text = "> " + tr + " <"
			}
			btn := menu.Data(text, "time", tr)
			row = append(row, btn)
		}
		rows = append(rows, row)
	}

	if len(selected) > 0 {
		confirm := menu.Data(messages.M.UI.Buttons.Confirm, "confirm_time")
		rows = append(rows, menu.Row(confirm))
	}

	menu.Inline(rows...)
	return menu
}

func ResubmitKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	btn := menu.Data(messages.M.UI.Buttons.Resubmit, "resubmit")
	menu.Inline(menu.Row(btn))
	return menu
}

func HowItWorksKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	btn := menu.Data(messages.M.UI.Buttons.HowItWorks, "how_it_works")
	menu.Inline(menu.Row(btn))
	return menu
}

func MeetingKeyboard(meetingID string) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	confirm := menu.Data(messages.M.UI.Buttons.ConfirmMeeting, "confirm_meeting", meetingID)
	cancel := menu.Data(messages.M.UI.Buttons.CancelMeeting, "cancel_meeting", meetingID)

	menu.Inline(menu.Row(confirm, cancel))
	return menu
}

func CancelSupportKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	btn := menu.Data(messages.M.UI.Buttons.CancelSupport, "cancel_support")
	menu.Inline(menu.Row(btn))
	return menu
}

func CancelKeyboard(meetingID string) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	cancel := menu.Data(messages.M.UI.Buttons.CancelMeeting, "cancel_meeting", meetingID)
	menu.Inline(menu.Row(cancel))
	return menu
}

func ArrivedKeyboard(meetingID string) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	arrived := menu.Data(messages.M.UI.Buttons.Arrived, "arrived_meeting", meetingID)
	menu.Inline(menu.Row(arrived))
	return menu
}

func CantFindKeyboard(meetingID string) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}
	btn := menu.Data(messages.M.UI.Buttons.CantFind, "cant_find_partner", meetingID)
	menu.Inline(menu.Row(btn))
	return menu
}
