package entities

import (
	"database/sql"

	"github.com/google/uuid"
)

type Report struct {
	ReportID             uuid.UUID `db:"report_id"`
	ContentJob           string    `db:"content_job"`
	DoneAt               sql.NullTime
	PetitionID           uuid.UUID `db:"petition_id"`
	Location             string    `db:"location"`
	Description_petition string    `db:"description_petition"`
}
