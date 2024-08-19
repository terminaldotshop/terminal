FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ssh ./cmd/ssh/main.go

# Second stage: build the runtime image.
FROM alpine
WORKDIR /root/
COPY --from=builder /app/ssh .
CMD TERM=xterm-256color ./ssh
