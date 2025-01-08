package schemas

import (
	"time"

	"chrono/db/repo"
)

type BatchRequest struct {
	StartDate  time.Time
	EndDate    time.Time
	EventCount int
	Request    *repo.GetPendingRequestsRow
}
