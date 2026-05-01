package rating

import "github.com/josimar/ranking/backend/pkg/apperror"

// ErrIncompleteRatings is returned when not all active attributes have scores.
var ErrIncompleteRatings = apperror.NewBadRequestError(apperror.IncompleteRatings, "all active attributes must have a score")

// ErrInvalidScore is returned when a score is outside the 0–100 range.
var ErrInvalidScore = apperror.NewBadRequestError(apperror.InvalidScore, "score must be between 0 and 100")
