package domain

import (
	"errors"
	"strings"
	"time"
)

const (
	ApplicationAdmitted    = "admitted"
	ApplicationScreenedOut = "screened_out"
)

var (
	ErrApplicationPostingUnavailable = errors.New("posting is not accepting applications")
	ErrApplicationLimitReached       = errors.New("posting application limit reached")
	ErrDuplicateApplication          = errors.New("candidate has already applied to this posting")
	ErrCandidateUnavailable          = errors.New("candidate is not eligible to apply")
)

// Application is JOB's submission and initial Amarat screening result. Later pipeline state is owned by INVW.
type Application struct {
	ID                   string    `json:"id"`
	ReferenceNumber      int32     `json:"reference_number"`
	PostingID            string    `json:"posting_id"`
	Status               string    `json:"status"`
	SubmittedAt          time.Time `json:"submitted_at"`
	ScreenedOut          bool      `json:"screened_out"`
	ScreeningCompletedAt time.Time `json:"screening_completed_at"`
	FailedAmarat         []string  `json:"failed_amarat"`
}

// MissingAmarat returns every posting Amara that the candidate does not have.
// Matching is case-insensitive and ignores surrounding whitespace.
func MissingAmarat(required, candidate []string) []string {
	have := make(map[string]struct{}, len(candidate))
	for _, title := range candidate {
		have[strings.ToLower(strings.TrimSpace(title))] = struct{}{}
	}
	missing := make([]string, 0)
	for _, title := range required {
		if _, ok := have[strings.ToLower(strings.TrimSpace(title))]; !ok {
			missing = append(missing, title)
		}
	}
	return missing
}

// ApplicationCapReached reports whether accepted submissions filled the configured cap.
func ApplicationCapReached(maxApplications *int32, accepted int64) bool {
	return maxApplications != nil && accepted >= int64(*maxApplications)
}

// ApplicationIsOpen checks whether a published, non-archived posting accepts applications now.
func ApplicationIsOpen(status string, opening, deadline *time.Time, now time.Time) bool {
	return status == PostingPublished &&
		(opening == nil || !now.Before(*opening)) &&
		(deadline == nil || now.Before(*deadline))
}
