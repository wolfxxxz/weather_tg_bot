FROM golang:1.20-alpine

#COPY . /app_bot
#WORKDIR /app_bot
COPY . /weather_bot
WORKDIR /weather_bot

RUN go mod tidy
RUN go build -o main cmd/main.go
CMD [ "./main" ]


