package repository

import (
	"context"
	"database/sql"
	"fmt"
	"petition/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PetitionPgRepository struct {
	dbPool *pgxpool.Pool
}

func NewPetitionPgRepository(dpPool *pgxpool.Pool) *PetitionPgRepository {
	return &PetitionPgRepository{
		dbPool: dpPool,
	}
}

func (r *PetitionPgRepository) Create(ctx context.Context, p *entities.Petition) (string, error) {
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return "", err
	}

	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("Rolling back transaction")
			_ = tx.Rollback(ctx)
		}
	}()

	var petitionID string
	query := `
		INSERT INTO petitions (location, description, status)
		VALUES ($1, $2, $3)
		RETURNING petition_id
	`
	err = tx.QueryRow(ctx, query, p.Location, p.Description, "created").Scan(&petitionID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return "", err
	}

	return petitionID, nil
}

func (r *PetitionPgRepository) Update(ctx context.Context, petition_id string, p *entities.Petition) (bool, error) {
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
	}
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("Rolling back transaction")
			_ = tx.Rollback(ctx)
		}
	}()

	query := `
	UPDATE petitions
	SET
		location = COALESCE($1, location),
		description = COALESCE($2, description),
		status = COALESCE($3, status),
		department = COALESCE($4, department),
		done_at = COALESCE($5, done_at),
		report_id = COALESCE($6, report_id),
		content_job = COALESCE($7, content_job)
	WHERE petition_id = $8;
	`
	id, err := uuid.Parse(petition_id)
	if err != nil {
		return false, err
	}

	if p.ReportID != uuid.Nil {
		p.Status = "in progress"
	}

	result, err := tx.Exec(ctx, query,
		emptyStringToNull(p.Location),
		emptyStringToNull(p.Description),
		emptyStringToNull(p.Status),
		emptyStringToNull(p.Department),
		timeToNull(p.DoneAt),
		uuidToNull(p.ReportID),
		emptyStringToNull(p.ContentJob),
		id)
	fmt.Println(p.DoneAt)
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		fmt.Println(
			"Update failed: no rows were affected",
		)
		return false, nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Commit transaction failed")
		return false, err
	}
	return true, nil
}

func (r *PetitionPgRepository) Get(ctx context.Context, petition_id string) (*entities.Petition, error) {
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
	}
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("Rolling back transaction")
			_ = tx.Rollback(ctx)
		}
	}()

	petition := &entities.Petition{}
	query := `
	SELECT petition_id, location, description, created_at, status, department, done_at, report_id, content_job 
		FROM petitions 
		WHERE petition_id = $1;
	`

	err = tx.QueryRow(ctx, query, petition_id).Scan(
		&petition.PetitionID,
		&petition.Location,
		&petition.Description,
		&petition.CreatedAt,
		&petition.Status,
		&petition.Department,
		&petition.DoneAt,
		&petition.ReportID,
		&petition.ContentJob,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Warn().Str("petition_id", petition_id).Msg("No petition found for petition_id")
			return nil, nil
		}
		log.Error().Err(err).Msg("Failed to execute query")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Commit transaction failed")
		return nil, err
	}

	return petition, nil
}

func (r *PetitionPgRepository) Delete(ctx context.Context, petition_id string) (bool, error) {
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to start transaction")
	}
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("Rolling back transaction")
			_ = tx.Rollback(ctx)
		}
	}()

	query := `
	DELETE FROM petitions
	WHERE petition_id = $1;
	`

	result, err := tx.Exec(ctx, query, petition_id)
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return false, nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Commit transaction failed")
		return false, err
	}

	return true, nil
}

func emptyStringToNull(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func uuidToNull(value uuid.UUID) interface{} {
	if value == uuid.Nil {
		return nil
	}
	return value
}

func timeToNull(value sql.NullTime) interface{} {
	if !value.Valid {
		return nil
	}
	return value.Time
}
