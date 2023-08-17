FROM golang:1.20-alpine as base

ARG DATABASE_DRIVER
ARG DATABASE_URL

FROM base as dev 

WORKDIR /opt/app/server

RUN go install github.com/cosmtrek/air@latest 
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

ENV GOOSE_DRIVER=${DATABASE_DRIVER}
ENV GOOSE_DBSTRING=${DATABASE_URL}

COPY ./go.mod ./go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
