FROM golang:1.11rc2 as builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

FROM alpine:latest
COPY --from=builder /src/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]
