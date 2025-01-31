package repository

import (
	"context"
	"database/sql"
	"depart/internal/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type ReportPgRepository struct {
	dbPool *pgxpool.Pool
}

func NewReportPgRepository(dpPool *pgxpool.Pool) *ReportPgRepository {
	return &ReportPgRepository{
		dbPool: dpPool,
	}
}

func (r *ReportPgRepository) Create(ctx context.Context, petition_id uuid.UUID, location string,
	description_petition string) (uuid.UUID, error) {
	tx, err := r.dbPool.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start transaction")
		return uuid.Nil, err
	}

	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("Rolling back transaction")
			_ = tx.Rollback(ctx)
		}
	}()

	var reportID uuid.UUID
	query := `
		INSERT INTO reports (petition_id, location, description_petition)
		VALUES ($1, $2, $3)
		RETURNING report_id
	`
	err = tx.QueryRow(ctx, query, petition_id, location, description_petition).Scan(&reportID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute query")
		return uuid.Nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to commit transaction")
		return uuid.Nil, err
	}

	return reportID, nil
}

func (r *ReportPgRepository) Update(ctx context.Context, petition_id uuid.UUID, report *entities.Report) (bool, error) {
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
	UPDATE reports
	SET
		content_job = COALESCE($1, content_job),
		done_at = COALESCE($2, done_at)
	WHERE petition_id = $3;
	`

	result, err := tx.Exec(ctx, query,
		emptyStringToNull(report.ContentJob),
		timeToNull(report.DoneAt),
		petition_id)
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

func timeToNull(value sql.NullTime) interface{} {
	if !value.Valid {
		return nil
	}
	return value.Time
}
