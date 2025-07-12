FROM golang:1.24-alpine AS deps

WORKDIR /app

RUN CGO_ENABLED=0 go install github.com/cosmtrek/air@v1.49.0

FROM golang:1.24-alpine AS final

WORKDIR /app

COPY --from=deps /go/bin /go/bin

COPY go.mod go.sum ./
COPY air.toml ./

RUN go mod tidy

COPY . .

ENTRYPOINT ["air", "-c", "air.toml"]