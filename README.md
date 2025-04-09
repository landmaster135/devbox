# go_dev_template
![Go](https://img.shields.io/badge/Go-1.23-%2300ADD8?logo=go)
![Coverage](https://img.shields.io/badge/Coverage-61.5%25-yellow)
![License](https://img.shields.io/badge/license-MIT-blue)

# Cloud Functions

## Entry points
- any_function

## Example
Request
```go
{
  "any": "ANY"
}
```
Example of response
```
{"data": "any_data"}
```

# Usage
- Go 1.23.5 or later

# Deployment
Set env-vars.
```bash
gcloud functions deploy any-function \
  --set-env-vars SCRIPT_MANAGER_API_CLIENT_ID=AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,SCRIPT_MANAGER_API_CLIENT_SECRET=BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB,SCRIPT_MANAGER_API_ENDPOINT=https://CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC
```

# Development

## Tools
Refer `go.mod` and etc..

## Installing
go version
```bash
go version
go list -m -u all
go mod init example.com/mymodule
go mod tidy
```

## Environment settings
Sets if deployed on cloud function
```bash
export SCRIPT_MANAGER_API_CLIENT_ID='AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA'
export SCRIPT_MANAGER_API_CLIENT_SECRET='BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB'
export SCRIPT_MANAGER_API_ENDPOINT='https://CCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCCC'
```

# Contributing
Contributions to the project are welcome. Please fork the repository and submit a pull request with your changes.

# License
MIT License
