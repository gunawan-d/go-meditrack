# Stage 1: Build the Go application 
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
## Build with arch amd64
ARG DIR
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o service *.go

# Stage 2: for production
FROM alpine:3.20 AS production
RUN apk --no-cache add ca-certificates && update-ca-certificates
RUN apk add --no-cache tzdata 
WORKDIR /app
COPY --from=builder /app/service .
RUN chmod +x service 
ENV TZ=Asia/Jakarta
EXPOSE 8080
ENTRYPOINT [ "./service" ]