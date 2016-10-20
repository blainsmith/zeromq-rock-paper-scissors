# ZeroMQ Rock/Paper/Scissors

Written for https://github.com/agnostechvalley/katas/tree/master/2016-10

Since ZeroMQ uses TCP as its protocol a server and client can be written in any language as long as it uses the `PAIR` socket type. Then the server and client can be started in any order. If you start the server first you do have to know the address of the server ahead of time.

## Running as a Server

Specify `GAMES` and `PORT` to start as a server and wait for a connection.

```bash
$ GAMES=7 PORT=1313 go run main.go
Address: tcp://10.0.0.13:1313
Games: 7
...
```

## Running as a Client

Specify `ADDRESS` of a running server.

```bash
$ ADDRESS=tcp://10.0.0.13:1313 go run main.go
Games: 7
...
```
