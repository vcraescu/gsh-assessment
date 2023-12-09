# Assessment ![](https://github.com/vcraescu/gsh-assessment/actions/workflows/go.yml/badge.svg)

## How to run the tests

`make test`

## How to run the web server

### Local

The local web server is initiated and running on http://localhost:3000.

`make start-local`

### Docker

A docker container is started and the web server is running on http://localhost:3000.

`make start`

## How to deploy

### AWS Lambda

The lambda function is deployed using your credentials on the local environment.

`make sls-deploy`
