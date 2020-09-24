FROM golang:alpine AS build

WORKDIR /go/src/github.com/ortymid/market

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /go/bin/httpgateway ./cmd/httpgateway

FROM golang:alpine
COPY --from=build /go/bin/httpgateway /go/bin/httpgateway

CMD ["/go/bin/httpgateway"]
