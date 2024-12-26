FROM golang:alpine as compiler
WORKDIR /app
COPY . .
RUN go build ./cmd/server.go

FROM alpine as server
COPY --from=compiler /app/server .
COPY --from=compiler /app/.env .
CMD ./server