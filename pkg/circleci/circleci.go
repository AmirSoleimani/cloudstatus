package circleci

import (
	"context"
	"fmt"

	"github.com/amirsoleimani/cloudstatus"
)

// Operator interface
type Operator interface {
	cloudstatus.StatusOp
}

// circleCI structure
type circleCI struct {
	name    string
	pageURL string
}

// New create new `CircleCI` status operator
func New() Operator {
	return &circleCI{
		name:    "CircleCI",
		pageURL: "https://status.circleci.com/",
	}
}

func (c *circleCI) CheckContext(ctx context.Context) (*cloudstatus.Result, error) {

	statusResp, err := fetchStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get service status: %w", err)
	}

	result := new(cloudstatus.Result)
	result.Name = c.name
	result.MoreInfoURL = c.pageURL
	result.Description = statusResp.Status.Description

	if statusResp.Status.Indiactor != "none" {
		result.Status = cloudstatus.StatusServiceOutage
		result.Title = "Service Unavailable"
		return result, nil
	}

	result.IsHealthy = true
	result.Status = cloudstatus.StatusAvailable
	result.Title = "All Services Available"
	return result, nil
}
