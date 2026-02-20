FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk --no-cache add ca-certificates && \
  adduser -D verifier && \
  chown -R verifier /app

USER verifier

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o verifier ./cmd/api/main.go


FROM alpine:latest
WORKDIR /root/

RUN apk --no-cache add ca-certificates
COPY --from=builder /app/verifier .
EXPOSE 8080

CMD [ "./verifier" ]

