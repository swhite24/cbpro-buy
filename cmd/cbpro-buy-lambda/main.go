package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/swhite24/dca-bot/pkg/cmd"
)

// Handle is the handler for lambda
func Handle() error {
	return cmd.CBProBuyCmd.Execute()
}

func main() {
	lambda.Start(Handle)
}
