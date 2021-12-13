FROM golang:1.17
WORKDIR /app/EmailSender
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o emailSender ./main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/EmailSender ./
EXPOSE 88
CMD ["./emailSender"]