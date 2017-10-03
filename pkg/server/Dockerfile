FROM golang:1.9 as builder

MAINTAINER Alex Hong Pin Tan <alextan220990@gmail.com>

LABEL version="0.0.1-beta"
LABEL description="A grpc-server serving private events data"
LABEL language="golang:1.9"
LABEL keywords="go, grpc, grpc-server, event"

WORKDIR /go/src/github.com/alextanhongpin/go-grpc-event/

COPY . .

WORKDIR /go/src/github.com/alextanhongpin/go-grpc-event/pkg/server-private

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Multi-stage Build for Production Image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/src/github.com/alextanhongpin/go-grpc-event/pkg/server-private/app .
CMD ["./app"]