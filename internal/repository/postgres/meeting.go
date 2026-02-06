package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jus1d/kypidbot/internal/domain"
)

type MeetingRepo struct {
	db *sql.DB
}

func NewMeetingRepo(d *DB) *MeetingRepo {
	return &MeetingRepo{db: d.db}
}

func (r *MeetingRepo) SaveMeeting(ctx context.Context, m *domain.Meeting) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO meetings (dill_id, doe_id, pair_score, is_fullmatch)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		m.DillID, m.DoeID, m.PairScore, m.IsFullmatch,
	).Scan(&m.ID)
}

func (r *MeetingRepo) GetMeetingByID(ctx context.Context, id int64) (*domain.Meeting, error) {
	var m domain.Meeting

	err := r.db.QueryRowContext(ctx, `
		SELECT id, dill_id, doe_id, pair_score, is_fullmatch,
		       place_id, time, dill_state, doe_state
		FROM meetings WHERE id = $1`, id).Scan(
		&m.ID, &m.DillID, &m.DoeID, &m.PairScore, &m.IsFullmatch,
		&m.PlaceID, &m.Time, &m.DillState, &m.DoeState,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *MeetingRepo) GetRegularMeetings(ctx context.Context) ([]domain.Meeting, error) {
	return r.getMeetingsByFullmatch(ctx, false)
}

func (r *MeetingRepo) GetFullMeetings(ctx context.Context) ([]domain.Meeting, error) {
	return r.getMeetingsByFullmatch(ctx, true)
}

func (r *MeetingRepo) getMeetingsByFullmatch(ctx context.Context, fullmatch bool) ([]domain.Meeting, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, dill_id, doe_id, pair_score, is_fullmatch,
		       place_id, time, dill_state, doe_state
		FROM meetings WHERE is_fullmatch = $1`, fullmatch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []domain.Meeting
	for rows.Next() {
		var m domain.Meeting
		if err := rows.Scan(
			&m.ID, &m.DillID, &m.DoeID, &m.PairScore, &m.IsFullmatch,
			&m.PlaceID, &m.Time, &m.DillState, &m.DoeState,
		); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}
	return meetings, rows.Err()
}

func (r *MeetingRepo) AssignPlaceAndTime(ctx context.Context, id int64, placeID int64, time string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE meetings SET place_id = $1, time = $2 WHERE id = $3`,
		placeID, time, id)
	return err
}

func (r *MeetingRepo) UpdateState(ctx context.Context, meetingID int64, isDill bool, state domain.ConfirmationState) error {
	col := "doe_state"
	if isDill {
		col = "dill_state"
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE meetings SET `+col+` = $1 WHERE id = $2`, state, meetingID)
	return err
}

func (r *MeetingRepo) ClearMeetings(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM meetings`)
	return err
}
