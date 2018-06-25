# osquery-s3-config

an [osquery](https://osquery.io) config plugin to read from a configuration file stored in an [AWS s3](https://aws.amazon.com/s3/) bucket.

## building

to build the osquery extension you will need to have the following installed:
* [go](https://golang.org/)  
* [make](https://www.gnu.org/software/make/)

to build the extension use the following commands:
```
make deps

make
```

## configuration

to the run the extension the following environment variables are required to be set:
```
OSQUERY_S3_CONFIG_BUCKET_NAME
OSQUERY_S3_CONFIG_BUCKET_REGION
OSQUERY_S3_CONFIG_PATH // optional - defaults to `osquery.conf`
```

## AWS configuration
standard AWS SDK mechanisms for AWS; this includes env vars (AWS_ACCESS_KEY_ID) and profiles (AWS_PROFILE) and IAM authentication.

## troubleshooting
when troubleshooting, ensure you are running osqueryd/osqueryi with the --verbose flag.

note if running osuqery as root you will have to change the ownership of `build/osquery-s3-config.ext` to root or by passing the `--allow_unsafe` flag.

## thanks
[groob](https://twitter.com/wikiwalk) for the example in his blog post [Extending osquery with Go](https://blog.gopheracademy.com/advent-2017/osquery-sdk/).
