## Setup

Notes:
* [Cobra](https://github.com/spf13/cobra) allows for sharing one binary with multiple applications (jwt_maker, seed_database, server, and twirp_server).
* [GoDotEnv](https://github.com/joho/godotenv) loads env vars from .env for local development

### .env

* Copy .env.SAMPLE to .env at the base of the project

### Create JWT Keypair

* `openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub`
* reference your keypair in .env

### Create the PostgreSQL Database

* Create a PostgreSQL database, and reference it in your .env
* `dbmate up`
* `go build -o charlie-go`
* `./charlie-go seed_database`

### Twirp Server

#### Auth

##### Create a JWT
`> go build -o charlie-go`

`> ./charlie-go generate_jwt 1h <users.id from database>`

* Claims (see `internal/jwt_maker/jwt_maker`)
  * exp: expiration (unix time)
  * iat: issued at timestamp (unix time)
  * iss: issuer ("charlie.com")
  * aud: audience ("user")
  * user_id: user_id (uuid, users.id from the database)

### A Note on Server Hooks
Create as many different `ChainHook` objects as you need. Here we create just one:

**cmd/server/main.go**

```
chainHooks := twirp.ChainHooks(
    provider.AuthHooks(),
)
```
With Twirp's generated code, each `ChainHook` object has access to these parts of the request / response lifecycle:

See [ServerHooks](https://github.com/twitchtv/twirp/blob/90c6a70b98cf6a201d6ebc060924cd6c0800f1bd/server_options.go#L96)

```
type ServerHooks struct {
	// RequestReceived is called as soon as a request enters the Twirp
	// server at the earliest available moment.
	RequestReceived func(context.Context) (context.Context, error)

	// RequestRouted is called when a request has been routed to a
	// particular method of the Twirp server.
	RequestRouted func(context.Context) (context.Context, error)

	// ResponsePrepared is called when a request has been handled and a
	// response is ready to be sent to the client.
	ResponsePrepared func(context.Context) context.Context

	// ResponseSent is called when all bytes of a response (including an error
	// response) have been written. Because the ResponseSent hook is terminal, it
	// does not return a context.
	ResponseSent func(context.Context)

	// Error hook is called when an error occurs while handling a request. The
	// Error is passed as argument to the hook.
	Error func(context.Context, Error) context.Context
}
```

#### Put it All Together

**cmd/twirp_server.go**

```
// POST http(s)://<host>/api/v1/charlie_go.CharlieGo/CreateItem
// POST http(s)://<host>/api/v1/charlie_go.CharlieGo/GetItem
handler := charlie_go.NewCharlieGoServer(provider, twirp.WithServerPathPrefix("/api/v1"), chainHooks)
mux.Handle(
    handler.PathPrefix(), twirp_server.AddJwtTokenToContext(
        handler,
    ),
)

http.ListenAndServe(httpPort, mux)
```

#### Run the Twirp Server

`> go build -o charlie-go`

`> ./charlie-go twirp_server`

Your server should be running on port 8081

##### Create an Item

```
> curl --location 'http://localhost:8081/api/v1/charlie_go.CharlieGo/CreateItem' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <your JWT>' \
--data '{
  "name": "Widget 08"
}'
```

##### Get an Item

```
> curl --location 'http://localhost:8081/api/v1/charlie_go.CharlieGo/GetItem' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer <your JWT>' \
--data '{
        "id": "<items.id owned by this user"
}'
```

#### Run the REST Server

`> go build -o charlie-go`

`> ./charlie-go rest_server`

##### Get items

```
> curl --location --request GET 'http://localhost:8080/api/v1/items' \
--header 'Content-Type: application/json' \
--header 'X-User-Id: <users.id from the database>' \
--data '{}'
```