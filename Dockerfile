FROM golang:1.24 AS build

WORKDIR /app

ENV GOTOOLCHAIN=auto

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/api

FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=build /app/server .

ENV PORT=8080 \
    GIN_MODE=release

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["/app/server"]
