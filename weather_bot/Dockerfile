FROM golang:1.20-alpine

COPY . /weather_bot
WORKDIR /weather_bot

RUN go mod tidy
RUN go build -o main cmd/main.go
CMD [ "./main" ]


