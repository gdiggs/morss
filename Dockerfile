FROM golang:1.18 as builder

WORKDIR /morss

COPY go.mod .
COPY go.sum .
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build -o ./app morss.go

FROM gcr.io/distroless/base-debian11

COPY --from=builder /morss/app /app

ENTRYPOINT ["/app"]