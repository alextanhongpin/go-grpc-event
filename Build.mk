VERSION=0.0.2

# Your dockerhub name
ORGANIZATION=alextanhongpin

# Your docker image name
SERVER_IMAGE=event-server
GATEWAY_IMAGE=event-gateway

# The folder path where the files resides
SERVER_PATH=pkg/server
GATEWAY_PATH=pkg/gateway

# The constructed docker image with organization name
SERVER_REPO=${ORGANIZATION}/${SERVER_IMAGE}
GATEWAY_REPO=${ORGANIZATION}/${GATEWAY_IMAGE}

# Docker multi-stage production build
prod-server:
	docker build -f ${SERVER_PATH}/Dockerfile -t ${SERVER_REPO}:${VERSION} .
	docker tag ${SERVER_REPO}:${VERSION} ${SERVER_REPO}:latest
	docker push ${SERVER_REPO}:${VERSION}
	docker push ${SERVER_REPO}:latest

prod-gateway:
	docker build -f ${GATEWAY_PATH}/Dockerfile -t ${GATEWAY_REPO}:${VERSION} .
	docker tag ${GATEWAY_REPO}:${VERSION} ${GATEWAY_REPO}:latest
	docker push ${GATEWAY_REPO}:${VERSION}
	docker push ${GATEWAY_REPO}:latest

# Local build
dev-gateway:
	@echo "Building development gateway"
	cd ${GATEWAY_PATH} && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	rm -rf app && \
	docker build -f Dockerfile.dev -t ${GATEWAY_REPO}:${VERSION} . && \
	docker tag ${GATEWAY_REPO}:${VERSION} ${GATEWAY_REPO}:latest && \
	docker push ${GATEWAY_REPO}:${VERSION} && \
	docker push ${GATEWAY_REPO}:latest

dev-server:
	@echo "Building development server"
	cd ${SERVER_PATH} && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	rm -rf app && \
	docker build -f Dockerfile.dev -t ${SERVER_REPO}:${VERSION} . && \
	docker tag ${SERVER_REPO}:${VERSION} ${SERVER_REPO}:latest && \
	docker push ${SERVER_REPO}:${VERSION} && \
	docker push ${SERVER_REPO}:latest