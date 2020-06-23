package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	//"strings"


	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/prometheus/client_golang/prometheus"
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



func newAWSClient(config *AWSConfig) *Client {
	c := &Client{
		config: config,
	}
	c.api = *c.NewClient(c.config)
	return c
}

//NewClient daf 
func (c *Client) NewClient(config *AWSConfig) *ec2.Client{
	cfg, err := external.LoadDefaultAWSConfig(config.cfg)
	if err != nil {
		panic("failed to load config, " + err.Error())
	}
	svc := ec2.New(cfg)
	return svc

}

func (c *Client) GetSubnets() ([]eni, error){
	input := &ec2.DescribeSubnetsInput{
		Filters: []ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []string{
					"vpc-cf458dab",
				},
			},
		},
	}
	req := c.api.DescribeSubnetsRequest(input)
	result, err := req.Send(context.Background())
	a := make([]eni,0)
	for _, subnet := range result.Subnets{
		for _, v := range subnet.Tags{
			if *v.Key == "Name"{
				b:= eni{name: *v.Value, 
					ipsAvailable: int(*subnet.AvailableIpAddressCount)}
				a = append(a, b)
			}
		}
		
	}

	return a, err
}

func CreateMetrics(eni []eni)  {
	go func() {
		for _,v := range(eni){

		foo := newEniCollector(v)
		prometheus.MustRegister(foo)
		}
	}()
}

func AWSconfiguration() []eni{
	region := "us-west-1"
	awsConfiguration := aws.Config(aws.Config{Region: &region})
	
	configuration := AWSConfig{
		cfg: awsConfiguration,
		vpc: "vpc-******",
	}
	cli := newAWSClient(&configuration)
	eni, err := cli.GetSubnets()
	if err != nil {
		fmt.Println(err)
	}
	return eni
}

func main(){ 
	
	eni := AWSconfiguration()
	CreateMetrics(eni)
	
	fmt.Println("HERE")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))

}

