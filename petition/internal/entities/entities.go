package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Petition struct {
	PetitionID  uuid.UUID    `db:"report_id"`
	Location    string       `db:"location"`
	Description string       `db:"description"`
	CreatedAt   time.Time    `db:"created_at"`
	Status      string       `db:"status"`
	Department  string       `db:"department"`
	DoneAt      sql.NullTime `db:"done_at"`
	ReportID    uuid.UUID    `db:"report_id"`
	ContentJob  string       `db:"content_job"`
}
