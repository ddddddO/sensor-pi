package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(test)
}

func test(ctx context.Context, snsEvent events.SNSEvent) {
	for i, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Printf("%d: %s\n", i, snsRecord.Message)
	}
}
