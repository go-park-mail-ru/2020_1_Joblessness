package vacancyInterfaces

import "errors"

var (
	ErrOrgIsNotOwner = errors.New("organization doesnt own this vacancy")
)
