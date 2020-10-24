package circleci

import (
	"errors"
	"fmt"
	"testing"

	"bou.ke/monkey"
	"github.com/amirsoleimani/cloudstatus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestCheckContext(t *testing.T) {
	tests := []struct {
		name    string
		result  *cloudstatus.Result
		err     error
		apiFunc func() (*response, error)
	}{
		{
			name: "with api error",
			apiFunc: func() (*response, error) {
				return nil, errors.New("error occurred")
			},
			result: nil,
			err:    fmt.Errorf("failed to get service status: %w", errors.New("error occurred")),
		},
		{
			name: "healthy",
			apiFunc: func() (*response, error) {
				return &response{
					Status: respStatus{
						Indiactor:   "none",
						Description: "alles goed",
					},
				}, nil
			},
			result: &cloudstatus.Result{
				Name:        "CircleCI",
				MoreInfoURL: "https://status.circleci.com/",
				IsHealthy:   true,
				Status:      cloudstatus.StatusAvailable,
				Description: "alles goed",
				Title:       "All Services Available",
			},
			err: nil,
		},
		{
			name: "unhealthy",
			apiFunc: func() (*response, error) {
				return &response{
					Status: respStatus{
						Indiactor:   "outage",
						Description: "service outage",
					},
				}, nil
			},
			result: &cloudstatus.Result{
				Name:        "CircleCI",
				MoreInfoURL: "https://status.circleci.com/",
				IsHealthy:   false,
				Status:      cloudstatus.StatusServiceOutage,
				Description: "service outage",
				Title:       "Service Unavailable",
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		monkey.Patch(fetchStatus, tt.apiFunc)
		t.Run(tt.name, func(t *testing.T) {
			resp, err := New().CheckContext(context.Background())
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.result, resp)
		})
	}
}
