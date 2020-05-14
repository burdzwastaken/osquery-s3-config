# osquery-s3-config

A [osquery](https://osquery.io) config plugin to read from a configuration file stored in an [AWS s3](https://aws.amazon.com/s3/) bucket.

## Building

To build the osquery extension you will need to have the following installed:
* [go](https://golang.org/)  
* [make](https://www.gnu.org/software/make/)

To build the extension use the following commands:
```
make
```

## Configuration

To the run the extension the following environment variables are required to be set:
```
OSQUERY_S3_CONFIG_BUCKET_NAME
OSQUERY_S3_CONFIG_BUCKET_REGION
OSQUERY_S3_CONFIG_PATH // optional - defaults to `osquery.conf`
```

## AWS configuration
Standard AWS SDK mechanisms for AWS; This includes env vars (AWS_ACCESS_KEY_ID) and profiles (AWS_PROFILE) and IAM authentication.

## Troubleshooting
When troubleshooting, ensure you are running osqueryd/osqueryi with the --verbose flag.

Note if running osquery as root you will have to change the ownership of `build/osquery-s3-config.ext` to root or by passing the `--allow_unsafe` flag.

## Thanks
[groob](https://twitter.com/wikiwalk) for the example in his blog post [Extending osquery with Go](https://blog.gopheracademy.com/advent-2017/osquery-sdk/).
