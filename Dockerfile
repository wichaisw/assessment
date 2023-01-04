FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v

RUN go build -o ./out/go-expense-app .

# ====================

FROM alpine:3.16.2
COPY --from=build-base /app/out/go-expense-app /app/go-expense-app
EXPOSE 2565

CMD ["/app/go-expense-app"]
