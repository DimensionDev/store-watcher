FROM golang:alpine AS builder
LABEL stage=builder
RUN apk add --no-cache git
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o watcher -v ./cmd/watcher

FROM python:3.7-alpine AS final
WORKDIR /
RUN pip install requests
COPY --from=builder /workspace/watcher .
CMD ["./watcher"]