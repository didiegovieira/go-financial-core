FROM golang:1.25.6-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/payments-api ./apps/payments/cmd/api

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=build /out/payments-api /payments-api

EXPOSE 8080

ENTRYPOINT ["/payments-api"]
