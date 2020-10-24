package cloudstatus

import "context"

// Status types
type Status string

var (
	// StatusAvailable available status
	StatusAvailable Status = "available"
	// StatusServiceDisruption disruption status
	StatusServiceDisruption Status = "service_disruption"
	// StatusServiceOutage outage status
	StatusServiceOutage Status = "service_outage"
)

// Result cloud status structure
type Result struct {
	Name        string
	IsHealthy   bool
	Status      Status
	Title       string
	Description string
	MoreInfoURL string
}

// StatusOp interface
type StatusOp interface {
	CheckContext(ctx context.Context) (*Result, error)
}
