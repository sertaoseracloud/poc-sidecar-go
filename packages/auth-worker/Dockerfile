FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod .
COPY ../identity-adapters ../identity-adapters
RUN go mod download
COPY . .
RUN go build -o /auth-worker

FROM gcr.io/distroless/base-debian12
COPY --from=build /auth-worker /auth-worker
CMD ["/auth-worker"]
