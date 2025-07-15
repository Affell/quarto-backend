FROM golang:1.24.1-alpine3.21 as builder

WORKDIR /app

WORKDIR /builder
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /builder/build

FROM alpine:3.21
    
WORKDIR /app
COPY --from=builder /builder/build /app/app
ENTRYPOINT ["/app/app"]