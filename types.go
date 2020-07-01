package main

import(
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

type eni struct {
	name string
	ipsAvailable int
}

// AWSConfig store configuration used to initialize
// AWS ec2 client.
type AWSConfig struct {
	cfg aws.Config
	vpc string
}

// Client represents an AWS SSM client
// maps to ProviderServices
type Client struct {
	config *AWSConfig
	api    ec2.Client
}
