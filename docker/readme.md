Docker

Env vars:
1. `S3_BUCKET`: S3 bucket containing expected object key: flags.tar.gz
1. `AWS_ACCESS_KEY_ID`: AWS access key id (or omit if using instance role etc)
1. `AWS_SECRET_ACCESS_KEY`: AWS secret access key (or omit if using instance role etc)

To Build and push: (platform flag is required for AWS Lightsail)
`docker build --platform=linux/amd64 -t  drichelson/brightly:latest . && docker push drichelson/brightly:latest`

You'll want to expose port 8030.
Health check is at port `8030/status`