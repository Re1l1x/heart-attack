package domain

import "context"

type ConfirmationState string

const (
	StateNotConfirmed ConfirmationState = "not_confirmed"
	StateConfirmed    ConfirmationState = "confirmed"
	StateCancelled    ConfirmationState = "cancelled"
)

type Meeting struct {
	ID          int64
	DillID      int64
	DoeID       int64
	PairScore   float64
	IsFullmatch bool
	PlaceID     *int64
	Time        *string
	DillState   ConfirmationState
	DoeState    ConfirmationState
}

type MeetingRepository interface {
	SaveMeeting(ctx context.Context, m *Meeting) error
	GetMeetingByID(ctx context.Context, id int64) (*Meeting, error)
	GetRegularMeetings(ctx context.Context) ([]Meeting, error)
	GetFullMeetings(ctx context.Context) ([]Meeting, error)
	AssignPlaceAndTime(ctx context.Context, id int64, placeID int64, time string) error
	UpdateState(ctx context.Context, meetingID int64, isDill bool, state ConfirmationState) error
	ClearMeetings(ctx context.Context) error
}
