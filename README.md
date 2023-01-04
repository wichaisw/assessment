# README
- โจทย์เดิมอยู่ในไฟล์ Assignment.md

# Run Application
- Replace {DATABASE_URL} in the command below with your connection string
- Replace {PORT} in the command with your desired port, post must be in format :port e.g. :2565
- Run this command
```
DATABASE_URL={DATABASE_URL} PORT={PORT} go run server.go
```

# Container
## Run Container
- Replace {APP_NAME} with your desired container name build your docker file with 
```
docker build -t {APP_NAME} .
```
## Start Application in the Container
- Replace the command variable according to the description then run the command below
  - Replace DATABASE_URL: your database connection string
  - Replace PORT: your desired port, post must be in format :port e.g. :2565
  - Replace APP_NAME: container name build your docker file with
  - Replace HOST_PORT: your host's port you want to expose
  - Replace CONTAINER_PORT: application's port in the container
```
docker run -e {DATABASE_URL} -e {PORT} -p {HOST_PORT}:{CONTAINER_PORT} {APP_NAME}
```
- e.g. docker run -e DATABASE_URL=postgres://xxx:xxx@host.com/database -e PORT=:2565 expense-app

# Testing
## Unit Testing
```
go test -v -tags=unit ./...
```

## Integration Testing in Docker Sandbox
```
docker compose -f docker-compose.test.yaml up
```