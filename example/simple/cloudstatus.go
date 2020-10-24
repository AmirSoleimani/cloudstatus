package main

import (
	"context"
	"log"
	"time"

	"github.com/amirsoleimani/cloudstatus"
	"github.com/amirsoleimani/cloudstatus/pkg/circleci"
	"github.com/amirsoleimani/cloudstatus/pkg/cloudflare"
	"github.com/amirsoleimani/cloudstatus/pkg/gcp"
)

var statuses = []cloudstatus.StatusOp{}

func init() {
	// CircleCI
	circleciStatus := circleci.New()
	statuses = append(statuses, circleciStatus)

	// GCP
	gcpStatus := gcp.New()
	statuses = append(statuses, gcpStatus)

	// Cloudflare
	cloudflareStatus := cloudflare.New()
	statuses = append(statuses, cloudflareStatus)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, status := range statuses {
		st, err := status.CheckContext(ctx)
		if err != nil {
			log.Println(err)
		}
		log.Printf("%s: %v\n", st.Name, st)
	}
}
