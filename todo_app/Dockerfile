FROM golang

WORKDIR /app
COPY . .
RUN go mod download
RUN go build

RUN CGO_ENABLED=1 GOOS=linux go build -o tic_tac_toe.fun

CMD ["./todo_app"]
