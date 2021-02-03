FROM golang:1.14
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /app/app .
COPY --from=0 /app/swagger-ui .
EXPOSE 8080
CMD ["./app"]
