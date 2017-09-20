public-server:
	docker build -f pkg/server-public/Dockerfile -t alextanhongpin/public-event-server:0.0.1-beta .
	docker push alextanhongpin/public-event-server:0.0.1-beta

private-server:
	docker build -f pkg/server-private/Dockerfile -t alextanhongpin/private-event-server:0.0.1-beta .
	docker push alextanhongpin/private-event-server:0.0.1-beta

public-gateway:
	docker build -f pkg/gateway-public/Dockerfile -t alextanhongpin/public-event-gateway:0.0.1-beta .
	docker push alextanhongpin/public-event-gateway:0.0.1-beta

private-gateway:
	docker build -f pkg/gateway-private/Dockerfile -t alextanhongpin/private-event-gateway:0.0.1-beta .
	docker push alextanhongpin/private-event-gateway:0.0.1-beta


