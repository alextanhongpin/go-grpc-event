# go-grpc-event

Ensure you have the following in your `.env` file:

```bash
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=password
MONGO_INITDB_DATABASE=dbname


GATEWAY_PORT=:3000
PHOTO_PORT=:5000
EVENT_PORT=:9000

MONGO_HOST=localhost:27017
MONGO_USER=user
MONGO_PASS=password

GATEWAY_TRACER=gateway
EVENT_TRACER=event_service
PHOTO_TRACER=photo_service

AUTH0_JWK=<AUTH0_JWK>
AUTH0_ISS=<AUTH0_ISS>
AUTH0_AUD=<AUTH0_ISS>
AUTH0_WHITELIST=<AUTH0_ISS>
```


## Swagger

Swagger definitions are available at

- /swagger/event.swagger.json
- /swagger/photo.swagger.json
- /swagger/user.swagger.json

## Monitoring

- For serious services, alert immediately
- Else, just send a daily report of number of errors that occured daily
