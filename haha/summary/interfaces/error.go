package summaryInterfaces

import "fmt"

type ErrorSummaryAlreadySent struct{}

func NewErrorSummaryAlreadySent() *ErrorSummaryAlreadySent {
	return &ErrorSummaryAlreadySent{}
}

func (e *ErrorSummaryAlreadySent) Error() string {
	return "Summary already sent"
}

type ErrorNoSummaryToRefresh struct{}

func NewErrorNoSummaryToRefresh() *ErrorNoSummaryToRefresh {
	return &ErrorNoSummaryToRefresh{}
}

func (e *ErrorNoSummaryToRefresh) Error() string {
	return "No summary to refresh"
}

type ErrorPersonIsNotOwner struct {
	PersonID  uint64
	SummaryID uint64
}

func NewErrorPersonIsNotOwner(personID, summaryID uint64) *ErrorPersonIsNotOwner {
	return &ErrorPersonIsNotOwner{
		PersonID:  personID,
		SummaryID: summaryID,
	}
}

func (e *ErrorPersonIsNotOwner) Error() string {
	return fmt.Sprintf("Person with id %d does't own summary with id %d", e.PersonID, e.SummaryID)
}

type ErrorOrganizationIsNotOwner struct{}

func NewErrorOrganizationIsNotOwner() *ErrorOrganizationIsNotOwner {
	return &ErrorOrganizationIsNotOwner{}
}

func (e *ErrorOrganizationIsNotOwner) Error() string {
	return "Organization doesn't own this vacancy"
}

type ErrorSummaryNotFound struct {
	ID uint64
}

func NewErrorSummaryNotFound(id uint64) *ErrorSummaryNotFound {
	return &ErrorSummaryNotFound{ID: id}
}

func (e *ErrorSummaryNotFound) Error() string {
	return fmt.Sprintf("Summary with id %d not found", e.ID)
}
