# aws-eni-exporter

Prometheus exporter to get how many ip availables are in each subnet in the vpc selected.

## install 

As docker:
```sh
docker build -t . 
```
or
```
docker run -ti -p 8888:8888 -v $HOME/.aws:/root/.aws -e AWS_PROFILE=your_AWS_profile -e REGION=your_AWS_region -e VPC=your_VPCID ismaelfm/eni-exporter:0.1.4
```

As binary:
```
go build -o  eni-exporter ./src/
```

As helm:
```
fill up vars in [vars](eni-exporter/values.yaml)
run helm 
```

## RUN
### Locally

You **MUST** have setted as ENVAR
- AWS_PROFILE
- REGION
- VPC

or run it as docker:
```
docker run -ti -p 8888:8888 -v $HOME/.aws:/root/.aws -e AWS_PROFILE=your_AWS_profile -e REGION=your_AWS_region -e VPC=your_VPCID ismaelfm/eni-exporter:0.1.4
```

### Kubernetes
- fill up vars in [vars](eni-exporter/values.yaml)
- run helm like:
```
helm upgrade eni-exporter . --debug --install --version 0.1.0 --values values.yaml --tls
```

## TODO
- tests

- Makefile