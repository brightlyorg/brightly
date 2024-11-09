![go workflow](https://github.com/brightlyorg/brightly/actions/workflows/go.yml/badge.svg)
# Brightly Flags Backend Bits
[More info on the Brightly project](https://github.com/brightlyorg/brightly/wiki)

This repo contains the backend bits for the Brightly Flags project:
1. Go code that runs in GitHub Actions converting human-friendly yaml files to a format that can be consumed by the ld-relay appliance.
2. Dockerfile + friends to build the image used in the deployed backend service.

### Docker
The Dockerfile is used to build the image used in the deployed backend service. It is built on top of the ld-relay image.
To build and publish it: (requires docker login with permissions to push to drichelson)
```bash
export TAG=0.0.6 && docker build --platform=linux/amd64 -t drichelson/brightly:$TAG ./docker/ && docker push drichelson/brightly:$TAG
```
