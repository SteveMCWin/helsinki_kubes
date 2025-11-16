FROM golang

WORKDIR /usr/src/app
COPY . .
RUN go mod download
RUN go build

RUN CGO_ENABLED=1 GOOS=linux go build -o todo_app

CMD ["./todo_app"]
