FROM golang:1.21.6-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY main.go .

RUN CGO_ENABLED=0 go build -o /main

FROM alpine

COPY --from=builder /main /main

ENTRYPOINT ["sh", "-c", "echo $INPUT_GITHUB_TOKEN"]
