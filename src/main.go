package main

import (
	"context"
	"net/http"
	"os"
	"log"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
)

var region string = os.Getenv("REGION")
var vpc string = os.Getenv("VPC")
var port string = os.Getenv("PORT")

//newAWSClient return a new AWS client
func newAWSClient(config *AWSConfig) *Client {
	c := &Client{
		config: config,
	}
	c.api = *c.NewClient(c.config)
	return c
}

//NewClient return a EC2 Client
func (c *Client) NewClient(config *AWSConfig) *ec2.Client {
	cfg, err := external.LoadDefaultAWSConfig(config.cfg)
	if err != nil {
		log.Println("failed to load config, " + err.Error())
	}
	svc := ec2.New(cfg)
	return svc

}

//GetSubnets return the subnets filtered by VPC
func (c *Client) GetSubnets(vpc string) ([]eni, error) {
	input := &ec2.DescribeSubnetsInput{
		Filters: []ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []string{
					vpc,
				},
			},
		},
	}
	req := c.api.DescribeSubnetsRequest(input)
	result, err := req.Send(context.Background())
	a := make([]eni, 0)
	for _, subnet := range result.Subnets {
		for _, v := range subnet.Tags {
			if *v.Key == "Name" {
				b := eni{name: *v.Value,
					ipsAvailable: int(*subnet.AvailableIpAddressCount)}
				a = append(a, b)
			}
		}

	}
	return a, err
}

// CreateMetrics will register AvailableIpAddressCount metric per each subnet
func CreateMetrics(eni []eni) {
	go func() {
		for _, v := range eni {
			foo := newEniCollector(v)
			prometheus.MustRegister(foo)
		}
	}()
}

//AWSconfiguration return the all config to AWSConfig struct
func AWSconfiguration(region string, vpc string) *AWSConfig {
	awsConfiguration := aws.Config(aws.Config{Region: &region})

	configuration := AWSConfig{
		cfg: awsConfiguration,
		vpc: vpc,
	}

	return &configuration
}

//main function
func main() {
	log.Print(region)
	log.Print(vpc)
	awsConfig := AWSconfiguration(region, vpc)
	cli := newAWSClient(awsConfig)
	eni, err := cli.GetSubnets(vpc)
	if err != nil {
		log.Println(err)
	}
	CreateMetrics(eni)
	log.Print("starting web server at ", port)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":" + port  , nil))
}
