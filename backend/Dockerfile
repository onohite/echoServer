FROM golang:1.17-alpine
WORKDIR /app/backend
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o upserver ./cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app/backend ./
EXPOSE 80
CMD ["./upserver"]