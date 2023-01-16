// Copyright 2023, e-inwork.com. All rights reserved.

package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

// Team
type Team struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	TeamUser    uuid.UUID
	TeamName    string
	TeamPicture string
	IsIndexed   bool
	IsDeleted   bool
	Version     int
}

type TeamModel struct {
	DB *sql.DB
}

// Get a Team by the ID field from the database,
// and convert the record to the Team.
func (m TeamModel) Get(id uuid.UUID) (*Team, error) {
	query := `
        SELECT id, created_at, team_user, team_name,
					team_picture, is_indexed, is_deleted, version
        FROM teams
        WHERE id = $1`

	var team Team

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&team.ID,
		&team.CreatedAt,
		&team.TeamUser,
		&team.TeamName,
		&team.TeamPicture,
		&team.IsIndexed,
		&team.IsDeleted,
		&team.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &team, nil
}

func (m TeamModel) IsIndexedTrue(team *Team) error {
	// SQL Update
	query := `
        UPDATE teams
        SET is_indexed = true
				WHERE id = $1 AND version = $2
        RETURNING version`

	// Assign arguments
	args := []interface{}{
		team.ID,
		team.Version,
	}

	// Create a context of the SQL Update
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Run SQL Update
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&team.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m TeamModel) Delete(team *Team) error {
	// SQL Update
	query := `
        DELETE FROM teams
				WHERE id = $1`

	// Assign arguments
	args := []interface{}{
		team.ID,
	}

	// Create a context of the SQL Update
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Run SQL Update
	err := m.DB.QueryRowContext(ctx, query, args...).Scan()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
