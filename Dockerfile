FROM golang:1.24.5 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -o /bin/app -ldflags="-s -w" cmd/server/main.go

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=builder /bin/app /app/app

USER nonroot

ENTRYPOINT ["/app/app"]