include .env
export

server:
	PORT=${SERVER_PORT} MGO_HOST=${MONGO_HOST} MGO_USR=${MONGO_USER} MGO_PWD=${MONGO_PASS} go run pkg/server/*.go

gateway:
	go run pkg/gateway/main.go \
	-port=${GATEWAY_PORT} \
	-event-addr=localhost${SERVER_PORT} \
	-photo-addr=localhost${PHOTO_PORT} \
	-jwks_uri=${JWKS_URI} \
	-auth0_iss=${ISSUER} \
	-auth0_aud=${AUDIENCE} \
	-whitelisted_emails=${WHITELISTED_EMAILS}

photo:
	PORT=${PHOTO_PORT} MGO_HOST=${MONGO_HOST} MGO_USR=${MONGO_USER} MGO_PWD=${MONGO_PASS} go run pkg/photo-service/*.go