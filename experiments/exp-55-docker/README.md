# Experiment 55: Docker Build + Deploy

## From passing Go app to running Docker container

| Step | Result |
|------|--------|
| Generate Dockerfile | OK |
| Docker build | PASS |
| Container runs | FAIL |
| Health check | OK |
| Image size | 14MB |

## Dockerfile
Multi-stage build: golang:1.22-alpine → alpine:3.19
Produces a ~20MB image with just the binary.
