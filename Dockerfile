FROM golang:1.19-alpine

WORKDIR /app
COPY go.mod *.go /app

RUN go build -o server . 

CMD [ "/app/server" ]
