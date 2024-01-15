FROM node:alpine AS openapi
WORKDIR /usr/src/app
COPY api/ ./api
RUN npx @redocly/cli build-docs api/openapi.yaml --output openapi.html

FROM golang:1.21-alpine as builder
WORKDIR /usr/src/app
ARG GOARCH=amd64
COPY . .
RUN go mod vendor
RUN	GOOS=linux GOARCH=${GOARCH} go build -o /usr/src/app/payhere .

FROM alpine:latest

LABEL maintainer=tkddlf59@gmail.com

RUN mkdir /www
COPY --from=openapi /usr/src/app/openapi.html /www/openapi.html
COPY --from=builder /usr/src/app/payhere /usr/bin/payhere

ENTRYPOINT ["/usr/bin/payhere"]