FROM golang:1.17
WORKDIR /app/consumer
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rabbitConsumer ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/consumer ./
EXPOSE 88
CMD ["./rabbitConsumer"]