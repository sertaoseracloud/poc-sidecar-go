FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod .
COPY ../sidecar-go ../sidecar-go
RUN go mod download
COPY . .
RUN go build -o /mock-app

FROM gcr.io/distroless/base-debian12
COPY --from=build /mock-app /mock-app
CMD ["/mock-app"]
