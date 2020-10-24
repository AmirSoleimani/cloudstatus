package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/amirsoleimani/cloudstatus"
	"github.com/amirsoleimani/cloudstatus/pkg/circleci"
	"github.com/amirsoleimani/cloudstatus/pkg/cloudflare"
	"github.com/amirsoleimani/cloudstatus/pkg/gcp"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var snsTopicARN = os.Getenv("SNS_TOPIC_ARN")
var svc *sns.SNS

var healthyState = 1
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

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file. (~/.aws/credentials).
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = sns.New(sess)
}

func main() {
	lambda.Start(Handler)
}

// Handler Lambda handler
func Handler(ctx context.Context) {

	unhealthyServices := make([]cloudstatus.Result, 0)
	for _, status := range statuses {
		st, err := status.CheckContext(ctx)
		if err != nil {
			log.Println(err)
			continue // ?
		}
		log.Printf("%s: %v\n", st.Name, st)

		if !st.IsHealthy {
			unhealthyServices = append(unhealthyServices, *st)
		}
	}

	var message string

	// * WE DONT SEND DUPLICATE MESSAGE IF HEALTH STATE HAS NOT CHANGED
	if len(unhealthyServices) > 0 && healthyState == 1 {
		log.Println("Unhealthy")
		// Send a message to the SNS topic (unhealthy)
		healthyState = 0 // change to unhealthy
		message = snsMessageServiceIsUnhealthy(unhealthyServices)
	} else if len(unhealthyServices) == 0 && healthyState == 0 {
		log.Println("Healthy")
		// Send a message to the SNS topic (healthy)
		healthyState = 1 // change to healthy
		message = snsMessageServiceIsBack()
	}

	if message != "" {
		_, err := svc.Publish(&sns.PublishInput{
			Message:  &message,
			TopicArn: &snsTopicARN,
		})
		if err != nil {
			log.Println(err.Error())
		}
	}

	log.Println("Good Bye!")
}

func snsMessageServiceIsBack() string {
	return fmt.Sprintf(`
Hello,

All services are back and status is operational and available.

Enjoy!`)
}

func snsMessageServiceIsUnhealthy(services []cloudstatus.Result) string {
	var srv string
	for _, s := range services {
		if s.Description == "" {
			s.Description = "-"
		}
		srv = fmt.Sprintf("%s\n  - %s (%s)\n     Description: %s\n", srv, s.Name, s.MoreInfoURL, s.Description)
	}
	return fmt.Sprintf(`
OMG! Check it out

Unhealthy services: %s
`, srv)
}
