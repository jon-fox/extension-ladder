# Building the binary of the App
FROM golang:1.26 AS build

WORKDIR /go/src/extension-ladder

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o extension-ladder cmd/main.go

FROM debian:12-slim as release

WORKDIR /app

COPY --from=build /go/src/extension-ladder/extension-ladder .
RUN chmod +x /app/extension-ladder

RUN apt update && apt install -y ca-certificates && rm -rf /var/lib/apt/lists/*

#EXPOSE 8080

#ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["sh", "-c", "/app/extension-ladder"]