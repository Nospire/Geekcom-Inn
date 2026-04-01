FROM golang:1.24-alpine AS build
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
ENV GOTOOLCHAIN=auto
RUN go mod download
COPY . .
RUN go build -o /tavrn ./cmd/tavrn-admin

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build /tavrn /usr/local/bin/tavrn
WORKDIR /app
EXPOSE 2222 8090
CMD ["tavrn"]
