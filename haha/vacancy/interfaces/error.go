package vacancyInterfaces

import "errors"

var (
	ErrOrgIsNotOwner = errors.New("organization doesn't own this vacancy")
)
