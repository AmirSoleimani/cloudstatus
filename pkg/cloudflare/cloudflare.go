package cloudflare

import (
	"context"
	"fmt"

	"github.com/amirsoleimani/cloudstatus"
)

// Operator interface
type Operator interface {
	cloudstatus.StatusOp
}

// cloudflare structure
type cloudflare struct {
	name    string
	pageURL string
}

// New create new `cloudflare` status operator
func New() Operator {
	return &cloudflare{
		name:    "Cloudflare",
		pageURL: "https://www.cloudflarestatus.com/",
	}
}

func (g *cloudflare) CheckContext(ctx context.Context) (*cloudstatus.Result, error) {

	statusResp, err := fetchStatus()
	if err != nil {
		return nil, fmt.Errorf("failed to get service status: %w", err)
	}

	result := new(cloudstatus.Result)
	result.Name = g.name
	result.MoreInfoURL = g.pageURL

	if !statusResp {
		// TODO use the last status description
		result.Title = "Service Unavailable"
		result.Status = cloudstatus.StatusServiceOutage
		return result, nil
	}

	result.IsHealthy = true
	result.Title = "All Services Available"
	result.Status = cloudstatus.StatusAvailable
	return result, nil
}
