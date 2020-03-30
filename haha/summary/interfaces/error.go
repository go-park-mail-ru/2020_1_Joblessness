package summaryInterfaces

import "errors"

var (
	ErrSummaryAlreadySend       = errors.New("summary already send")
	ErrNoSummaryToRefresh       = errors.New("no summary to refresh")
	ErrPersonIsNotOwner       = errors.New("person doesnt own this summary")
	ErrOrgIsNotOwner       = errors.New("organization doesnt own this vacancy")
)