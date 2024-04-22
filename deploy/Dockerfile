FROM golang:1.22.2-alpine3.18 as builder

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
#COPY configs/ ./configs
#COPY resources/ ./resources

RUN go build -o ./settings-loader ./cmd/server/main.go

FROM alpine:latest as runner

COPY --from=builder /app /app

#ENV DB_HOST="postgres_container"

VOLUME ["/app/configs"]
VOLUME ["/app/resources"]

#CMD ["/bin/sh", "-c", "tail -f /dev/null"]
CMD ["/app/settings-loader"]
