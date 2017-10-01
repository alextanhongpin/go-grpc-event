PRIVATE_SERVER_PORT=:8081
PRIVATE_GATEWAY_PORT=:9091

PUBLIC_SERVER_PORT=:8090
PUBLIC_GATEWAY_PORT=:9090

up:
	docker-compose -f docker-compose/production/docker-compose.yml up -d

down:
	docker-compose -f docker-compose/production/docker-compose.yml down

build-server:
	docker build -f server/Dockerfile -t alextanhongpin/go-grpc-event-server .

build-gateway:
	docker build -f gateway/Dockerfile -t alextanhongpin/go-grpc-event-gateway .

# setup
setup:
	brew install glide
	brew install protobuf
	go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u -v github.com/golang/protobuf/protoc-gen-go
	go get -u -v github.com/favadi/protoc-go-inject-tag


# stub generates the grpc server file from the proto file
stub:
	find . -name "*.proto" -exec \
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--go_out=plugins=grpc:. \
	--proto_path=. "{}" \;

# proto/**/*.proto

# Setting proto path in order to import them (cannot be absolute)
# -I/--proto_path=. \

# gw generates the grpc gateway file from the proto file
gw:
	find . -name "*.proto" -exec \
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. "{}" \;
# proto/**/.proto

# tag generates the inline tag for structs
tag:
	find . -name "*.pb.go" -exec protoc-go-inject-tag -input "{}" \;


# serve-grpc serves the grpc server at the specified port
serve-private:
	go run pkg/server-private/main.go -port ${PRIVATE_SERVER_PORT}


serve-public:
	go run pkg/server-public/main.go -port ${PUBLIC_SERVER_PORT}


# serve-gateway serves the grpc gateway at the specified port while listening to the server endpoint
serve-gateway-private:
	go run pkg/gateway-private/main.go -port ${PRIVATE_GATEWAY_PORT} -addr localhost${PRIVATE_SERVER_PORT} -jwks_uri "" -auth0_aud -auth0_iss

serve-gateway-public:
	go run pkg/gateway-public/main.go -port ${PUBLIC_GATEWAY_PORT} -addr localhost${PUBLIC_SERVER_PORT}


# ===============
# API Endpoints #
# ===============

post:
	curl -X POST -d '{"data": {"name":"hello", "start_date": 1505254724582, "uri": "test_uri"} }' http://localhost:9090/v1/events


get-public-all:
	curl http://localhost:9090/public/v1/events


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

push-server:
	docker push alextanhongpin/go-grpc-event-server

push-gateway:
	docker push alextanhongpin/go-grpc-event-gateway


compile:
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o app main.go