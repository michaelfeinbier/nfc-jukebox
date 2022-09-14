# syntax=docker/dockerfile:1
FROM node:18-alpine AS app

ADD app /app
WORKDIR /app
RUN yarn install && yarn build

FROM golang:1.18-alpine AS golang

ARG GOOS=linux
ARG GOARCH=arm64

WORKDIR /app

COPY server/go.mod ./
COPY server/go.sum ./
RUN go mod download

COPY server/ ./

#ENV GOOS=${GOOS} GOARCH=${GOARCH}

RUN go build


FROM alpine:3.16

WORKDIR /app
COPY --from=app /server/static /app/static
COPY --from=golang /app/vinyl-player /app/

ENV GIN_MODE=release

EXPOSE 8080
# HEALTHCHECK CMD curl --fail http://localhost:8080 || exit 1

ENTRYPOINT [ "/app/vinyl-player" ]
