include .env
export

event:
	PORT=${EVENT_PORT} \
	MGO_HOST=${MONGO_HOST} \
	MGO_USR=${MONGO_USER} \
	MGO_PWD=${MONGO_PASS} go run pkg/event-service/*.go

photo:
	PORT=${PHOTO_PORT} \
	MGO_HOST=${MONGO_HOST} \
	MGO_USR=${MONGO_USER} \
	MGO_PWD=${MONGO_PASS} go run pkg/photo-service/*.go

gateway:
	PORT=${GATEWAY_PORT} \
	PHOTO_ADDR=localhost${PHOTO_PORT} \
	EVENT_ADDR=localhost${EVENT_PORT} \
	AUTH0_JWK=${AUTH0_JWK} \
	AUTH0_ISS=${AUTH0_ISS} \
	AUTH0_AUD=${AUTH0_AUD} \
	AUTH0_WHITELIST=${AUTH0_WHITELIST} \
	TRACER=${GATEWAY_TRACER} go run pkg/gateway/*.go


