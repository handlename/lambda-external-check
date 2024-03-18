# lambda-external-check

HTTP external checker running on AWS Lambda.

## Usage

1. Deploy to AWS Lambda
2. Monitor CloudWatch Metrics `AWS/Lambda Errors`

## Configuration

Needs environment variables:

- `EXTERNAL_CHECK_TARGET`: URL to check
- `EXTERNAL_CHECK_TIMEOUT`: Timeout string like `10s` or `1m`

## License

MIT

## Author

@handlename
