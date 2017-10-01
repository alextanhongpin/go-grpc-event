VERSION=0.0.2

public-server:
	docker build -f pkg/server-public/Dockerfile -t alextanhongpin/public-event-server:0.0.1-beta .
	docker push alextanhongpin/public-event-server:${VERSION}

private-server:
	docker build -f pkg/server-private/Dockerfile -t alextanhongpin/private-event-server:0.0.1-beta .
	docker push alextanhongpin/private-event-server:${VERSION}

public-gateway:
	docker build -f pkg/gateway-public/Dockerfile -t alextanhongpin/public-event-gateway:0.0.1-beta .
	docker push alextanhongpin/public-event-gateway:${VERSION}

private-gateway:
	docker build -f pkg/gateway-private/Dockerfile -t alextanhongpin/private-event-gateway:0.0.1-beta .
	docker push alextanhongpin/private-event-gateway:${VERSION}


push-all:
	docker push alextanhongpin/public-event-server:${VERSION} && \
	docker push alextanhongpin/private-event-server:${VERSION} && \
	docker push alextanhongpin/public-event-gateway:${VERSION} && \
	docker push alextanhongpin/private-event-gateway:${VERSION}


# Does a local build

local-gateway:
	@echo "building local gateway"
	cd pkg/gateway-public/ && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.production -t alextanhongpin/public-event-gateway:${VERSION} . && \
	rm -rf app

	cd pkg/gateway-private/ && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.production -t alextanhongpin/private-event-gateway:${VERSION} . && \
	rm -rf app

local:
	@echo "Building public event server"
	cd pkg/server-public/ && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.production -t alextanhongpin/public-event-server:${VERSION} . && \
	rm -rf app
	@echo "Done building public event server"

	@echo "Building privte event server"
	cd pkg/server-private/ && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.production -t alextanhongpin/private-event-server:${VERSION} . && \
	rm -rf app
	@echo "Done building private server"
