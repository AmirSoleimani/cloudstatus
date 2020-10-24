package cloudflare

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
		apiFunc func() (bool, error)
	}{
		{
			name: "with api error",
			apiFunc: func() (bool, error) {
				return false, errors.New("error occurred")
			},
			result: nil,
			err:    fmt.Errorf("failed to get service status: %w", errors.New("error occurred")),
		},
		{
			name: "healthy",
			apiFunc: func() (bool, error) {
				return true, nil
			},
			result: &cloudstatus.Result{
				Name:        "Cloudflare",
				MoreInfoURL: "https://www.cloudflarestatus.com/",
				IsHealthy:   true,
				Status:      cloudstatus.StatusAvailable,
				Title:       "All Services Available",
			},
			err: nil,
		},
		{
			name: "unhealthy",
			apiFunc: func() (bool, error) {
				return false, nil
			},
			result: &cloudstatus.Result{
				Name:        "Cloudflare",
				MoreInfoURL: "https://www.cloudflarestatus.com/",
				IsHealthy:   false,
				Status:      cloudstatus.StatusServiceOutage,
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
