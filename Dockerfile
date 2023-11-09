FROM golang:1.19-alpine as builder
LABEL authors="Rafael Leonan"
WORKDIR /app
COPY . .
COPY .env.production .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

FROM scratch
COPY .env.production .
COPY --from=builder /app/server /server
ENTRYPOINT [ "/server" ]