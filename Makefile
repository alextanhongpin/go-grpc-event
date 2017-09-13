SERVER_PORT=:8081
GATEWAY_PORT=:9090

build-server:
	docker build -f server/Dockerfile -t alextanhongpin/go-grpc-event-server .

build-gateway:
	docker build -f server/Dockerfile -t alextanhongpin/go-grpc-event-gateway .

# stub generates the grpc server file from the proto file
stub:
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:. \
	proto/*.proto

# gw generates the grpc gateway file from the proto file
gw:
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. \
	proto/*.proto

# install-tag downloads the binary that does the tag injection
install-tag:
	go get github.com/favadi/protoc-go-inject-tag

# tag generates the inline tag for structs
tag:
	find . -name "*.pb.go" | xargs protoc-go-inject-tag -input

# serve-grpc serves the grpc server at the specified port
serve-grpc:
	go run server/main.go -port ${SERVER_PORT}

# serve-gateway serves the grpc gateway at the specified port 
# while listening to the server endpoint
serve-gateway:
	go run gateway/main.go -port ${GATEWAY_PORT} -addr localhost${SERVER_PORT}

post:
	curl -X POST -d '{"data": {"name":"hello", "start_date": 1505254724582, "uri": "test_uri"} }' http://localhost:9090/v1/events

get-all:
	curl http://localhost:9090/v1/events

get-one:
	curl http://localhost:9090/v1/events/59b6bbc47270e76be612ee81

delete:
	curl -X DELETE http://localhost:9090/v1/events/59b6babf7270e76be612ee5b

update:
	curl -X PATCH -d '{"data": {"name": "123", "uri": "hellp", "tags": ["1", "2"]}}' http://localhost:9090/v1/events/59b6baf47270e76be612ee61 

update-2:
	curl -X PATCH -d '{"data": {"uri": "cool kid"}}' http://localhost:9090/v1/events/59b6baf47270e76be612ee61 