package repository

import (
	"github.com/scheduler-prototype/dto"
	"github.com/scheduler-prototype/utility"
)

func (r *Repository) fetchEvents(query string, args ...interface{}) ([]dto.MGraphEventDto, error) {
	rows, err := r.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []dto.MGraphEventDto
	for rows.Next() {
		var event dto.MGraphEventDto
		if err := rows.Scan(
			&event.ID,
			&event.UserId,
			&event.ICalUid,
			&event.EventId,
			&event.Title,
			&event.Description,
			&event.LocationsCount,
			&event.StartTime,
			&event.EndTime,
			&event.IsOnline,
			&event.IsAllDay,
			&event.IsCancelled,
			&event.OrganizerUserId,
			&event.CreatedTime,
			&event.UpdatedTime,
			&event.Timezone,
			&event.PlatformUrl,
			&event.MeetingUrl,
			&event.Type,
			&event.IsRecurring,
			&event.SeriesMasterId,
			&event.CreatedAt,
			&event.UpdatedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (r *Repository) CreateEvent(event *dto.MGraphEventDto) error {
	query := `
				INSERT INTO events 
					(user_id, ical_uid, event_id, title, description, 
					locations_count, start_time, end_time, is_online, 
					is_all_day, is_cancelled, organizer_user_id, 
					created_time, updated_time, timezone, platform_url, 
					meeting_url, type, is_recurring, series_master_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20 , $21, $22) 
				RETURNING id
			 `

	if err := r.conn.QueryRow(
		query,
		event.UserId,
		event.ICalUid,
		event.EventId,
		event.Title,
		event.Description,
		event.LocationsCount,
		event.StartTime,
		event.EndTime,
		event.IsOnline,
		event.IsAllDay,
		event.IsCancelled,
		event.OrganizerUserId,
		event.CreatedTime,
		event.UpdatedTime,
		event.Timezone,
		event.PlatformUrl,
		event.MeetingUrl,
		event.Type,
		event.IsRecurring,
		event.SeriesMasterId,
		event.CreatedAt,
		event.UpdatedAt,
	).Scan(&event.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetEventByICalUid(iCalUid string) (dto.MGraphEventDto, error) {
	query := `
				SELECT * FROM events WHERE ical_uid = $1
			 `

	events, err := r.fetchEvents(query, iCalUid)
	if err != nil {
		return dto.MGraphEventDto{}, err
	}

	if len(events) == 0 {
		return dto.MGraphEventDto{}, utility.ErrNotFound
	}

	return events[0], nil
}
