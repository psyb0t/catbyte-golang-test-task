docker-compose up -d

## API

`go run cmd/api/main.go`

### Example request

```
curl -X POST -d '{"sender": "me", "receiver": "you", "message": "test"}' http://localhost:8080/message
```

### Example response

```
{
	"message": "OK"
}
```

## Message Processor

`go run cmd/message-processor/main.go`

## Reporting API

`go run cmd/reporting-api/main.go`

### Example request

```
curl -X GET "http://localhost:8081/message/list?sender=me&receiver=you"
```

### Example response

```
[{
	"sender": "me",
	"receiver": "you",
	"message": "test2"
}, {
	"sender": "me",
	"receiver": "you",
	"message": "test1"
}, {
	"sender": "me",
	"receiver": "you",
	"message": "test"
}]
```

## Notes

Normally I would've added the services to the docker-compose file by creating images for each of them and for the services exposing HTTP I would've put in place a reverse proxy like Nginx or Traefik and each service would've just ran on port 80 internally and would've just accepted requests on the `/` route path.

Also normally all of the logic code i would've put in a separate `internal` directory whihe `cmd` would be responsible just for handling shutdown signals and running the internal code.

Here's an example of how I usually structure services: https://github.com/psyb0t/telegram-logger

### Subnote

I just realised I used vars instead of constants when defining global strings like this one `var rabbitMQURI = "amqp://user:password@localhost:7001/"`. My bad.
