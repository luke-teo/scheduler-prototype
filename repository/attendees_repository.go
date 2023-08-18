package repository

import (
	"github.com/scheduler-prototype/dto"
	"github.com/scheduler-prototype/utility"
)

func (r *Repository) fetchAttendees(query string, args ...interface{}) ([]dto.MGraphAttendeeDto, error) {
	rows, err := r.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attendees []dto.MGraphAttendeeDto
	for rows.Next() {
		var attendee dto.MGraphAttendeeDto
		if err := rows.Scan(
			&attendee.ID,
			&attendee.UserId,
			&attendee.Name,
			&attendee.EmailAddress,
			&attendee.ICalUid,
			&attendee.CreatedAt,
			&attendee.UpdatedAt,
		); err != nil {
			return nil, err
		}
		attendees = append(attendees, attendee)
	}
	return attendees, nil
}

func (r *Repository) CreateAttendee(attendee *dto.MGraphAttendeeDto) error {
	query := `
				INSERT INTO attendees 
					(user_id, name, email_address, ical_uid, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6) 
				RETURNING id
			 `

	if err := r.conn.QueryRow(
		query,
		attendee.UserId,
		attendee.Name,
		attendee.EmailAddress,
		attendee.ICalUid,
		attendee.CreatedAt,
		attendee.UpdatedAt,
	).Scan(&attendee.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAttendeeByICalUidAndEmailAddress(iCalUid string, emailAddress string) (dto.MGraphAttendeeDto, error) {
	query := `
				SELECT * FROM attendees WHERE ical_uid = $1 AND email_address = $2
			 `

	attendees, err := r.fetchAttendees(query, iCalUid, emailAddress)
	if err != nil {
		return dto.MGraphAttendeeDto{}, err
	}

	if len(attendees) == 0 {
		return dto.MGraphAttendeeDto{}, utility.ErrNotFound
	}

	return attendees[0], nil
}

