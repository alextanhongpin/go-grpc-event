PRIVATE_SERVER_PORT=:8081
PRIVATE_GATEWAY_PORT=:9091

PUBLIC_SERVER_PORT=:8090
PUBLIC_GATEWAY_PORT=:9090

# up starts the docker 
up:
	docker-compose up -d

# down stops the docker
down:
	docker-compose down

# setup will prepare the installation of the binaries required for grpc and glide dep management
setup:
	brew install glide
	brew install protobuf
	go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	go get -u -v github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
	go get -u -v github.com/golang/protobuf/protoc-gen-go
	go get -u -v github.com/favadi/protoc-go-inject-tag
	go get github.com/gogo/protobuf/protoc-gen-gofast


# stub generates the grpc server file from the proto file
stub:
	find . -name "*.proto" -exec \
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--gofast_out=plugins=grpc:. \
	--proto_path=. "{}" \;

swagger:
	find . -name "*.proto" -exec \
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--swagger_out=logtostderr=true:. "{}" \;

	cp proto/**/**.json pkg/gateway/assets/

# gw generates the grpc gateway file from the proto file
gw:
	find . -name "*.proto" -exec \
	protoc -I/usr/local/include -I. \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--grpc-gateway_out=logtostderr=true:. "{}" \;

# tag generates the inline tag for structs
tag:
	find . -name "*.pb.go" -exec protoc-go-inject-tag -input "{}" \;

# compile contains the shortcut command to build a linux go binary
compile:
	GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -o app main.go
	
# NOTE: when building for linode, you do not need the goarch
# CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app

