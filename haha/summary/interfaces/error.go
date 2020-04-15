package summaryInterfaces

import (
	"errors"
)

var (
	ErrSummaryAlreadyExists = errors.New("summary already exists")
	ErrNoSummaryToRefresh = errors.New("no summary to refresh")
	ErrPersonIsNotOwner = errors.New("person doesn't own summary")
	ErrOrganizationIsNotOwner = errors.New("organization doesn't own this vacancy")
	ErrSummaryNotFound = errors.New("summary not found")
	ErrSummaryAlreadySent = errors.New("summary already sent")
)
