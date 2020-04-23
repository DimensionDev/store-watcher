FROM golang:alpine AS builder
LABEL stage=builder
RUN apk add --no-cache git
WORKDIR /workspace
COPY cmd/watcher .
RUN CGO_ENABLED=0 GOOS=linux go build -o watcher -v .

FROM alpine AS final
WORKDIR /
COPY --from=builder /workspace/watcher .
CMD ["./watcher"]