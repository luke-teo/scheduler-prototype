package repository

import (
	"github.com/scheduler-prototype/dto"
	"github.com/scheduler-prototype/utility"
)

func (r *Repository) fetchLocations(query string, args ...interface{}) ([]dto.MGraphLocationDto, error) {
	rows, err := r.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []dto.MGraphLocationDto
	for rows.Next() {
		var location dto.MGraphLocationDto
		if err := rows.Scan(
			&location.ID,
			&location.ICalUid,
			&location.DisplayName,
			&location.LocationUri,
			&location.Address,
			&location.CreatedAt,
			&location.UpdatedAt,
		); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}

func (r *Repository) CreateLocation(location *dto.MGraphLocationDto) error {
	query := `
				INSERT INTO locations 
					(ical_uid, display_name, location_uri, address, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6) 
				RETURNING id
			 `

	if err := r.conn.QueryRow(
		query,
		location.ICalUid,
		location.DisplayName,
		location.LocationUri,
		location.Address,
		location.CreatedAt,
		location.UpdatedAt,
	).Scan(&location.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLocationByICalUid(iCalUid string) (dto.MGraphLocationDto, error) {
	query := `
				SELECT * FROM locations WHERE ical_uid = $1
			 `

	locations, err := r.fetchLocations(query, iCalUid)
	if err != nil {
		return dto.MGraphLocationDto{}, err
	}

	if len(locations) == 0 {
		return dto.MGraphLocationDto{}, utility.ErrNotFound
	}

	return locations[0], nil
}
