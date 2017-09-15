FROM golang:1.9 as builder

WORKDIR /go/src/github.com/alextanhongpin/go-grpc-event

COPY . .

WORKDIR /go/src/github.com/alextanhongpin/go-grpc-event/server

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /go/src/github.com/alextanhongpin/go-grpc-event/server/app .
CMD ["./app"]