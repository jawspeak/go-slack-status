package requestsresponses

import (
	"github.com/jawspeak/go-slack-status/bitbucket/models"
	"fmt"
)

type GetPullRequestActivitiesOK struct {
	Payload *models.PullRequestActivitiesResponse
}

type GetPullRequestsOK struct {
	Payload *models.PullRequestsResponse
}


// NewAPIError creates a new API error
func NewAPIError(opName string, response interface{}, code int) APIError {
	return APIError{
		OperationName: opName,
		Response:      response,
		Code:          code,
	}
}

// APIError wraps an error model and captures the status code
type APIError struct {
	OperationName string
	Response      interface{}
	Code          int
}

func (a APIError) Error() string {
	return fmt.Sprintf("%s (status %d): %+v ", a.OperationName, a.Code, a.Response)
}

