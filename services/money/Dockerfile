FROM golang:1.22-alpine AS builder
WORKDIR /app
# Copy all go.mod files for workspace
COPY go.work go.work.sum* ./
COPY services/bookings/go.mod services/bookings/go.su[m] ./services/bookings/
COPY services/releases/go.mod services/releases/go.su[m] ./services/releases/
COPY services/publicity/go.mod services/publicity/go.su[m] ./services/publicity/
COPY services/social/go.mod services/social/go.su[m] ./services/social/
COPY services/money/go.mod services/money/go.su[m] ./services/money/
RUN cd services/money && go mod download
COPY services/money ./services/money
RUN cd services/money && CGO_ENABLED=0 GOOS=linux go build -o /money ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /money .
EXPOSE 8080
CMD ["./money"]
