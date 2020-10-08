FROM golang:1.15

WORKDIR /go/src/github.com/ortymid/market

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -v -o /go/bin/httpgateway ./cmd/httpgateway

CMD ["/go/bin/httpgateway"]