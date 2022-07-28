# builder
FROM golang:1.18-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /slashbot .

# alpine
FROM ffmpeg-alpine

COPY --from=builder /slashbot /slashbot

ENTRYPOINT ["/slashbot"]
CMD ["--help"]
