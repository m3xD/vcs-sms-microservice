FROM golang:1.22.5-bullseye AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o app main.go

FROM alpine:3.20
RUN apk add --no-cache libc6-compat gcc libc-dev
COPY --from=build /app/app /usr/local/bin/app
RUN chmod +x /usr/local/bin/app
CMD ["app"]