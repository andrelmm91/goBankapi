FROM golang:1.22.4-alpine3.20

WORKDIR /app

COPY . .

RUN go build -o main main.go

EXPOSE 8080

CMD [ "/app/main" ]