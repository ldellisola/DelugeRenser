FROM golang:1.22-alpine as builder
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o ./out/DelugeRenser

FROM alpine:latest
WORKDIR /app
COPY --from=builder /src/out/ /app/
CMD ["./DelugeRenser"]