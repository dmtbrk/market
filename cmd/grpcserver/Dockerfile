FROM golang:alpine AS build

WORKDIR /go/src/github.com/ortymid/market

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /go/bin/grpcserver ./cmd/grpcserver

FROM golang:alpine
COPY --from=build /go/bin/grpcserver /go/bin/grpcserver

CMD ["/go/bin/grpcserver"]
