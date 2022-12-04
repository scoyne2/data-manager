FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./
COPY feed ./feed
COPY schemas ./schemas

RUN go build -o /data-manager-api

CMD [ "/data-manager-api" ]
