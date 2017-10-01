GOOS=linux GOARCH=arm go build -o app main.go && 
docker build -f Dockerfile.production -t alextanhongpin/private-event-gateway:0.0.1-beta . &&
rm -rf app