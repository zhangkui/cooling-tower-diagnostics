FROM golang:1.22-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/cooling-tower ./cmd/server

FROM golang:1.22-bookworm
WORKDIR /app
COPY --from=build /out/cooling-tower /app/cooling-tower
EXPOSE 8080
ENV DATA_DIR=/app/data HTTP_ADDR=:8080
ENTRYPOINT ["/app/cooling-tower"]
