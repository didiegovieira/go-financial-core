FROM golang:1.25.6-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/financial-core-api ./apps/core/cmd/api

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=build /out/financial-core-api /financial-core-api

EXPOSE 8080

ENTRYPOINT ["/financial-core-api"]
