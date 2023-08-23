package repository

import (
	"github.com/google/uuid"
	"github.com/scheduler-prototype/dto"
	"github.com/scheduler-prototype/utility"
)

func (r *Repository) fetchUsers(query string, args ...interface{}) ([]dto.UserDto, error) {
	rows, err := r.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []dto.UserDto
	for rows.Next() {
		var user dto.UserDto
		if err := rows.Scan(
			&user.ID,
			&user.UserId,
			&user.CurrentDelta,
			&user.PreviousDelta,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *Repository) CreateUser(user *dto.UserDto) error {
	query := `
				INSERT INTO users 
					(user_id, current_delta, previous_delta, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5) 
				RETURNING id
			 `

	if err := r.conn.QueryRow(
		query,
		user.UserId,
		user.CurrentDelta,
		user.PreviousDelta,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserByUserId(userId *uuid.UUID) (dto.UserDto, error) {
	query := `
						SELECT * FROM users WHERE user_id = $1 
					`
	users, err := r.fetchUsers(query, userId)
	if err != nil {
		return dto.UserDto{}, err
	}

	if len(users) == 0 {
		return dto.UserDto{}, utility.ErrNotFound
	}

	return users[0], nil
}

func (r *Repository) UpdateCurrentDeltaByUser(userDto *dto.UserDto) error {
	query := ` 
						UPDATE users SET current_delta = $2 WHERE user_id = $1
					`

	if _, err := r.conn.Exec(
		query,
		userDto.UserId,
		*userDto.CurrentDelta,
	); err != nil {
		return err
	}

	return nil
}
