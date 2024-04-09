FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ssh ./cmd/ssh/main.go

# Second stage: build the runtime image.
FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/ssh .
COPY id_ed25519 .
COPY id_ed25519.pub .
CMD ["./ssh"]
