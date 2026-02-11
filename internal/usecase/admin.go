package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/domain"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrAlreadyAdmin = errors.New("user is already admin")
	ErrNotAdmin     = errors.New("user is not admin")
)

type Admin struct {
	users    domain.UserRepository
	meetings domain.MeetingRepository
}

func NewAdmin(users domain.UserRepository, meetings domain.MeetingRepository) *Admin {
	return &Admin{users: users, meetings: meetings}
}

func (a *Admin) Promote(ctx context.Context, username string) error {
	user, err := a.users.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	if user.IsAdmin {
		return ErrAlreadyAdmin
	}
	return a.users.SetAdmin(ctx, user.TelegramID, true)
}

func (a *Admin) Demote(ctx context.Context, username string) error {
	user, err := a.users.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}
	if !user.IsAdmin {
		return ErrNotAdmin
	}
	return a.users.SetAdmin(ctx, user.TelegramID, false)
}

func (a *Admin) GetStatistics(ctx context.Context) (domain.Statistics, error) {
	total, registered, optedOut, err := a.users.GetUserCounts(ctx)
	if err != nil {
		return domain.Statistics{}, err
	}

	daily, weekly, err := a.users.GetLastRegisteredCount(ctx)
	if err != nil {
		return domain.Statistics{}, err
	}

	males, females, err := a.users.GetSexCounts(ctx)
	if err != nil {
		return domain.Statistics{}, err
	}

	meetingStats, err := a.meetings.GetMeetingStats(ctx)
	if err != nil {
		return domain.Statistics{}, err
	}

	return domain.Statistics{
		TotalUsers:       total,
		RegisteredUsers:  registered,
		OptedOutUsers:    optedOut,
		RegisteredDaily:  daily,
		RegisteredWeekly: weekly,
		MaleCount:        males,
		FemaleCount:      females,
		Meetings:         meetingStats,
	}, nil
}

func (a *Admin) FormatPanel(ctx context.Context) (string, error) {
	s, err := a.GetStatistics(ctx)
	if err != nil {
		return "", err
	}

	active := s.RegisteredUsers - s.OptedOutUsers

	return messages.Format(messages.M.Command.Admin, map[string]string{
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
	}), nil
}
