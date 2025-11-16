FROM golang

WORKDIR /app
COPY . .
RUN go mod download
RUN go build

RUN CGO_ENABLED=1 GOOS=linux go build -o log_output

CMD ["./log_output", "gen"]
