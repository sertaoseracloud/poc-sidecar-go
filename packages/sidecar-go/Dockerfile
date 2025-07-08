FROM golang:1.22 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o sidecar-proxy

FROM debian:12-slim
WORKDIR /
COPY --from=build /app/sidecar-proxy /sidecar-proxy
EXPOSE 3000 50051
ENTRYPOINT ["/sidecar-proxy"]
