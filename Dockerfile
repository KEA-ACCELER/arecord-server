FROM golang:alpine

WORKDIR /Backend

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN go build -o /run-server

WORKDIR /
COPY .env .
VOLUME [ "/data" ]

EXPOSE 18080

CMD [ "/run-server" ]