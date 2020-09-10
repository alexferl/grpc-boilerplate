FROM golang:1.15.2-alpine as builder
MAINTAINER Alexandre Ferland <aferlandqc@gmail.com>

ENV GO111MODULE=on

WORKDIR /build

RUN apk add --no-cache git

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/grpcd

FROM scratch
COPY --from=builder /build/grpcd /app
COPY --from=builder /build/configs /configs

ENTRYPOINT ["/app"]

EXPOSE 50051 50052
CMD ["--env-name", "prod"]
