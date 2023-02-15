FROM golang:1.20 as teststage
WORKDIR /app
COPY ./ ./
RUN go test ./...
RUN CGO_ENABLED=0 go build -o hashd ./cmd/hashd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=teststage /app/hashd /app
CMD ["./hashd"]