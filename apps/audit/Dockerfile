FROM golang:1.25.6-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /out/audit-worker ./apps/audit/cmd/worker

FROM alpine:3.22

RUN addgroup -S app && adduser -S app -G app
USER app

COPY --from=build /out/audit-worker /audit-worker

EXPOSE 8081

ENTRYPOINT ["/audit-worker"]
