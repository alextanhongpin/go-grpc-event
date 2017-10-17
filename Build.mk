GATEWAY_VERSION=1.0.2
EVENT_VERSION=1.0.2
PHOTO_VERSION=1.0.1

# Your dockerhub name
ORGANIZATION=alextanhongpin

# Your docker image name
EVENT_IMAGE=event
PHOTO_IMAGE=photo
GATEWAY_IMAGE=gateway

# The folder path where the files resides
EVENT_PATH=pkg/event-service
PHOTO_PATH=pkg/photo-service
GATEWAY_PATH=pkg/gateway

# The constructed docker image with organization name
EVENT_REPO=${ORGANIZATION}/${EVENT_IMAGE}
GATEWAY_REPO=${ORGANIZATION}/${GATEWAY_IMAGE}
PHOTO_REPO=${ORGANIZATION}/${PHOTO_IMAGE}

prod-gateway:
	docker build -f ${GATEWAY_PATH}/Dockerfile -t ${GATEWAY_REPO}:${GATEWAY_VERSION} .
	docker tag ${GATEWAY_REPO}:${GATEWAY_VERSION} ${GATEWAY_REPO}:latest
	docker push ${GATEWAY_REPO}:${GATEWAY_VERSION}
	docker push ${GATEWAY_REPO}:latest

# Local build
dev-gateway:
	@echo "Building development gateway"
	cd ${GATEWAY_PATH} && \
	go generate && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.dev -t ${GATEWAY_REPO}:${GATEWAY_VERSION} . && \
	rm -rf app

dev-event:
	@echo "Building development event gRPC server"
	cd ${EVENT_PATH} && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.dev -t ${EVENT_REPO}:${EVENT_VERSION} . && \
	rm -rf app

	@echo Pushing a new version ${EVENT_REPO}:${EVENT_VERSION} to dockerhub...
	docker tag ${EVENT_REPO}:${EVENT_VERSION} ${EVENT_REPO}:latest && \
	docker push ${EVENT_REPO}:${EVENT_VERSION} && \
	docker push ${EVENT_REPO}:latest
	@echo ${EVENT_REPO}:${EVENT_VERSION} is now available at dockerhub

dev-photo:
	@echo "Building development photo gRPC server"
	cd ${PHOTO_PATH} && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.dev -t ${PHOTO_REPO}:${PHOTO_VERSION} . && \
	rm -rf app

dev-all:
	@echo "Building development gateway"
	cd ${GATEWAY_PATH} && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.dev -t ${GATEWAY_REPO}:${GATEWAY_VERSION} . && \
	rm -rf app
	@echo "Building development event gRPC server"
	cd ${EVENT_PATH} && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.dev -t ${EVENT_REPO}:${EVENT_VERSION} . && \
	rm -rf app
	@echo "Building development photo gRPC server"
	cd ${PHOTO_PATH} && \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app && \
	docker build -f Dockerfile.dev -t ${PHOTO_REPO}:${PHOTO_VERSION} . && \
	rm -rf app

	docker tag ${GATEWAY_REPO}:${GATEWAY_VERSION} ${GATEWAY_REPO}:latest && \
	docker push ${GATEWAY_REPO}:${GATEWAY_VERSION} && \
	docker push ${GATEWAY_REPO}:latest

	docker tag ${EVENT_REPO}:${EVENT_VERSION} ${EVENT_REPO}:latest && \
	docker push ${EVENT_REPO}:${EVENT_VERSION} && \
	docker push ${EVENT_REPO}:latest

	docker tag ${PHOTO_REPO}:${PHOTO_VERSION} ${PHOTO_REPO}:latest && \
	docker push ${PHOTO_REPO}:${PHOTO_VERSION} && \
	docker push ${PHOTO_REPO}:latest

push-gateway:
	docker tag ${GATEWAY_REPO}:${GATEWAY_VERSION} ${GATEWAY_REPO}:latest && \
	docker push ${GATEWAY_REPO}:${GATEWAY_VERSION} && \
	docker push ${GATEWAY_REPO}:latest