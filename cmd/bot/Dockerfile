FROM golang:1.25.11-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/bot ./cmd/bot

FROM alpine:3.21
RUN apk --no-cache add ca-certificates tzdata \
 && adduser -D -u 10001 app
USER app
WORKDIR /home/app
COPY --from=builder --chown=app:app /out/bot ./bot

# EXPOSE 7070
CMD ["./bot"]