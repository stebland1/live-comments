FROM golang:1.25-aline AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
	-trimpath \
	-ldflags="-s -w" \
	-o /app/live-comments \
	./cmd/server/main.go

FROM alpine:3.20

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /app/live-comments /usr/local/bin/live-comments

USER app

EXPOSE 8081

CMD ["live-comments"]
