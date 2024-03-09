FROM golang:1.22.0-bullseye as build
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY ./cmd /app/cmd
COPY ./src /app/src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/rinha ./cmd/rinha.go

FROM debian:bullseye-slim as final
WORKDIR /app
COPY --from=build /app/bin/rinha /app
ENV GOGC 1000
# ENV GOMAXPROCS 3
CMD [ "/app/rinha" ]
