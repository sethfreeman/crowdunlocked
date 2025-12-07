FROM golang:1.22-alpine AS builder
WORKDIR /app
# Copy all go.mod files for workspace
COPY go.work go.work.sum* ./
COPY services/bookings/go.mod services/bookings/go.su[m] ./services/bookings/
COPY services/releases/go.mod services/releases/go.su[m] ./services/releases/
COPY services/publicity/go.mod services/publicity/go.su[m] ./services/publicity/
COPY services/social/go.mod services/social/go.su[m] ./services/social/
COPY services/money/go.mod services/money/go.su[m] ./services/money/
RUN cd services/social && go mod download
COPY services/social ./services/social
RUN cd services/social && CGO_ENABLED=0 GOOS=linux go build -o /social ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /social .
EXPOSE 8080
CMD ["./social"]
