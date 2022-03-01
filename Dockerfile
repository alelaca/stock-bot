FROM golang:1.16-alpine

WORKDIR /stock-bot

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o /build/stock-bot.go src/main.go

EXPOSE 8080

CMD [ "/build/stock-bot.go" ]