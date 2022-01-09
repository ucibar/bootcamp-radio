# Bootcamp Radio
`Bootcamp Radio` is a webrtc audio stream server. Users can live stream their voices to other users.

## Used Technologies
- [pion/webrtc](https://github.com/pion/webrtc) for audio streaming
- [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) for HTTP routing
- PostgreSQL for persistence layer
- Alpinejs, TailwindCSS and Axios used for frontend showcase

## Install
First create a PostgreSQL database, after that run with your config:
```
$ go build -o server ./api
$ DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASS=postgres DB_NAME=whereami ./server
```

You need to pass config variables as env to executable as seen above.

## Endpoints
- `GET /public` frontend for showcase:
- 
- `POST /api/auth/register` for register a user:

```
{
    "username": "yourUsername",
    "email": "yourEmail",
    "password": "yourPassword",
    "password_repeat": "yourPassword",
}
```

- `POST /api/auth/login` for login and get auth token:

```
{
    "username": "yourUsername",
    "password": "yourPassword",
}
```

- `GET /broadcast/list` get list of active broadcasts
- `POST /broadcast/create` create a broadcast, this endpoint needs a webRTC compatible client(eg. a browser):

```
{
    "broadcast_title": "title of your broadcast",
    "offer": "webrtc offer with base64 encoding",
}
```
**NOTE:** currently `broadcast_title` not implemented correctly.

- `POST /broadcast/join` join a broadcast, this endpoint needs a webRTC compatible client(eg. a browser):

```
{
    "broadcast_id": "broadcast id",
    "offer": "webrtc offer with base64 encoding",
}
```