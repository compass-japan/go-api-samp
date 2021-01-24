FROM golang:1.15 AS builder

WORKDIR /project

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-api-samp .

FROM centos:centos7

WORKDIR /app

COPY --from=builder /project/config ./config
COPY --from=builder /project/go-api-samp ./app

ENTRYPOINT ["./app"]
