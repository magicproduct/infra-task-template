FROM golang:1.21.4

WORKDIR /go/src/app

COPY . .

RUN go build -v -o dataset-cleaner cmd/dataset-cleaner/main.go

CMD ["./dataset-cleaner"]
