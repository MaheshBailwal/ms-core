FROM golang:1.19

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

RUN go build -o /ms-core

EXPOSE 8080

CMD [ "/ms-core" ]