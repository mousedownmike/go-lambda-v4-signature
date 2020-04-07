package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	S3Endpoint = "s3.us-east-1.amazonaws.com"
	S3Region   = "us-east-1"
)

func main() {
	awsSess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		log.Fatalf("failed creating session: %s", err)
	}
	signer := v4.NewSigner(awsSess.Config.Credentials)

	bl := &BucketLister{
		endpoint: S3Endpoint,
		region:   S3Region,
		signer:   signer,
		client:    &http.Client{},
	}

	// If this environment variable is set, assume we're running in a Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		lambda.Start(handleList(bl))
	} else {
		handleList(bl)()
	}
}

func handleList(bl *BucketLister) func() {
	return func() {
		err := bl.List()
		if err != nil {
			log.Fatalf("failed: %s", err)
		}
	}
}

type BucketLister struct {
	endpoint string
	region   string
	signer   *v4.Signer
	client   *http.Client
}

/**
Log all the buckets available to the assumed role.
*/
func (b *BucketLister) List() error {
	url := fmt.Sprintf("https://%s/", b.endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request, %w", err)
	}

	_, err = b.signer.Sign(req, nil, "s3", b.region, time.Now())
	if err != nil {
		return fmt.Errorf("error signing request: %w", err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resposne: %w", err)
	}
	log.Printf("buckets: %s", string(body))
	return nil
}
